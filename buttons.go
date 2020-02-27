package gowave

import (
	"time"

	"golang.org/x/xerrors"
)

const ButtonEventID uint8 = 4

type ButtonEvent struct {
	ID        string
	Action    string
	Timestamp int64
}

func NewButtonEvent(q Query) (ButtonEvent, error) {
	var be ButtonEvent
	if q.ID != "ButtonEvent" {
		return be, xerrors.Errorf("Given Query was not a button event")
	}
	be = ButtonEvent{
		ID:        ButtonIDs[q.Payload[0]],
		Action:    ButtonActions[q.Payload[1]],
		Timestamp: time.Now().Unix(),
	}
	return be, nil
}

//
var ButtonIDs = map[uint8]string{
	0: "A",
	1: "B",
	2: "C",
	3: "D",
}

//
var ButtonActions = map[uint8]string{
	0: "Up",
	1: "Down",
	2: "Long",
	3: "LongUp",
	4: "ExtraLong",
	5: "ExtraLongUp",
	6: "Click",
	7: "DoubleClick",
}
