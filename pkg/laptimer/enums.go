package laptimer

// LapRecordingType indicates the type of recording.
type LapRecordingType int

const (
	// LapRecordingUnknown is used for old recordings with unknown
	// state or as null value.
	LapRecordingUnknown LapRecordingType = iota

	// LapRecordingManual indicates user triggered manual recording.
	LapRecordingManual

	// LapRecordingTriggered indicates fully triggered recording with
	// no information on status of track set used.
	LapRecordingTriggered

	// LapRecordingIncomplete is used for recordings with one end
	// triggered and the other manual (both situations happen).
	LapRecordingIncomplete
)

// DifferentialStatus indicates the status of differential GPS used for a Fix.
type DifferentialStatus int

const (
	// DifferentialStatusUnknown indicates the status is unknown.
	DifferentialStatusUnknown DifferentialStatus = iota

	// DifferentialStatus2D3D indicates either 2D or 3D is used typically
	// depending on the number of satellites 2D with three or fewer and 3D
	// with four or more.
	DifferentialStatus2D3D

	// DifferentialStatusDGPS indicates differential GPS is used.
	DifferentialStatusDGPS

	// DifferentialStatusInvalid indicates that status is invalid.
	DifferentialStatusInvalid
)

// PositionFixing indicates the type of position fixing used for a Fix.
type PositionFixing int

const (
	// PositionFixingNoFix indicates an invalid fix.
	PositionFixingNoFix PositionFixing = iota

	// PositionFixing2D either a fix measured with a too small number of
	// satellites, or a triangulated fix.
	PositionFixing2D

	// PositionFixing3D indicates a standard measured fix.
	PositionFixing3D

	// PositionFixingVirtual2D3D unknown quality of a fix, mostly used
	// to signal a virtual (generated) fix.
	PositionFixingVirtual2D3D

	// PositionFixing2DIndoor position delivered by some (non-GPS)
	// indoor positioning system.
	PositionFixing2DIndoor
)

// DriveWheels represents the wheels which drive the vehicle.
type DriveWheels string

const (
	// FrontWheelDrive is used for vehicles which use front wheel drive only.
	FrontWheelDrive DriveWheels = "front"

	// RearWheelDrive is used for vehicles which use rear wheel drive only.
	RearWheelDrive DriveWheels = "rear"

	// AllWheelDrive is used for vehicles which drive all wheels.
	AllWheelDrive DriveWheels = "all"
)

// IntakeType represents engine intake type.
type IntakeType string

const (
	// UnspecifiedIntake represents a unspecified intake type.
	UnspecifiedIntake IntakeType = ""

	// NaturallyAspirated represents a naturally aspirated engine.
	NaturallyAspirated IntakeType = "Naturally Aspirated"

	// Turbocharged represents a turbocharged engine.
	Turbocharged IntakeType = "Turbocharged"

	// Supercharged represents a supercharged engine.
	Supercharged IntakeType = "Supercharged"
)

// EngineType represents the type of the engine powering the vehicle.
type EngineType string

const (
	// UnspecifiedEngine represents a unspecified engine type.
	UnspecifiedEngine EngineType = ""

	// Otto represents and engine with uses an Otto cycle, also
	// know as gasoline or petrol engine.
	Otto EngineType = "Otto"

	// Diesel represents an engine which runs on diesel.
	Diesel EngineType = "Diesel"

	// Electric represents an electric motor.
	Electric EngineType = "Electric"
)

// FixType TODO(steve): document.
type FixType int

// FixType enums TODO(steve): document.
const (
	FixTypeCombustion FixType = iota
	FixTypeSpeedAndCadence
	FixTypeElectric
	FixTypeHybrid
)

// TyrePosition represents a tyre position in TPMS monitoring.
type TyrePosition int

// TyrePosistions for TPMS.
const (
	FrontLeft TyrePosition = 2 + iota
	FrontRight
	RearLeft
	RearRight
)
