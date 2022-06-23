package convert

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stevenh/tracktools/pkg/laptimer"
	"github.com/stevenh/tracktools/pkg/trackaddict"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/interp"
)

func TestTrackAddict(t *testing.T) {
	f, err := os.Open("../../test/Log-20220531-085930 Goodwood Motorcircuit - 2.57.527.csv")
	require.NoError(t, err)
	defer f.Close() // nolint: errcheck

	dec, err := trackaddict.NewDecoder(f)
	require.NoError(t, err)

	sess, err := dec.Decode()
	require.NoError(t, err)

	enc, err := laptimer.NewEncoder(ioutil.Discard)
	require.NoError(t, err)

	var pl interp.PiecewiseLinear
	conv, err := NewTrackAddict(
		TrackOpt("Goodwood"),
		TagsOpt("Me"),
		VehicleOpt("'19 McLaren 720s"),
		PredictorOpt(&pl),
	)
	require.NoError(t, err)

	db, err := conv.LapTimer(sess)
	require.NoError(t, err)

	err = enc.Encode(db)
	require.NoError(t, err)
}
