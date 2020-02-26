package gowave

type Query struct {
	ID          string
	Type        string
	PayloadSize int
	Payload     []byte
}

func NewQuery(data []byte) Query {
	payload := data[2:]
	q := Query{
		ID:          QueryIDs[data[1]],
		Type:        QueryTypes[data[0]],
		PayloadSize: len(payload),
		Payload:     payload,
	}
	return q
}

var QueryIDs = map[uint8]string{
	0:             "Unknown",
	DatastreamID:  "Datastream",
	2:             "BatteryStatus",
	3:             "DeviceInfo",
	ButtonEventID: "ButtonEvent",
	5:             "DeviceMode",
	6:             "Identify",
	7:             "Recenter",
	8:             "DisplayFrame",
	9:             "MAX_VAL",
}

var QueryTypes = map[uint8]string{
	0: "Unknown",
	1: "Request",
	2: "Response",
	3: "Stream",
	4: "MAX_VAL",
}
