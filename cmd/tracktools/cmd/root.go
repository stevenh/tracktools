package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = ".tracktools"
	defaultConfigType = "toml"
	debugLevel        = 1
	traceLevel        = 2
)

var (
	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "tracktools",
		Short: "A set of tools for creating track videos",
		Long: `A set of tools for creating track videos including converting
between different track app formats and joining GoPro chaptered videos.`,
	}
)

type rootCommand struct {
	Config  string
	Verbose int
}

func newRoot() {
	r := &rootCommand{}
	rootCmd.PersistentPreRunE = r.PersistentPreRunE

	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&r.Config, "config", "c", "", "config file (Default .tracktools.toml)")
	pf.CountVarP(&r.Verbose, "verbose", "v", "verbose output")
}

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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	newRoot()
}

// PersistentPreRunE initialises our config.
func (r *rootCommand) PersistentPreRunE(cmd *cobra.Command, args []string) error {
	v := viper.GetViper()

	switch r.Verbose {
	case debugLevel:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case traceLevel:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if err := r.setConfigSpec(v); err != nil {
		return err
	}

	v.AutomaticEnv()

	if err := r.loadConfig(v); err != nil {
		return err
	}

	return nil
}

// configSpec sets our config spec on v.
// If we have a specified Config this takes preference otherwise
// we search in the following directories in order:
// * Current working directory
// * Users home directory.
func (r *rootCommand) setConfigSpec(v *viper.Viper) error {
	if r.Config != "" {
		// Use config file from the flag.
		v.SetConfigFile(r.Config)
		return nil
	}

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("home directory: %w", err)
	}

	v.AddConfigPath(".")
	v.AddConfigPath(home)
	v.SetConfigName(defaultConfigName)

	return nil
}

// loadConfig loads the config from file falling back to
// the default embedded one if the config file location wasn't
// specified and we didn't find a file, so we get sensible defaults.
func (r *rootCommand) loadConfig(v *viper.Viper) error {
	if err := v.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) {
			return fmt.Errorf("load config: %w", err)
		}

		if r.Config != "" {
			// Config was specified so don't fall back to default.
			return err
		}

		// ReadConfig needs to be told the type of config to expect.
		v.SetConfigType(defaultConfigType)

		buf := bytes.NewBuffer(defaultConfig)
		if err = v.ReadConfig(buf); err != nil {
			return fmt.Errorf("load default config: %w", err)
		}

		// If the config type is invalid it will just produce no results
		// so check we got a valid config.
		if len(v.AllKeys()) == 0 {
			return fmt.Errorf("load default config: %w", err)
		}

		log.Print("Using default embedded config")
		log.Trace().Msg(string(defaultConfig))

		return nil
	}

	log.Print("Using config file: ", v.ConfigFileUsed())

	return nil
}

// RootCmd returns the root command for doc generation.
func RootCmd() *cobra.Command {
	return rootCmd
}
