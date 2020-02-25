package go-wave


type Datastream struct {
	Data       SensorData,
	MotionData MotionData
}

type SensorData struct {
	Gyro  []float64,
	Accel []float64,
	Mag   []float64,
}

type MotionData struct {
	RawPos     []float64,
	CurrentPos []float64,
	Euler       []float64,
	Timestamp   float64,
	Tap         Peak,
}

type Peak struct {
	Detected     bool,
	NormVelocity float64,
}