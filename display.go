package gowave

type Pixel struct {
	X uint8
	Y uint8
}

func GetTestDisplayQuery() Query {
	// Default API frame
	frame := [][]byte{
		{000, 255, 000, 000, 255, 255, 000, 000, 255},
		{255, 000, 255, 000, 255, 000, 255, 000, 255},
		{255, 255, 255, 000, 255, 255, 000, 000, 255},
		{255, 000, 255, 000, 255, 000, 000, 000, 255},
		{000, 000, 000, 000, 000, 000, 000, 000, 000},
	}

	// Convert frame to bytes payload
	var payload []byte
	for _, row := range frame {
		payload = append(payload, row...)
	}

	q := Query{
		ID:          "DisplayFrame",
		Type:        "Request",
		PayloadSize: 45,
		Payload:     payload,
	}

	return q
}

func BlankDisplayFrame() [][]byte {
	frame := [][]byte{
		{000, 000, 000, 000, 000, 000, 000, 000, 000},
		{000, 000, 000, 000, 000, 000, 000, 000, 000},
		{000, 000, 000, 000, 000, 000, 000, 000, 000},
		{000, 000, 000, 000, 000, 000, 000, 000, 000},
		{000, 000, 000, 000, 000, 000, 000, 000, 000},
	}
	return frame
}

func GetDisplayFrameQuery(frame [][]byte) Query {
	// Convert frame to bytes payload
	var payload []byte
	for _, row := range frame {
		payload = append(payload, row...)
	}

	q := Query{
		ID:          "DisplayFrame",
		Type:        "Request",
		PayloadSize: 45,
		Payload:     payload,
	}

	return q
}

func (w *Wave) SetDisplay(frame [][]byte) error {
	q := GetDisplayFrameQuery(frame)
	err := w.SendQuery(q)
	return err
}
