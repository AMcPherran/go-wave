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

func (q Query) ToBytes() []byte {
	data := make([]byte, 2)
	data[0] = ReverseQueryIDs[q.ID]
	data[1] = ReverseQueryTypes[q.Type]
	// Payload size conversion
	ps := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(ps, int64(q.PayloadSize))
	b := ps[:n]
	//
	data = append(data, b...)
	data = append(data, q.Payload...)
	return data
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

var ReverseQueryIDs = map[string]uint8{
	"Unknown":       0,
	"Datastream":    DatastreamID,
	"BatteryStatus": 2,
	"DeviceInfo":    3,
	"ButtonEvent":   ButtonEventID,
	"DeviceMode":    5,
	"Identify":      6,
	"Recenter":      7,
	"DisplayFrame":  8,
	"MAX_VAL":       9,
}

var QueryTypes = map[uint8]string{
	0: "Unknown",
	1: "Request",
	2: "Response",
	3: "Stream",
	4: "MAX_VAL",
}

var ReverseQueryTypes = map[string]uint8{
	"Unknown":  0,
	"Request":  1,
	"Response": 2,
	"Stream":   3,
	"MAX_VAL":  4,
}
