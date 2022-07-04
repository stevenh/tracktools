package image

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
)

const (
	// defaultWeight is the default render weight.
	defaultWeight = 2
)

var (
	// Red represents the color Red.
	Red = color.RGBA{R: 255, A: 0xff}

	// Green represents the color Green.
	Green = color.RGBA{G: 255, A: 0xff}

	// Blue represents the color Blue.
	Blue = color.RGBA{B: 255, A: 0xff}
)

// Option is a option to a Image.
type Option func(*Image)

// Render is function which writes a image img to w.
type Render func(w io.Writer, img image.Image) error

// Weight sets the for the paths rendered.
// Default: 2.
func Weight(weight float64) Option {
	return func(i *Image) {
		i.weight = weight
	}
}

// Image represents an track session which can render an image.
type Image struct {
	*sm.Context // Embedded to provide flexibility.

	weight float64
	render Render
}

// New returns a new Image.
func New(width, height int, options ...Option) *Image {
	i := &Image{
		render: png.Encode,
		weight: defaultWeight,
	}
	i.Reset(width, height)

	for _, f := range options {
		f(i)
	}

	return i
}

// Reset resets the image context.
func (i *Image) Reset(width, height int) {
	i.Context = sm.NewContext()
	i.SetSize(width, height)
}

// AddPath add a path represented by positions with color c to the image.
func (i *Image) AddPath(positions []s2.LatLng, c color.Color) {
	i.AddObject(sm.NewPath(positions, c, i.weight))
}

// Provider sets the provider title.
func (i *Image) Provider(title string) {
	prov := sm.NewTileProviderOpenStreetMaps()
	prov.Attribution = fmt.Sprintf("%s | %s", prov.Attribution, title)
	i.SetTileProvider(prov)
}

// Start adds the start line to image.
func (i *Image) Start(lat1, lon1, lat2, lon2 float64) {
	i.AddObject(sm.NewPath([]s2.LatLng{
		s2.LatLngFromDegrees(lat1, lon1),
		s2.LatLngFromDegrees(lat2, lon2),
	}, Blue, i.weight))
}

// Render renders the image to file.
func (i *Image) Render(file string) (err error) {
	img, err := i.Context.Render()
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	return i.render(f, img)
}
