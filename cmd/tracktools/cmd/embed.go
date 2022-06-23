package cmd

import (
	_ "embed"
)

//go:embed .tracktools.toml
var defaultConfig []byte
