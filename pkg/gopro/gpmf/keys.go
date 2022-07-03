package gpmf

// Known keys.
// Additional information can be found here:
// https://exiftool.org/TagNames/GoPro.html
var (
	// KeyDevice unique device source for metadata.
	KeyDevice = "DEVC"

	// KeyDeviceID device/track ID.
	KeyDeviceID = "DVID"

	// KeyDeviceName device name.
	KeyDeviceName = "DVNM"

	// KeyStream nested signal stream of metadata / telemetry.
	KeyStream = "STRM"

	// KeyStreamName stream name.
	KeyStreamName = "STNM"

	// KeyRemark stream comment (debugging).
	KeyRemark = "RMRK"

	// KeyScale scaling factor (divisor).
	KeyScale = "SCAL"

	// KeyStandardUnits Standard Units (like SI).
	KeyStandardUnits = "SIUN"

	// KeyDisplayUnits display units.
	KeyDisplayUnits = "UNIT"

	// KeyTypeDef is a type definition for complex structures.
	KeyTypeDef = "TYPE"

	// KeyTimeOffset time offset of the metadata stream that follows (single 4 byte float).
	KeyTimeOffset = "TIMO"

	// KeyEmpty payloads that are empty since the device start
	// e.g. BLE disconnect.
	KeyEmpty = "EMPT"

	// KeyShutter shutter exposure times.
	KeyShutter = "SHUT"

	// KeyAccel IMU data accelerometer.
	KeyAccel = "ACCL"

	// KeyGyro IMU data gyroscope.
	KeyGyro = "GYRO"

	// KeyGPS GPS location.
	KeyGPS = "GPS5"

	// KeyImageSensorGain Image sensor gain.
	KeyImageSensorGain = "ISOG"

	// KeyGPSTime UTC time and data for GPS.
	KeyGPSTime = "GPSU"

	// KeyGPSFix GPS fix.
	KeyGPSFix = "GPSF"

	// KeyGPSDoP GPS Precision - Dilution of Precision (DOP x 100).
	KeyGSPDoP = "GPSP"

	// KeyTimeStamp time stamp for the first sample in microseconds timestamps.
	KeyTimeStamp = "STMP"

	// KeyMagnetometer magnetometer.
	KeyMagnetometer = "MAGN"

	// KeyFace face detection bounding boxes.
	KeyFace = "FACE"

	// KeyFaces faces counted per frame.
	KeyFaces = "FCNM"

	// KeySensorISO sensor ISO.
	KeySensorISO = "ISOE"

	// KeyAutoLowLight Auto Low Light frame.
	KeyAutoLowLight = "ALLD"

	// KeyWhiteBalance white balance in kelvin.
	KeyWhiteBalance = "WBAL"

	// WRGB white balance RGB gains.
	KeyWhiteBalanceRGB = "WRGB"

	// KeyFrameLuma luma (Y) average over frame.
	KeyFrameLuma = "YAVG"

	// KeyFrameHues predominant hues over the frame.
	KeyFrameHues = "HUES"

	// KeyImageUniformity image uniformity.
	KeyImageUniformity = "UNIF"

	// KeySceneClassifier scene classifier in probabilities.
	KeySceneClassifier = "SCEN"

	// KeySensourReadOut sensor read out time.
	KeySensorReadOut = "SROT"

	// KeyCameraOrientation camera orientation.
	KeyCameraOrientation = "CORI"

	// IORI Image orientation.
	KeyImageOrientation = "IORI"

	// KeyGravityVector gravity vector.
	KeyGavityVector = "GRAV"

	// KeyWindProcessing wind processing.
	KeyWindProcessing = "WNFM"

	// KeyMicrophoneWet microphone is wet.
	KeyMicrophoneWet = "MWET"

	// KeyDisparityTrack disparity track (360 modes).
	KeyDisparityTrack = "DISP"

	// KeyMainVideoFrameSkip main video frame skip.
	KeyMainVideoFrameSkip = "MSKP"

	// KeyLowResVideoFrameSkip Low res video frame skip.
	KeyLowResVideoFrameSkip = "LSKP"

	// KeyBeginTimingData beginning of data timing (arrival) in milliseconds.
	KeyBeginTimingData = "TICK"

	// KeyEndTimingData end of data timing (arrival)  in milliseconds.
	KeyEndTimingData = "TOCK"

	// KeyTotalSamples total sample count including the current payload.
	KeyTotalSamples = "TSMP"

	// KeyDeviceTemperature device temperature in Celsius.
	KeyDeviceTemperature = "TMPC"

	// KeyQuantize quantize used to enable stream compression.
	// 1 - enable.
	// 2+ enable and quantize by this value.
	KeyQuantize = "QUAN"

	// KeyVersion version of the metdata stream (debugging).
	KeyVersion = "VERS"

	// KeyFree n bytes reserved for more metadata added to an existing stream.
	KeyFree = "FREE"

	// KeyOrientationIn - input 'n' channel data orientation, lowercase is negative
	// e.g. "Zxy" or "ABGR".
	KeyOrientationIn = "ORIN"

	// KeyOrientationOut - output 'n' channel data orientation, e.g. "XYZ" or "RGBA".
	KeyOrientationOut = "ORIO"

	// KeyMatrix - 2D matrix for any sensor calibration.
	KeyMatrix = "MTRX"

	// KeyPreformatted - GPMF data.
	KeyPreformatted = "PFRM"

	// KeyTimeStamps stream of all the timestamps delivered.
	// Generally don't use this. This would be if your sensor has no periodic times,
	// yet precision is required, or for debugging.
	KeyTimeStamps = "STPS"
)

