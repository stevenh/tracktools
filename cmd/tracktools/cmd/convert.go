package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/stevenh/tracktools/pkg/convert"
	"github.com/stevenh/tracktools/pkg/laptimer"
	"github.com/stevenh/tracktools/pkg/trackaddict"
)

type convertCmd struct {
	decoder  string
	encoder  string
	compress bool
}

func (c *convertCmd) RunE(cmd *cobra.Command, args []string) (err error) { // nolint: nonamedreturns
	// Validate encoder / decoder.
	switch c.decoder {
	case "trackaddict":
	default:
		return fmt.Errorf("convert: unknown decoder: %q", c.decoder)
	}

	switch c.encoder {
	case "laptimer":
	default:
		return fmt.Errorf("convert: unknown encoder: %q", c.encoder)
	}

	// Open input / output if needed.
	var input io.Reader
	switch args[0] {
	case "-":
		input = os.Stdin
	default:
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("convert: input %w", err)
		}
		defer f.Close() // nolint: errcheck
		input = f
	}

	var output io.Writer
	switch args[1] {
	case "-":
		output = os.Stdout
	default:
		f, err := os.Create(args[1])
		if err != nil {
			return fmt.Errorf("convert: output: %w", err)
		}

		defer func() {
			// Check error is needed because of buffered writes.
			if cerr := f.Close(); cerr != nil && err == nil {
				err = cerr
			}
		}()
		output = f
	}

	return c.trackAddict2LapTimer(input, output)
}

// trackAddict2LapTimer decodes TrackAddict data from r converts it and
// encodes it to w in LapTimer format.
func (c *convertCmd) trackAddict2LapTimer(r io.Reader, w io.Writer) error {
	dec, err := trackaddict.NewDecoder(r)
	if err != nil {
		return fmt.Errorf("convert: new trackaddict decoder: %w", err)
	}

	sess, err := dec.Decode()
	if err != nil {
		return fmt.Errorf("convert: trackaddict decode: %w", err)
	}

	ta, err := convert.NewTrackAddict()
	if err != nil {
		return fmt.Errorf("convert: new trackaddict converter: %w", err)
	}

	db, err := ta.LapTimer(sess)
	if err != nil {
		return fmt.Errorf("convert: laptimer: %w", err)
	}

	enc, err := laptimer.NewEncoder(w)
	if err != nil {
		return fmt.Errorf("convert: new laptimer encoder: %w", err)
	}

	if err = enc.Encode(db); err != nil {
		return fmt.Errorf("convert: laptimer encode: %w", err)
	}

	return nil
}

// addConvertCmd adds the convert command.
func addConvertCmd() {
	c := convertCmd{}
	cmd := &cobra.Command{
		Use:   "convert input-file output-file",
		Short: "Convert between track app formats.",
		Long:  `Convert between different track app logging formats`,
		Args:  cobra.ExactArgs(2),
		RunE:  c.RunE,
	}

	f := cmd.Flags()
	f.StringVar(&c.decoder, "decoder", "trackaddict", "The decoder to use for the input")
	f.StringVar(&c.encoder, "encoder", "laptimer", "The encoder of use for the output")
	f.BoolVar(&c.compress, "compress", false, "Compress the output")

	rootCmd.AddCommand(cmd)
}

func init() { // nolint: gochecknoinits
	addConvertCmd()
}
