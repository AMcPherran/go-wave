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
