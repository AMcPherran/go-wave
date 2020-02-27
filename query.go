package gowave

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/xerrors"
)

type Query struct {
	ID          string
	Type        string
	PayloadSize uint64
	Payload     []byte
}

func NewQuery(data []byte) (Query, error) {
	var q Query
	buf := bytes.NewBuffer(data[2:3])
	payloadSize, err := binary.ReadUvarint(buf)
	if err != nil {
		return q, xerrors.Errorf("Failed to decode payload size bytes, message is likely corrupted: ", err)
	}
	payload := data[4:]
	q = Query{
		ID:          QueryIDs[data[1]],
		Type:        QueryTypes[data[0]],
		PayloadSize: payloadSize,
		Payload:     payload,
	}
	if int(payloadSize) != len(payload) {
		err := xerrors.Errorf("Length of the payload did not match expected size, expected %d bytes received %d", payloadSize, len(payload))
		return q, err
	}

	return q, nil
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
