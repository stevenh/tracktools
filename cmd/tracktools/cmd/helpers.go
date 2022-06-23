package cmd

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cmdName = "tracktools_annotation_cmd"
)

// loadConfig loads a config section into cfg skipping values which
// have already been set by flags.
func loadConfig(name string, cfg interface{}, fs *pflag.FlagSet) error {
	data := viper.GetStringMap(name)
	fs.VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			return
		}

		if v, ok := f.Annotations[cmdName]; ok && v[0] == name {
			n := strings.ToLower(strings.ReplaceAll(f.Name, "-", ""))
			delete(data, n)
		}
	})

	if err := mapstructure.Decode(data, &cfg); err != nil {
		return fmt.Errorf("decode cfg: %w", err)
	}

	return nil
}

// configName returns the fully qualified name of f according to its
// annotations.
func configName(f *pflag.Flag) string {
	if v, ok := f.Annotations[cmdName]; ok {
		n := strcase.ToCamel(f.Name)
		return fmt.Sprintf("%s.%s", v[0], n)
	}

	return f.Name
}

// annotate annotates all flags in fs with the command name.
func annotate(fs *pflag.FlagSet, name string) {
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Annotations == nil {
			f.Annotations = map[string][]string{}
		}
		f.Annotations[cmdName] = []string{name}
	})
}
