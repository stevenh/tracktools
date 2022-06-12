package main

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/stevenh/gopro-tools/pkg/gopro"
)

func main() {
	var sourceDir, cfgFile string

	flag.StringVar(&sourceDir, "source", "", "source directory containing GoPro videos (config override)")
	flag.StringVar(&cfgFile, "cfg", "gpconvert.json", "configuration for joiner")
	flag.Parse()

	cfg, err := gopro.LoadConfig(cfgFile)
	if err != nil {
		log.Print(err)
		return
	}

	if sourceDir != "" {
		cfg.SourceDir = sourceDir
	}

	p, err := gopro.NewProcessor(gopro.Cfg(*cfg))
	if err != nil {
		log.Print(err)
		return
	}

	files, err := p.Process()
	if err != nil {
		log.Print(err)
		return
	}

	log.Info().Msg(fmt.Sprintf("Processed %d files\n", len(files)))
	for _, f := range files {
		log.Info().Msg(f)
	}
}