type parserFunc func(*Element) error

// We need two hashes to avoid an initialisation loop.
var (
	// keyParsers has nil entries so we can check if know about
	// a given key.
	keyParsers = map[string]parserFunc{
		KeyDevice:               nil,
		KeyDeviceID:             parseMetadata,
		KeyDeviceName:           parseMetadata,
		KeyStream:               nil,
		KeyStreamName:           parseMetadata,
		KeyRemark:               nil,
		KeyScale:                parseScale,
		KeyStandardUnits:        parseMetadata,
		KeyDisplayUnits:         parseMetadata,
		KeyTypeDef:              parseMetadata,
		KeyTimeOffset:           nil,
		KeyEmpty:                nil,
		KeyShutter:              nil,
		KeyAccel:                parseAccel,
		KeyGyro:                 parseGyro,
		KeyGPS:                  parseGPS,
		KeyImageSensorGain:      nil,
		KeyGPSTime:              parseMetadata,
		KeyGPSFix:               parseGPSFix,
		KeyGSPDoP:               parseGPSDoP,
		KeyTimeStamp:            nil,
		KeyMagnetometer:         parseMagnetometer,
		KeyFace:                 parseFace,
		KeyFaces:                parseHasMetadata,
		KeySensorISO:            parseHasMetadata,
		KeyAutoLowLight:         nil,
		KeyWhiteBalance:         nil,
		KeyWhiteBalanceRGB:      parseWhiteBalanceRGB,
		KeyFrameLuma:            nil,
		KeyFrameHues:            nil,
		KeyImageUniformity:      nil,
		KeySceneClassifier:      nil,
		KeySensorReadOut:        nil,
		KeyCameraOrientation:    nil,
		KeyImageOrientation:     nil,
		KeyGavityVector:         nil,
		KeyWindProcessing:       nil,
		KeyMicrophoneWet:        nil,
		KeyDisparityTrack:       nil,
		KeyMainVideoFrameSkip:   nil,
		KeyLowResVideoFrameSkip: nil,
		KeyBeginTimingData:      nil,
		KeyEndTimingData:        nil,
		KeyTotalSamples:         parseMetadata,
		KeyDeviceTemperature:    parseMetadata,
		KeyQuantize:             nil,
		KeyVersion:              nil,
		KeyFree:                 nil,
		KeyOrientationIn:        nil,
		KeyOrientationOut:       nil,
		KeyMatrix:               nil,
		KeyPreformatted:         nil,
		KeyTimeStamps:           nil,
	}

	keyNames = map[string]string{
		KeyDeviceID:          "device_id",
		KeyDeviceName:        "device_name",
		KeyStreamName:        "stream_name",
		KeyScale:             "scale",
		KeyStandardUnits:     "standard_units",
		KeyDisplayUnits:      "display_units",
		KeyTypeDef:           "type_def",
		KeyAccel:             "acceleration",
		KeyGyro:              "gyroscope",
		KeyGPS:               "gps",
		KeyGPSTime:           "gps_time",
		KeyGPSFix:            "gps_fix",
		KeyGSPDoP:            "gps_dilution_of_precision",
		KeyMagnetometer:      "magnetometer",
		KeyFace:              "face_detection",
		KeyFaces:             "faces",
		KeyWhiteBalanceRGB:   "white_balance_rgb",
		KeyTotalSamples:      "samples",
		KeyDeviceTemperature: "device_temperature",
	}
)

func friendlyName(key string) string {
	if f := keyNames[key]; f != "" {
		return f
	}

	return key
}
