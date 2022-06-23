package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stevenh/tracktools/pkg/gopro"
)

// goproConvertCmd represents the convert command.
type goproConvertCmd struct {
	cfg gopro.Config
}

func (c *goproConvertCmd) RunE(cmd *cobra.Command, args []string) error {
	if err := loadConfig("gopro.convert", &c.cfg, cmd.Flags()); err != nil {
		return err
	}

	p, err := gopro.NewProcessor(gopro.Cfg(c.cfg))
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
}

func addGoproConvert() {
	c := goproConvertCmd{}

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Converts GoPro videos",
		Long:  `Converts GoPro videos between formats and joins multi chapters.`,
		RunE:  c.RunE,
	}

	fs := cmd.Flags()
	fs.StringVar(&c.cfg.SourceDir, "source-dir", "", "override source directory")
	fs.StringVar(&c.cfg.OutputDir, "output-dir", "", "override output directory")
	annotate(fs, "gopro.convert")

	goproCmd.AddCommand(cmd)
}

func init() { // nolint: gochecknoinits
	addGoproConvert()
}
