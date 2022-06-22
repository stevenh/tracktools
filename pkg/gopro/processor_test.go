package gopro

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"testing/fstest"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type tmpFile struct {
	dir     string
	pattern string
	data    []byte
}

func (f *tmpFile) Write(data []byte) (int, error) {
	f.data = append(f.data, data...)
	return len(data), nil
}

func (f *tmpFile) Close() error {
	return nil
}

func (f *tmpFile) Name() string {
	return filepath.Join(f.dir, f.pattern)
}

type testFS struct {
	fstest.MapFS
	tempFiles    []*tmpFile
	createdFiles []string
}

func newTestFS() *testFS {
	return &testFS{
		MapFS: make(fstest.MapFS),
	}
}

func (t *testFS) CreateTemp(dir, pattern string) (tempFile, error) { // nolint: ireturn
	f := &tmpFile{dir: dir, pattern: pattern}
	t.tempFiles = append(t.tempFiles, f)
	return f, nil
}

func (t *testFS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	f, ok := t.MapFS[name]
	if !ok {
		return &fs.PathError{
			Op:   "chtimes",
			Path: name,
			Err:  fs.ErrNotExist,
		}
	}

	f.ModTime = mtime

	return nil
}

func (t *testFS) Remove(name string) error {
	if _, ok := t.MapFS[name]; !ok {
		return &fs.PathError{
			Op:   "remove",
			Path: name,
			Err:  fs.ErrNotExist,
		}
	}

	delete(t.MapFS, name)

	return nil
}

func (t *testFS) Create(name string) {
	f := &fstest.MapFile{ModTime: time.Now()}
	t.MapFS[name] = f
	t.createdFiles = append(t.createdFiles, name)
}

func (p *Processor) testLog() {
	p.log.Warn().Msg("test")
}

func TestNewProcessor(t *testing.T) {
	validCfg := &Config{
		SourceDir:  ".",
		Binary:     "ffmpeg",
		Args:       []string{"-i", ""},
		SkipNames:  []string{"skip1", "skip2"},
		inputIndex: 1,
	}
	validData, err := json.Marshal(validCfg)
	require.NoError(t, err)

	invalidCfg := validCfg
	invalidCfg.SourceDir = ""
	invalidData, err := json.Marshal(invalidCfg)
	require.NoError(t, err)

	// Use an in memory file system for testing.
	baseFS = &testFS{
		MapFS: fstest.MapFS{
			"valid-config.json": {
				Data: validData,
			},
			"invalid-config.json": {
				Data: invalidData,
			},
		},
	}

	tests := []struct {
		name     string
		options  []Option
		expected *Config
		err      error
	}{
		{
			name:     "default-config",
			expected: DefaultConfig,
		},
		{
			name:    "config-error",
			options: []Option{Cfg(cfgNoBinary)},
			err:     configError("Binary"),
		},
		{
			name:    "config-error",
			options: []Option{Cfg(cfgNoBinary)},
			err:     configError("Binary"),
		},
		{
			name:    "valid-cfg-file",
			options: []Option{CfgFile("valid-config.json")},
			expected: &Config{
				SourceDir:  ".",
				Binary:     "ffmpeg",
				Args:       []string{"-i", ""},
				SkipNames:  []string{"skip1", "skip2"},
				inputIndex: 1,
				logLevel:   zerolog.WarnLevel,
				skip: map[string]struct{}{
					"skip1": struct{}{},
					"skip2": struct{}{},
				},
			},
		},
		{
			name:    "invalid-cfg-file",
			options: []Option{CfgFile("invalid-config.json")},
			err:     configError("SourceDir"),
		},
		{
			name:    "missing-cfg-file",
			options: []Option{CfgFile("missing-config.json")},
			err: fmt.Errorf("config load: %w", &fs.PathError{
				Op:   "open",
				Path: "missing-config.json",
				Err:  fs.ErrNotExist,
			}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			p, err := NewProcessor(append(tc.options, Logger(&buf))...)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)
			p.testLog()

			p.cfg.outputTmpl = nil // This is calculated so don't try and match.

			require.Equal(t, tc.expected, p.cfg)
			require.Equal(t, `{"level":"warn","message":"test"}
`, buf.String())
		})
	}
}

// procOut is an Option which sets procOut for a Processor.
func procOut(w io.Writer) Option {
	return func(p *Processor) error {
		p.procOut = w

		return nil
	}
}

