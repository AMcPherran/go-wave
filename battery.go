package gowave

import "golang.org/x/xerrors"

func GetBatteryStatusQuery() Query {
	q := Query{
		ID:          "BatteryStatus",
		Type:        "Request",
		PayloadSize: 0,
	}

	return q
}

func (w *Wave) SendBatteryStatusRequest() error {
	q := GetBatteryStatusQuery()
	err := w.SendQuery(q)
	return err
}

type BatteryStatus struct {
	Voltage    float32 `json:"voltage"`
	Percentage float32 `json:"percentage"`
	Charging   bool    `json:"charging"`
}

func NewBatteryStatus(q Query) (BatteryStatus, error) {
	var bs BatteryStatus
	if q.ID != "BatteryStatus" {
		return bs, xerrors.Errorf("Given Query was not a battery status")
	}
	bs = BatteryStatus{
		Voltage:    Float32frombytes(q.Payload[0:4]),
		Percentage: Float32frombytes(q.Payload[4:8]),
		Charging:   bool(q.Payload[8] == 1),
	}
	return bs, nil
}
