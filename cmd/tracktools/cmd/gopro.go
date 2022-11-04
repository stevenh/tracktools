package cmd

import (
	"github.com/spf13/cobra"
)

// goproCmd represents the gopro command.
var goproCmd = &cobra.Command{
	Use:   "gopro",
	Short: "Provides commands for manipulating GoPro videos",
	Long:  `Provides a set of commands for manipulating GoPro videos.`,
}

func init() { //nolint: gochecknoinits
	rootCmd.AddCommand(goproCmd)
}
