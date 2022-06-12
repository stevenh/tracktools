package gopro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

var (
	// DefaultConfig is the default configuration if not specified.
	DefaultConfig = &Config{
		Binary: "ffmpeg",
		Args: []string{
			"-y",
			"-safe", "0",
			"-f", "concat",
			"-i", "",
			"-c:a", "copy",
			"-c:d", "copy",
			"-c:v", "libx264",
			"-vf", "scale=1920:1080",
			"-copy_unknown",
			"-map_metadata", "0",
			"-movflags", "use_metadata_tags",
			"-map", "0:v",
			"-map", "0:a",
			"-map", "0:m:handler_name:\tGoPro TCD",
			"-map", "0:m:handler_name:\tGoPro MET",
			"-map", "0:m:handler_name:\tGoPro SOS",
		},
		SourceDir:      ".",
		OutputTemplate: "{{.Name}}-JOINED{{.Ext}}",
		LogLevel:       "warn",
	}
)

// Config represents a GoPro video processing configuration.
type Config struct {
	LogLevel       string
	SourceDir      string
	Binary         string
	Args           []string
	Skip           map[string]bool
	OutputTemplate string
	OutputDir      string
	Overwrite      bool

	logLevel   zerolog.Level
	inputIndex int
	outputTmpl *template.Template
}

// Validate validates c calculating internal values.
func (c *Config) Validate() error {
	switch {
	case c.SourceDir == "":
		return configError("SourceDir")
	case c.Binary == "":
		return configError("Binary")
	}

	tmpl, err := template.New("output").Parse(c.OutputTemplate)
	if err != nil {
		return fmt.Errorf("output template: %w", err)
	}

	c.outputTmpl = tmpl

	if c.LogLevel == "" {
		c.logLevel = zerolog.WarnLevel
	} else {
		c.logLevel, err = zerolog.ParseLevel(c.LogLevel)
		if err != nil {
			return fmt.Errorf("parse log level %q: %w", c.LogLevel, err)
		}
	}

	// Locate the index of the input argument.
	var iidx int
	for i, v := range c.Args {
		switch v {
		case "-i":
			iidx = i
		case "":
			if i == iidx+1 {
				c.inputIndex = i

				return nil
			}
		}
	}

	return configError(`Arg: -i ""`)
}

// Output returns the output filename calculated from .OutputTmpl and name.
func (c *Config) Output(name string) (string, error) {
	var buf bytes.Buffer
	ext := filepath.Ext(name)
	t := OutputFile{
		Name: strings.TrimSuffix(name, ext),
		Ext:  ext,
	}
	if err := c.outputTmpl.Execute(&buf, t); err != nil {
		return "", fmt.Errorf("output template: %w", err)
	}

	return buf.String(), nil
}

// OutputFile represents the template used for output file processing.
type OutputFile struct {
	// Name excluding extension.
	Name string

	// Extension including the dot.
	Ext string
}

// LoadConfig loads a json configuration from file.
func LoadConfig(file string) (*Config, error) {
	f, err := baseFS.Open(file)
	if err != nil {
		return nil, fmt.Errorf("config load: %w", err)
	}

	defer f.Close() // nolint: errcheck

	cfg := &Config{}
	d := json.NewDecoder(f)
	if err = d.Decode(cfg); err != nil {
		return nil, fmt.Errorf("config decode: %w", err)
	}

	return cfg, nil
}
