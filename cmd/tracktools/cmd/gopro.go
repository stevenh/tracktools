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

func init() { // nolint: gochecknoinits
	rootCmd.AddCommand(goproCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// goproCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// goproCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
