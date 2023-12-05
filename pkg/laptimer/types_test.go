package laptimer

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestThresholdXML(t *testing.T) {
	v1 := Threshold(10)
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Threshold>10%</Threshold>", string(data))

	var v2 Threshold
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestGearXML(t *testing.T) {
	v1 := Gear{Number: 2, Ratio: 3.980000}
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Gear>2,3.980000</Gear>", string(data))

	var v2 Gear
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestTyreXML(t *testing.T) {
	v1 := Tyre{Width: 245, Profile: 35, Size: 19, SpeedRating: "R"}
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Tyre>245 / 35 R 19</Tyre>", string(data))

	var v2 Tyre
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestTagsXML(t *testing.T) {
	v1 := Tags{"Me", "Them"}
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Tags>Me,Them</Tags>", string(data))

	var v2 Tags
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestCoordinateXML(t *testing.T) {
	v1 := Coordinate{Latitude: 50.85793657, Longitude: -0.75273358}
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Coordinate>50.85793657,-0.75273358</Coordinate>", string(data))

	var v2 Coordinate
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestAltitudeCoordinateXML(t *testing.T) {
	v1 := AltitudeCoordinate{Coordinate: Coordinate{Latitude: 50.85793657, Longitude: -0.75273358}, Altitude: 28.2}
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<AltitudeCoordinate>50.85793657,-0.75273358,28.2</AltitudeCoordinate>", string(data))

	var v2 AltitudeCoordinate
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestPositioningXML(t *testing.T) {
	v1 := Positioning{DifferentialStatus: DifferentialStatusDGPS, PositionFixing: PositionFixing3D, Interpolated: true}
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Positioning>2,2,1</Positioning>", string(data))

	var v2 Positioning
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestRelativeToStartXML(t *testing.T) {
	v1 := RelativeToStart{Distance: 32.1, Offset: Duration(time.Minute + time.Second + 210*time.Millisecond)}
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<RelativeToStart>32.1,01:01.21</RelativeToStart>", string(data))

	var v2 RelativeToStart
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestLapDateXML(t *testing.T) {
	v1 := LapDate(time.Date(2022, time.May, 31, 8, 1, 33, 0, time.UTC))
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<LapDate>31-MAY-22,08:01:33</LapDate>", string(data))

	var v2 LapDate
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestFixDateXML(t *testing.T) {
	v1 := FixDate(time.Date(2022, time.May, 31, 8, 1, 33, 800000000, time.UTC))
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<FixDate>31-MAY-22,08:01:33.80</FixDate>", string(data))

	var v2 FixDate
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestDurationXML(t *testing.T) {
	v1 := Duration(time.Minute + time.Second + 210*time.Millisecond)
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Duration>01:01.21</Duration>", string(data))

	var v2 Duration
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestFloat0dpXML(t *testing.T) {
	v1 := Float0dp(123.0)
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Float0dp>123</Float0dp>", string(data))

	var v2 Float0dp
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2) //nolint: testifylint
}

func TestFloat1dpXML(t *testing.T) {
	v1 := Float1dp(123.3)
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Float1dp>123.3</Float1dp>", string(data))

	var v2 Float1dp
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2) //nolint: testifylint
}

func TestFloat2dpXML(t *testing.T) {
	v1 := Float2dp(123.3)
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<Float2dp>123.30</Float2dp>", string(data))

	var v2 Float2dp
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2) //nolint: testifylint
}

func TestSyncPointXML(t *testing.T) {
	v1 := SyncPoint(145*time.Second + 210*time.Millisecond)
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, "<SyncPoint>145.21</SyncPoint>", string(data))

	var v2 SyncPoint
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}

func TestIntermediatesXML(t *testing.T) {
	v1 := Intermediates{
		{
			Time:     Duration(27*time.Second + 900*time.Millisecond),
			Distance: 1142.1,
		},
		{
			Time:     Duration(44*time.Second + 950*time.Millisecond),
			Distance: 1862.0,
		},
	}
	data, err := xml.Marshal(v1)
	require.NoError(t, err)
	require.Equal(t, `<Intermediates>&#xA;&#x9;&#x9;&#x9;00:27.90,1142.1&#xA;&#x9;&#x9;&#x9;00:44.95,1862.0&#xA;&#x9;&#x9;</Intermediates>`, string(data)) //nolint: lll

	var v2 Intermediates
	err = xml.Unmarshal(data, &v2)
	require.NoError(t, err)
	require.Equal(t, v1, v2)
}
