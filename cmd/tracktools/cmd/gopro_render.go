package cmd

import (
	"fmt"
	"os"

	"github.com/golang/geo/s2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/stevenh/tracktools/pkg/gopro/gpmf"
	"github.com/stevenh/tracktools/pkg/image"
)

// goproRenderCmd represents the gopro render command.
type goproRenderCmd struct {
	MinDoP        float64
	MinGood       int
	Width, Height int
	Start         Start

	data []s2.LatLng
	good int
}

func (c *goproRenderCmd) RunE(cmd *cobra.Command, args []string) error {
	if err := loadConfig(cmd, c); err != nil {
		return err
	}

	c.Start.calculate()

	dec := &gpmf.Decoder{}
	f, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("render: open %w", err)
	}

	defer f.Close() // nolint: errcheck

	data, err := dec.Decode(f)
	if err != nil {
		return fmt.Errorf("render: decode %q: %w", args[0], err)
	}

	if err := gpmf.Walk(data, c.walk); err != nil {
		return fmt.Errorf("render: walk %q: %w", args[0], err)
	} else if len(c.data) == 0 {
		return fmt.Errorf("render: walk %q: no gps data found", args[0])
	}

	r := image.New(c.Width, c.Height)
	r.Provider("tracktools")
	r.Start(c.Start.lat1, c.Start.lon1, c.Start.lat2, c.Start.lon2)
	r.AddPath(c.data, image.Red)

	if err := r.Render(args[1]); err != nil {
		return fmt.Errorf("render: image %q: %w", args[1], err)
	}

	log.Info().Str("file", args[1]).Msg("image rendered")

	return nil
}

// walk is a gpmf.WalkFunc which looks for GPS data, validates and
// stores for rendering if it passes.
// Validation is based off the GPS Dilution of Precision with only values
// above MinDoP being used. MinGood is also used to filter out bad data
// close to the start of the dataset.
func (c *goproRenderCmd) walk(e *gpmf.Element) error {
	data, ok := e.Data.(gpmf.GPSData)
	if !ok {
		return nil
	}

	for _, v := range data {
		var dop gpmf.GPSDoP = 100 // Default to bad data.
		if v, ok := e.MetadataByKey(gpmf.KeyGSPDoP); ok {
			if f, ok := v.(gpmf.GPSDoP); ok {
				dop = f
			}
		}

		if float64(dop) > c.MinDoP {
			// Bad data.
			c.good--

			switch {
			case c.good < 0:
				// Prevent it going negative.
				c.good = 0
				fallthrough
			case c.good < c.MinGood:
				// We have less than MinGood delete all previous
				// data to avoid a poor quality path in the render.
				c.data = c.data[:0]
				continue
			}
		} else {
			c.good++
		}

		c.data = append(c.data,
			s2.LatLngFromDegrees(v.Latitude, v.Longitude),
		)
	}

	return nil
}

func addGoproRender() {
	c := goproRenderCmd{data: make([]s2.LatLng, 0, 100)}
	cmd := &cobra.Command{
		Use:   "render [input mp4] [output image]",
		Short: "Renders image of GoPro GPS data",
		Long:  `Renders a image of GoPro GPS data.`,
		Args:  cobra.ExactArgs(2),
		RunE:  c.RunE,
	}

	fs := cmd.Flags()
	fs.IntVar(&c.MinGood, "min-good", 0, "override minimum good measurements")
	fs.Float64Var(&c.MinDoP, "min-dop", 0, "override GPS Dilution of Precision filter")
	fs.Float64Var(&c.Start.Latitude, "latitude", 0, "override start latitude")
	fs.Float64Var(&c.Start.Longitude, "longitude", 0, "override start longitude")
	fs.Float64Var(&c.Start.Bearing, "bearing", 0, "override start bearing")
	fs.Float64Var(&c.Start.Distance, "distance", 0, "override start distance")
	annotate(fs, "gopro.render")

	goproCmd.AddCommand(cmd)
}

func init() { // nolint: gochecknoinits
	addGoproRender()
}
