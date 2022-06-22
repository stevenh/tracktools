package cmd

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func loadConfig(name string, cfg interface{}) error {
	data := viper.GetStringMap(name)
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return fmt.Errorf("decode cfg: %w", err)
	}

	return nil
}
