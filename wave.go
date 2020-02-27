package gowave

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"golang.org/x/xerrors"
)

const WaveName = "Wave"
const ScanDuration = 30
const ApiServiceUUID = "f3402bdcd01711e9bb652a2ae2dbcce4"
const ApiCharacteristicUUID = "f3402ea2d01711e9bb652a2ae2dbcce4"

type Wave struct {
	BLE BLE
}

func Connect() (*Wave, error) {
	d, err := dev.NewDevice("WaveClient")
	if err != nil {
		return nil, err
	}
	ble.SetDefaultDevice(d)

	// Search for devices named "Wave"
	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.LocalName()) == strings.ToUpper(WaveName)
	}

	// Scan for devices
	fmt.Printf("Scanning for %d seconds...\n", ScanDuration)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), ScanDuration*time.Second))
	client, err := ble.Connect(ctx, filter)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Connected to Wave %s\n", client.Addr())

	// Set the MTU to 128 (Max used by Wave)
	_, err = client.ExchangeMTU(128)
	if err != nil {
		return nil, err
	}

	profile, err := client.DiscoverProfile(true)
	if err != nil {
		return nil, err
	}

	// Get the BLE Service for interacting with the Wave API
	service := getService(client, profile)
	if &service == nil {
		return nil, xerrors.Errorf("Failed to identify the Wave API Service on this device")
	}

	// Get the Wave API Characteristic
	var waveC *ble.Characteristic
	for _, c := range service.Characteristics {
		if c.UUID.Equal(ble.MustParse(ApiCharacteristicUUID)) {
			waveC = c
		}
	}
	if &waveC == nil {
		return nil, xerrors.Errorf("Failed to identify the Wave API Characterstic on this device")
	}

	wave := Wave{
		BLE: BLE{
			Client:         client,
			Characteristic: waveC,
			Profile:        profile,
		},
	}

	return &wave, nil
}

func (w Wave) Disconnect() error {
	_ = w.BLE.Client.ClearSubscriptions()
	err := w.BLE.Client.CancelConnection()
	return err
}

func (w Wave) Subscribe(handler ble.NotificationHandler) error {
	err := w.BLE.Client.Subscribe(w.BLE.Characteristic, false, handler)
	return err
}

func (w Wave) Unsubscribe() error {
	err := w.BLE.Client.Unsubscribe(w.BLE.Characteristic, false)
	return err
}

func getService(cln ble.Client, p *ble.Profile) *ble.Service {
	service := p.FindService(&ble.Service{
		UUID: ble.MustParse(ApiServiceUUID),
	})

	return service
}
