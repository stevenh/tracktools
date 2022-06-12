package gopro

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	cfgNoInput = func() Config {
		c := *DefaultConfig
		c.Args = nil
		return c
	}()

	cfgNoBinary = func() Config {
		c := *DefaultConfig
		c.Binary = ""
		return c
	}()

	cfgNoSourceDir = func() Config {
		c := *DefaultConfig
		c.SourceDir = ""
		return c
	}()
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		err  error
	}{
		{
			name: "default",
			cfg:  *DefaultConfig,
		},
		{
			name: "no-input",
			cfg:  cfgNoInput,
			err:  configError(`Arg: -i ""`),
		},
		{
			name: "no-binary",
			cfg:  cfgNoBinary,
			err:  configError("Binary"),
		},
		{
			name: "no-source-dir",
			cfg:  cfgNoSourceDir,
			err:  configError("SourceDir"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			require.Equal(t, tc.err, err)
		})
	}
}
