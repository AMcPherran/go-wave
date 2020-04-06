package gowave

import "sync"

// Concurrency-safe struct for managing state of the Wave
type WaveState struct {
	sensorData SensorData
	motionData MotionData
	Buttons    ButtonState
	display    DisplayState
	battery    BatteryStatus
	mux        sync.Mutex
}

// SensorData setter and getter

func (ws *WaveState) SetSensorData(sd SensorData) {
	ws.mux.Lock()
	ws.sensorData = sd
	ws.mux.Unlock()
}

func (ws *WaveState) GetSensorData() SensorData {
	ws.mux.Lock()
	defer ws.mux.Unlock()
	return ws.sensorData
}

// MotionData setter and getter

func (ws *WaveState) SetMotionData(md MotionData) {
	ws.mux.Lock()
	ws.motionData = md
	ws.mux.Unlock()
}

func (ws *WaveState) GetMotionData() MotionData {
	ws.mux.Lock()
	defer ws.mux.Unlock()
	return ws.motionData
}

// BatteryStatus setter and getter
func (ws *WaveState) SetBatteryStatus(bs BatteryStatus) {
	ws.mux.Lock()
	ws.battery = bs
	ws.mux.Unlock()
}

func (ws *WaveState) GetBatteryStatus() BatteryStatus {
	ws.mux.Lock()
	defer ws.mux.Unlock()
	return ws.battery
}

// DisplayState setter and getter
func (ws *WaveState) SetDisplayState(ds DisplayState) {
	ws.mux.Lock()
	ws.display = ds
	ws.mux.Unlock()
}

func (ws *WaveState) GetDisplayState() DisplayState {
	ws.mux.Lock()
	defer ws.mux.Unlock()
	return ws.display
}

// Concurrency-safe button handling

type ButtonState struct {
	top    ButtonEvent
	middle ButtonEvent
	bottom ButtonEvent
	mux    sync.Mutex
}

func (bs *ButtonState) Top() ButtonEvent {
	bs.mux.Lock()
	defer bs.mux.Unlock()
	return bs.top
}

func (bs *ButtonState) Middle() ButtonEvent {
	bs.mux.Lock()
	defer bs.mux.Unlock()
	return bs.middle
}

func (bs *ButtonState) Bottom() ButtonEvent {
	bs.mux.Lock()
	defer bs.mux.Unlock()
	return bs.bottom
}

func (bs *ButtonState) Set(be ButtonEvent) {
	bs.mux.Lock()
	switch be.ID {
	case "A":
		bs.top = be
	case "B":
		bs.middle = be
	case "C":
		bs.bottom = be
	}
	bs.mux.Unlock()
}
