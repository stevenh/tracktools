package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stevenh/tracktools/pkg/gopro/gpmf"
	"github.com/stevenh/tracktools/pkg/gopro/gpmf/geo"
)

// goproLapTimesCmd represents the gopro laptimes command.
type goproLapTimesCmd struct {
	Start     Start
	Tolerance float64

	p     *geo.Processor
	found int
}

func (c *goproLapTimesCmd) RunE(cmd *cobra.Command, args []string) error {
	if err := loadConfig(cmd, c); err != nil {
		return err
	}

	c.Start.calculate()
	c.p = geo.NewProcessor(geo.Tolerance(c.Tolerance))

	dec := &gpmf.Decoder{}
	for _, fn := range args {
		if err := c.process(dec, fn); err != nil {
			return err
		}
	}

	return nil
}

func (c *goproLapTimesCmd) process(dec *gpmf.Decoder, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("laptimes: open %w", err)
	}

	defer f.Close() // nolint: errcheck

	data, err := dec.Decode(f)
	if err != nil {
		return fmt.Errorf("laptimes: decode %q: %w", file, err)
	}

	if err := gpmf.Walk(data, c.walk); err != nil {
		return fmt.Errorf("laptimes: walk %q: %w", file, err)
	}

	if c.found == 0 {
		return fmt.Errorf("laptimes: walk %q: no laps found", file)
	}

	return nil
}

func (c *goproLapTimesCmd) walk(e *gpmf.Element) error {
	data, ok := e.Data.(gpmf.GPSData)
	if !ok {
		return nil
	}

	for _, v := range data {
		if c.p.OnLine(v.Latitude, v.Longitude, c.Start.lat1, c.Start.lon1, c.Start.lat2, c.Start.lon2) {
			log.Info().Object("gps", v).Msg("start line passed")
			c.found++
		}
	}

	return nil
}

func addGoproLapTimes() {
	c := goproLapTimesCmd{}
	cmd := &cobra.Command{
		Use:   "laptimes [file1] ... [fileN]",
		Short: "LapTimes reports laptimes of GoPro videos",
		Long:  `LapTimes reports laptimes of GoPro based on the GPS metadata information.`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  c.RunE,
	}

	fs := cmd.Flags()
	fs.Float64Var(&c.Start.Latitude, "latitude", 0, "override start latitude")
	fs.Float64Var(&c.Start.Longitude, "longitude", 0, "override start longitude")
	fs.Float64Var(&c.Start.Bearing, "bearing", 0, "override start bearing")
	fs.Float64Var(&c.Start.Distance, "distance", 0, "override start distance")
	fs.Float64Var(&c.Tolerance, "tolerance", 0, "override tolerance")
	annotate(fs, "gopro.laptimes")

	goproCmd.AddCommand(cmd)
}

func init() { // nolint: gochecknoinits
	addGoproLapTimes()
}
