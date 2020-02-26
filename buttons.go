package gowave

import (
	"time"

	"golang.org/x/xerrors"
)

const ButtonEventID = 2

type ButtonEvent struct {
	ID        string
	Action    string
	Timestamp int64
}

func NewButtonEvent(data []byte) (ButtonEvent, error) {
	var be ButtonEvent
	if data[0] != ButtonEventID {
		return be, xerrors.Errorf("Input was not a button event")
	}
	be = ButtonEvent{
		ID:        ButtonIDs[data[4]],
		Action:    ButtonActions[data[5]],
		Timestamp: time.Now().Unix(),
	}
	return be, nil
}

var ButtonIDs = map[uint8]string{
	0: "A",
	1: "B",
	2: "C",
	3: "D",
}

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
