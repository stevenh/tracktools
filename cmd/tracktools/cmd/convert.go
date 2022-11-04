package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/stevenh/tracktools/pkg/convert"
	"github.com/stevenh/tracktools/pkg/laptimer"
	"github.com/stevenh/tracktools/pkg/trackaddict"
)

// Encoders & Decoders.
const (
	trackAddict = "trackaddict"
	lapTimer    = "laptimer"
)

type convertCmd struct {
	Decoder  string
	Encoder  string
	Compress bool

	// Convert options.
	Track     string
	Vehicle   string
	Tags      []string
	Note      string
	StartDate date
}

func (c *convertCmd) RunE(cmd *cobra.Command, args []string) (err error) { //nolint: nonamedreturns
	if err := loadConfig(cmd, c); err != nil {
		return err
	}

	// Validate encoder / decoder.
	switch c.Decoder {
	case trackAddict:
	default:
		return fmt.Errorf("convert: unknown decoder: %q", c.Decoder)
	}

	switch c.Encoder {
	case lapTimer:
	default:
		return fmt.Errorf("convert: unknown encoder: %q", c.Encoder)
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
		defer f.Close() //nolint: errcheck
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

	taOpts := []convert.Option{
		convert.TrackOpt(c.Track),
		convert.VehicleOpt(c.Vehicle),
		convert.TagsOpt(c.Tags...),
		convert.NoteOpt(c.Note),
		convert.StartDateOpt(time.Time(c.StartDate)),
	}
	ta, err := convert.NewTrackAddict(taOpts...)
	if err != nil {
		return fmt.Errorf("convert: new trackaddict converter: %w", err)
	}

	db, err := ta.LapTimer(sess)
	if err != nil {
		return fmt.Errorf("convert: laptimer: %w", err)
	}

	encOpts := []laptimer.EncoderOpt{}
	if c.Compress {
		encOpts = append(encOpts, laptimer.Compress())
	}
	enc, err := laptimer.NewEncoder(w, encOpts...)
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

	fs := cmd.Flags()
	fs.StringVar(&c.Decoder, "decoder", "", "Override Decoder for the input")
	fs.StringVar(&c.Encoder, "encoder", "", "Override Encoder for the output")
	fs.StringVar(&c.Track, "track", "", "Override Track for the output")
	fs.StringVar(&c.Vehicle, "vehicle", "", "Override Vehicle for the output")
	fs.StringArrayVar(&c.Tags, "tags", nil, "Override Tags for the output")
	fs.StringVar(&c.Note, "note", "", "Override Note for the output")
	fs.BoolVar(&c.Compress, "compress", false, "Override Compress option for output")
	fs.Var(&c.StartDate, "start-date", "Override StartDate option for output (format YYYY-MM-DD)")
	annotate(fs, "convert")

	rootCmd.AddCommand(cmd)
}

func init() { //nolint: gochecknoinits
	addConvertCmd()
}
