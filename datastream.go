package gowave

import (
	"golang.org/x/xerrors"
)

const DatastreamID = 1

type Datastream struct {
	Data       SensorData `json:"data"`
	MotionData MotionData `json:"motionData"`
}

func NewDatastream(q Query) (Datastream, error) {
	var ds Datastream
	if q.ID != "Datastream" {
		return ds, xerrors.Errorf("Given Query was not a Datastream")
	}

	sd := SensorData{
		Gyro:      GetGyroVector(q.Payload),
		Accel:     GetAccelVector(q.Payload),
		Mag:       GetMagVector(q.Payload),
		Timestamp: Float32frombytes(q.Payload[36:40]), // Bytes 36,37,38,39
	}

	md := MotionData{
		Timestamp:  Float32frombytes(q.Payload[89:93]), // Bytes 89,90,91,92
		RawPos:     GetOnBoardOrientation(q.Payload),
		CurrentPos: GetCorrectedOrientation(q.Payload),
		Euler:      GetEulerOrientation(q.Payload),
		Tap: Peak{
			Detected:     q.Payload[84] > 0,
			NormVelocity: Float32frombytes(q.Payload[85:89]), // Bytes 85,86,87,88
		},
	}

	ds = Datastream{
		Data:       sd,
		MotionData: md,
	}

	return ds, nil
}

type SensorData struct {
	Gyro      Vector  `json:"gyro"`
	Accel     Vector  `json:"accel"`
	Mag       Vector  `json:"max"`
	Timestamp float32 `json:"timestamp"`
}

type MotionData struct {
	RawPos     Quaternion `json:"rawPos"`
	CurrentPos Quaternion `json:"currentPos"`
	Euler      Vector     `json:"euler"`
	Timestamp  float32    `json:"timestamp"`
	Tap        Peak       `json:"peak"`
}

type Peak struct {
	Detected     bool    `json:"detected"`
	NormVelocity float32 `json:"normVelocity"`
}

func GetGyroVector(payload []byte) Vector {
	v := Vector{
		X: Float32frombytes(payload[:4]),   // Bytes 0,1,2,3
		Y: Float32frombytes(payload[4:8]),  // Bytes 4,5,6,7
		Z: Float32frombytes(payload[8:12]), // Bytes 8,9,10,11
	}
	return v
}

func GetAccelVector(payload []byte) Vector {
	v := Vector{
		X: Float32frombytes(payload[12:16]), // Bytes 12,13,14,15
		Y: Float32frombytes(payload[16:20]), // Bytes 16,17,18,19
		Z: Float32frombytes(payload[20:24]), // Bytes 20,21,22,23
	}
	return v
}

func GetMagVector(payload []byte) Vector {
	v := Vector{
		X: Float32frombytes(payload[24:28]), // Bytes 24,25,26,27
		Y: Float32frombytes(payload[28:32]), // Bytes 28,29,30,31
		Z: Float32frombytes(payload[32:36]), // Bytes 32,33,34,35
	}
	return v
}

func GetOnBoardOrientation(payload []byte) Quaternion {
	q := Quaternion{
		X: Float32frombytes(payload[40:44]), // Bytes 40,41,42,43
		Y: Float32frombytes(payload[44:48]), // Bytes 44,45,46,47
		Z: Float32frombytes(payload[48:52]), // Bytes 48,49,50,51
		W: Float32frombytes(payload[52:56]), // Bytes 52,53,54,55
	}
	return q
}

func GetCorrectedOrientation(payload []byte) Quaternion {
	q := Quaternion{
		X: Float32frombytes(payload[56:60]), // Bytes 56,57,58,59
		Y: Float32frombytes(payload[60:64]), // Bytes 60,61,62,63
		Z: Float32frombytes(payload[64:68]), // Bytes 64,65,66,67
		W: Float32frombytes(payload[68:72]), // Bytes 68,69,70,71
	}
	return q
}

func GetEulerOrientation(payload []byte) Vector {
	v := Vector{
		X: Float32frombytes(payload[72:76]), // Bytes 72,73,74,75
		Y: Float32frombytes(payload[76:80]), // Bytes 76,77,78,79
		Z: Float32frombytes(payload[80:84]), // Bytes 80,81,82,83
	}
	return v
}
