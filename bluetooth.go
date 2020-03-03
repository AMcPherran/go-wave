package gowave

import "github.com/go-ble/ble"

// Struct for organizing the low-level BLE interfaces under the Wave
type BLE struct {
	Client         ble.Client
	Profile        *ble.Profile
	Characteristic *ble.Characteristic
	Disconnect     chan struct{}
}

func getService(cln ble.Client, p *ble.Profile) *ble.Service {
	service := p.FindService(&ble.Service{
		UUID: ble.MustParse(ApiServiceUUID),
	})

	return service
}
