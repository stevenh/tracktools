package cmd

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose int

	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "tracktools",
		Short: "A set of tools for creating track videos",
		Long: `A set of tools for creating track videos including converting
between different track app formats and joining GoPro chaptered videos.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Error will already have been output.
		os.Exit(1)
	}
}

func init() { // nolint: gochecknoinits
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	cobra.OnInitialize(initConfig)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.tracktools.yaml)")
	pf.CountVarP(&verbose, "verbose", "v", "verbose output")
}

func initConfig() {
	switch verbose {
	case 1:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to determine home directory")
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".tracktools")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Print("Using config file:", viper.ConfigFileUsed())
	}
}
