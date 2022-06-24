package cmd

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cmdNameAnno = "tracktools_annotation_cmd"
)

// loadConfig loads a config section for cmd into cfg skipping
// values which have already been set by flags.
func loadConfig(cmd *cobra.Command, cfg any) error {
	name := cmdConfigName(cmd)
	data := viper.GetStringMap(name)
	log.Trace().Str("cmd", name).Fields(data).Msg("Loading config")

	// Remove entries which have been overridden by command line flags.
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if v, ok := f.Annotations[cmdNameAnno]; ok && v[0] == name {
			n := strings.ToLower(strings.ReplaceAll(f.Name, "-", ""))
			delete(data, n)
			log.Trace().Str("flag", n).Msg("skipped")
		}
	})

	// Decode the remaining config into cfg.
	if err := mapstructure.Decode(data, cfg); err != nil {
		return fmt.Errorf("load config: decode: %w", err)
	}

	log.Trace().
		Str("cmd", name).
		Str("cfg", fmt.Sprintf("%#v", cfg)).
		Msg("Loaded config")

	return nil
}

// annotate annotates all flags in fs with the config name.
func annotate(fs *pflag.FlagSet, name string) {
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Annotations == nil {
			f.Annotations = map[string][]string{}
		}
		f.Annotations[cmdNameAnno] = []string{name}
	})
}

// cmdConfigName returns the config name of cmd.
func cmdConfigName(cmd *cobra.Command) string {
	parts := strings.Split(cmd.CommandPath(), " ")
	if len(parts) == 1 {
		return ""
	}

	return strings.Join(parts[1:], ".")
}
