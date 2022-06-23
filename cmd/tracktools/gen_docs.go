//go:build generate

package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	"github.com/stevenh/tracktools/cmd/tracktools/cmd"
)

func main() {
	r := cmd.RootCmd()
	r.DisableAutoGenTag = true
	if err := doc.GenMarkdownTree(r, "./docs"); err != nil {
		log.Fatal(err)
	}
}
