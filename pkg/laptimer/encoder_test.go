package laptimer

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	dt := time.Date(2009, time.November, 10, 23, 12, 16, 810000000, time.UTC)
	db := NewDB()
	db.Laps = append(db.Laps, Lap{
		ID:      1,
		Date:    LapDate(dt),
		LapTime: Duration(3*time.Minute + 12*time.Second),
		Vehicle: "testVehicle",
		Track:   "testTrack",
		Intermediates: Intermediates{
			{
				Time: Duration(
					2*time.Minute +
						13*time.Second +
						800*time.Millisecond,
				),
				Distance: 1022.4,
			},
		},
		LapRecordingType: LapRecordingTriggered,
		OverallDistance:  123.4,
		WeatherCondition: 2, // TODO: use enum.
		AmbientTemp:      15.2,
		AmbientPressure:  1012.1, // This is in hPa should it be?
		RelativeHumidity: 40.1,   // This is a percentage should it be?
		Videos: []Video{
			{
				Overlaid:  true,
				URL:       "http://example.com/test",
				SyncPoint: SyncPoint(23*time.Second + 810*time.Millisecond),
			},
		},
		Note: `First line
Second line`,
		Tags: Tags{"Me", "Other"},
		Recording: Recording{
			Fixes: []Fix{
				{
					ID:   1,
					Date: FixDate(dt),
					Coordinate: AltitudeCoordinate{
						Coordinate: Coordinate{
							Latitude:  50.85793657,
							Longitude: -0.75273358,
						},
						Altitude: 28.2,
					},
					Speed: 54.3,
					Positioning: Positioning{
						DifferentialStatus: DifferentialStatusDGPS,
						PositionFixing:     PositionFixing3D,
						Interpolated:       true,
					},
					Satellites: 3,
					Direction:  166.1,
					Hdop:       10.12,
					Accuracy:   3.5,
					RelativeToStart: RelativeToStart{
						// These aren't valid for initial entry but are better
						// testing.
						Distance: 10.0,
						Offset: Duration(
							3*time.Minute +
								25*time.Second +
								340*time.Millisecond,
						),
					},
					Acceleration: &Acceleration{
						Source:  1,
						Lateral: 0.21,
						Lineal:  1.32,
						Coordinate: Coordinate{
							Latitude:  50.85793657,
							Longitude: -0.75273358,
						},
					},
					OBD: &OBD{
						EngineRPM:                intp(2738),
						MassAirFlow:              float2dp(23.32),
						VehicleSpeed:             float1dp(86.2),
						Throttle:                 float2dp(34.21),
						FuelLevel:                float2dp(87.12),
						CoolantTemp:              float1dp(63.2),
						OilTemp:                  float1dp(76.2),
						ManifoldAbsolutePressure: float2dp(200.12),
						IntakeAirTemperature:     float0dp(20.1),
						Gear:                     intp(2),
						OilPressure:              float64p(220.3),
						BrakePressure:            float64p(22.1),
						WheelSpeedRearLeft:       float64p(23.52),
						WheelSpeedRearRight:      float64p(21.32),
						WheelSpeedFrontLeft:      float64p(22.12),
						WheelSpeedFrontRight:     float64p(21.56),
						YawRate:                  float64p(12.34),
						Odometer:                 float64p(44.2),
						SteerAngle:               float64p(14.3),
						SteerRate:                float64p(3.5),
						FixType:                  fixTypep(FixTypeCombustion),
						SupportLevel:             intp(1),
						DriverPower:              float64p(22.6),
						BatteryTemp:              float64p(35.7),
						BatteryVoltage:           float64p(15.2),
						ConsumedPower:            float64p(33.8),
					},
				},
			},
		},
	})

	var buf bytes.Buffer
	enc, err := NewEncoder(&buf)
	require.NoError(t, err)

	err = enc.Encode(db)
	require.NoError(t, err)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<LapTimerDB>
	<name>LapTimer Database</name>
	<lap index="1">
		<date>10-NOV-09,23:12:16</date>
		<lapTime>03:12.00</lapTime>
		<vehicle>testVehicle</vehicle>
		<track>testTrack</track>
		<intermediates>
			02:13.80,1022.4
		</intermediates>
		<lapRecordingType>2</lapRecordingType>
		<overallDistance>123.4</overallDistance>
		<weatherCondition>2</weatherCondition>
		<ambientTemp>15.2</ambientTemp>
		<ambientPressure>1012</ambientPressure>
		<relativeHumidity>40.10</relativeHumidity>
		<video>
			<overlaid>true</overlaid>
			<url>http://example.com/test</url>
			<syncPoint>23.81</syncPoint>
		</video>
		<note>First line
Second line</note>
		<tags>Me,Other</tags>
		<recording>
			<fix index="1">
				<date>10-NOV-09,23:12:16.81</date>
				<coordinate>50.85793657,-0.75273358,28.2</coordinate>
				<speed>54.3</speed>
				<positioning>2,2,1</positioning>
				<satellites>3</satellites>
				<direction>166.1</direction>
				<hdop>10.12</hdop>
				<accuracy>3.5</accuracy>
				<relativeToStart>10.0,03:25.34</relativeToStart>
				<acceleration>
					<source>1</source>
					<lateral>0.21</lateral>
					<lineal>1.32</lineal>
					<coordinate>50.85793657,-0.75273358</coordinate>
				</acceleration>
				<obd>
					<rpm>2738</rpm>
					<maf>23.32</maf>
					<speed>86.2</speed>
					<throttle>34.21</throttle>
					<fuelLevel>87.12</fuelLevel>
					<coolant>63.2</coolant>
					<oil>76.2</oil>
					<map>200.12</map>
					<iat>20</iat>
					<gear>2</gear>
					<oilp>220.3</oilp>
					<brake>22.1</brake>
					<speedrl>23.52</speedrl>
					<speedrr>21.32</speedrr>
					<speedfl>22.12</speedfl>
					<speedfr>21.56</speedfr>
					<yaw>12.34</yaw>
					<odometer>44.2</odometer>
					<steerangle>14.3</steerangle>
					<steerrate>3.5</steerrate>
					<fixtype>0</fixtype>
					<supportlevel>1</supportlevel>
					<driverpower>22.6</driverpower>
					<batterytemp>35.7</batterytemp>
					<voltage>15.2</voltage>
					<consumedpower>33.8</consumedpower>
				</obd>
			</fix>
		</recording>
	</lap>
</LapTimerDB>`
	require.Equal(t, expected, buf.String())
}

func intp(v int) *int {
	return &v
}

func float64p(v float64) *float64 {
	return &v
}

func float0dp(v float64) *Float0dp {
	f := Float0dp(v)
	return &f
}

func float1dp(v float64) *Float1dp {
	f := Float1dp(v)
	return &f
}

func float2dp(v float64) *Float2dp {
	f := Float2dp(v)
	return &f
}

func fixTypep(v FixType) *FixType {
	return &v
}
