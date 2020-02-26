package gowave

const DatastreamID = 1

type Datastream struct {
	Data       SensorData `json:"data"`
	MotionData MotionData `json:"motionData"`
}

type SensorData struct {
	Gyro  []float64 `json:"gyro"`
	Accel []float64 `json:"accel"`
	Mag   []float64 `json:"max"`
}

type MotionData struct {
	RawPos     []float64 `json:"rawPos"`
	CurrentPos []float64 `json:"currentPos"`
	Euler      []float64 `json:"euler"`
	Timestamp  float64   `json:"timestamp"`
	Tap        Peak      `json:"peak"`
}

type Peak struct {
	Detected     bool    `json:"detected"`
	NormVelocity float64 `json:"normVelocity"`
}
