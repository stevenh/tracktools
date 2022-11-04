package gopro

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

var (
	// baseFS is the file system used by Processor which defaults
	// to os file system.
	baseFS procFS = osFS{}
)

// Processor provides the ability to process GoPro video files.
// If it finds chaptered files they are joined during processing.
type Processor struct {
	cfg      *Config
	log      zerolog.Logger
	matchers []*Matcher
	handler  func(exe string, args ...string) error
}

// Option represents an configuration option for Processor.
type Option func(*Processor) error

// Cfg sets the config for the Processor.
// Either CfgFile or Cfg must be specified.
func Cfg(cfg Config) Option {
	return func(p *Processor) error {
		p.cfg = &cfg

		return p.cfg.Validate()
	}
}

// CfgFile loads the config from file for the Processor.
// Either CfgFile or Cfg must be specified.
func CfgFile(file string) Option {
	return func(p *Processor) error {
		c := &Config{}
		if err := c.Load(file); err != nil {
			return err
		}

		p.cfg = c

		return p.cfg.Validate()
	}
}

// Output sets output for log message for processor.
// Default is os.Stderr.
func Output(w io.Writer) Option {
	return func(p *Processor) error {
		p.log = p.log.Output(w)

		return nil
	}
}

// DefaultHandler is the default Handler function
// which runs ffmpeg wiring up Stdout and Stderr.
func DefaultHandler(exe string, args ...string) error {
	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s run: %w", exe, err)
	}

	return nil
}

// Handler sets a custom handler for running ffmpeg.
// Default is DefaultHandler.
func Handler(f func(exe string, args ...string) error) Option {
	return func(p *Processor) error {
		p.handler = f

		return nil
	}
}

// NewProcessor creates and returns a configured Processor.
func NewProcessor(options ...Option) (*Processor, error) {
	p := &Processor{
		log:      zerolog.New(os.Stderr),
		matchers: []*Matcher{Hero5, Hero10},
		cfg:      DefaultConfig,
		handler:  DefaultHandler,
	}

	for _, f := range options {
		if err := f(p); err != nil {
			return nil, err
		}
	}

	p.log = p.log.Level(p.cfg.logLevel)

	return p, nil
}

// fileSets searches the SourceDir and returns a map of FileSets.
func (p *Processor) fileSets() (map[string]*FileSet, error) {
	sets := make(map[string]*FileSet)
	if err := fs.WalkDir(baseFS, p.cfg.SourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		switch {
		case path == p.cfg.SourceDir:
			// Root dir, nothing to do.
			return nil
		case d.IsDir():
			// Don't process subdirectories.
			return filepath.SkipDir
		}

		fn := filepath.Base(path)
		for _, m := range p.matchers {
			f, err := m.Match(fn)
			if err != nil {
				if errors.Is(err, ErrNoMatch) {
					continue
				}

				return err
			}

			s, ok := sets[f.Index]
			if !ok {
				s = &FileSet{Number: f.Index}
				sets[f.Index] = s
			}
			s.Chapter(f)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("load file sets: %w", err)
	}

	return sets, nil
}

// Process finds and processes video files according to
// the Processor Config returning the resultant files.
// Files can be returned even if an error occurs.
func (p *Processor) Process() ([]string, error) {
	sets, err := p.fileSets()
	if err != nil {
		return nil, err
	}

	if len(sets) == 0 {
		return nil, ErrNoFiles
	}

	files := make([]string, 0, len(sets))
	for _, s := range sets {
		f, err := p.processSet(s)
		if err != nil {
			return files, err
		} else if f == "" {
			continue
		}

		files = append(files, f)
	}

	return files, nil
}

// writeInput writes a ffmpeg concat compatible input data to w for s.
func (p *Processor) writeInput(w io.WriteCloser, s *FileSet) error {
	for _, f := range s.Chapters {
		if _, err := fmt.Fprintf(w, "file '%s'\n", filepath.Join(p.cfg.SourceDir, f.Name)); err != nil {
			return fmt.Errorf("input file data: %w", err)
		}
	}

	return w.Close()
}

// Process processes s.
func (p *Processor) processSet(s *FileSet) (string, error) {
	p.log.Print("processing:", s)
	if err := s.Chapters.Validate(); err != nil {
		return "", err
	}

	if p.cfg.Skip(s.Chapters[0].Name) {
		p.log.Print("skip:", s.Chapters[0].Name)
		return "", nil
	}

	// TODO(steve): Avoid use of concat if there's only one file.

	f, err := baseFS.CreateTemp("", "gopro-process")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary input file: %w", err)
	}
	defer baseFS.Remove(f.Name()) //nolint: errcheck

	if err := p.writeInput(f, s); err != nil {
		return "", fmt.Errorf("write input file %q: %w", f.Name(), err)
	}

	// Templated naming.
	output, err := p.cfg.Output(s.Chapters[0].Name)
	if err != nil {
		return "", err
	}

	switch p.cfg.OutputDir {
	case ".":
		// Current directory no path needed.
	case "":
		// Default to source directory.
		output = filepath.Join(p.cfg.SourceDir, output)
	default:
		output = filepath.Join(p.cfg.OutputDir, output)
	}

	if _, err := baseFS.Stat(output); err == nil {
		if !p.cfg.Overwrite {
			p.log.Print("Skip existing:", output)
			return "", nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("stat output: %w", err)
	}

	// Set the input file name by index.
	p.cfg.Args[p.cfg.inputIndex] = f.Name()

	args := append(p.cfg.Args, output) //nolint: gocritic
	p.log.Print("handle:", p.cfg.Binary, args)
	if err = p.handler(p.cfg.Binary, args...); err != nil {
		return "", err
	}

	// Set the file times to that of first source file so that
	// file sorting by time matches that of the originals.
	fi, err := baseFS.Stat(filepath.Join(p.cfg.SourceDir, s.Chapters[0].Name))
	if err != nil {
		return "", fmt.Errorf("first source file stat: %w", err)
	}

	if err = baseFS.Chtimes(output, time.Now(), fi.ModTime()); err != nil {
		return "", fmt.Errorf("set output times: %w", err)
	}

	return output, nil
}
