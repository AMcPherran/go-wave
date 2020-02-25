package gowave

type ButtonEvent struct {
	ID        uint8
	Action    ButtonAction
	Timestamp float64
}

type ButtonAction uint8

const (
	Up          ButtonAction = 0
	Down        ButtonAction = 1
	Long        ButtonAction = 2
	LongUp      ButtonAction = 3
	ExtraLong   ButtonAction = 4
	ExtraLongUp ButtonAction = 5
	Click       ButtonAction = 6
	DoubleClick ButtonAction = 7
)