// processHelper processes the output from the helper process.
// It creates the required files in fs and signals the helper to
// quit once complete.
func processHelper(r io.Reader, fs *testFS) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		p := strings.SplitN(s.Text(), ":", 2)
		if len(p) != 2 {
			return fmt.Errorf("unexpected line %q", s.Text())
		}

		pid, err := strconv.Atoi(p[0])
		if err != nil {
			return fmt.Errorf("unexpected pid in line %q", s.Text())
		}

		// File must be created before the process exists.
		fs.Create(p[1])

		if err = syscall.Kill(pid, syscall.SIGINT); err != nil {
			return fmt.Errorf("signal pid %d:%w", pid, err)
		}
	}

	return s.Err()
}

func TestProcessorProcess(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		exe      string
		files    []string
		expected map[string]string
		err      error
	}{
		{
			name: "single",
			exe:  "ffmpeg",
			files: []string{
				"GOPR1234.mp4",
			},
			expected: map[string]string{
				"GOPR1234.mp4": "file 'GOPR1234.mp4'\n",
			},
		},
		{
			name: "multiple",
			exe:  "ffmpeg",
			files: []string{
				"GOPR0001.mp4",
				"GOPR0002.mp4",
			},
			expected: map[string]string{
				"GOPR0001.mp4": "file 'GOPR0001.mp4'\n",
				"GOPR0002.mp4": "file 'GOPR0002.mp4'\n",
			},
		},
		{
			name: "chapters",
			exe:  "ffmpeg",
			files: []string{
				"GOPR0001.mp4",
				"GP010001.mp4",
			},
			expected: map[string]string{
				"GOPR0001.mp4": "file 'GOPR0001.mp4'\nfile 'GP010001.mp4'\n",
			},
		},
		{
			name: "no-files",
			exe:  "ffmpeg",
			err:  ErrNoFiles,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tfs := newTestFS()
			for _, n := range tc.files {
				tfs.MapFS[n] = &fstest.MapFile{ModTime: now}
			}

			baseFS = tfs

			// Create a config which points to ourselves so we
			// can control the behaviour.
			cfg := *DefaultConfig
			cfg.Binary = os.Args[0]
			cfg.Args = append([]string{"-test.run=TestHelperProcess", "--", tc.exe}, cfg.Args...)
			cfg.Env = []string{"GO_WANT_HELPER_PROCESS=1"}

			// Process the output of our helper.
			pr, pw, err := os.Pipe()
			require.NoError(t, err)

			var once sync.Once
			closeWriter := func() {
				// Close our copy of the files to allow pipe to trigger.
				pw.Close() // nolint: errcheck
			}

			errs := make(chan error, 1)
			go func() {
				errs <- processHelper(pr, tfs)
				once.Do(closeWriter)
				close(errs)
			}()

			p, err := NewProcessor(Cfg(cfg), procOut(pw))

			require.NoError(t, err)

			files, err := p.Process()
			once.Do(closeWriter)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)
			for err := range errs {
				require.NoError(t, err)
			}

			require.Equal(t, tfs.createdFiles, files)
			require.Len(t, tfs.tempFiles, len(files))
			require.Len(t, tfs.createdFiles, len(files))

			for i, name := range tfs.createdFiles {
				// Check the file has the expected modification time.
				f, ok := tfs.MapFS[name]
				require.True(t, ok)
				require.Equal(t, now, f.ModTime)

				first := strings.Replace(name, "-JOINED", "", 1)
				expected, ok := tc.expected[first]
				require.True(t, ok, "missing file %q", name)

				tf := tfs.tempFiles[i]
				require.Equal(t, expected, string(tf.data))
			}
		})
	}
}

// TestHelperProcess isn't a real test, it's just used as a
// helper process.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		t.Skip("process helper")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "helper no command")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "ffmpeg":
		fmt.Printf("%d:%s\n", os.Getpid(), args[len(args)-1])
		os.Stdout.Close() // nolint: errcheck

		// Wait for signal to exit.
		select {
		case <-c:
		case <-time.After(time.Second * 10):
			fmt.Fprintln(os.Stderr, "helper timeout")
		}
		os.Exit(0)
	case "ffmpeg-fail":
		fmt.Fprintln(os.Stderr, "helper fail")
		os.Exit(2)
	default:
		fmt.Fprintf(os.Stderr, "helper unknown command %q\n", cmd)
		os.Exit(2)
	}
}
