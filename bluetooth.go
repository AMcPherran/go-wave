package gowave

import "github.com/go-ble/ble"

// Struct for organizing the low-level BLE interfaces under the Wave
type BLE struct {
	Client         ble.Client
	Profile        *ble.Profile
	Characteristic *ble.Characteristic
	Disconnect     chan struct{}
}
