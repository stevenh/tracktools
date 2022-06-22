package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stevenh/tracktools/pkg/gopro"
)

// goproConvertCmd represents the convert command.
var goproConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts GoPro videos",
	Long:  `Converts GoProv videos between formats and joins multi chapters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load gopro convert from global config file.
		var cfg gopro.Config
		if err := loadConfig("gopro.convert", &cfg); err != nil {
			return err
		}

		// Command line overrides.
		f := cmd.Flags()
		source, err := f.GetString("source")
		if err != nil {
			return err
		}

		if source != "" {
			cfg.SourceDir = source
		}

		output, err := f.GetString("output")
		if err != nil {
			return err
		}

		if output != "" {
			cfg.OutputDir = output
		}

		p, err := gopro.NewProcessor(gopro.Cfg(cfg))
		if err != nil {
			return err
		}

		files, err := p.Process()
		if err != nil {
			return err
		}

		log.Info().Msg(fmt.Sprintf("Processed %d files\n", len(files)))
		for _, f := range files {
			log.Info().Msg(f)
		}

		return nil
	},
}

func init() { // nolint: gochecknoinits
	goproCmd.AddCommand(goproConvertCmd)

	f := goproConvertCmd.Flags()
	f.String("source", "", "override config source directory")
	f.String("output", "", "override config output directory")
}
