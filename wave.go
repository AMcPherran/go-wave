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

var device ble.Device

type Wave struct {
	BLE   BLE
	State WaveState
}

func Connect() (*Wave, error) {
	d, err := dev.NewDevice("WaveCli")
	if err != nil {
		return nil, err
	}
	ble.SetDefaultDevice(d)
	device = d

	// Search for devices named "Wave"
	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.LocalName()) == strings.ToUpper(WaveName)
	}

	// Scan for devices
	fmt.Printf("Scanning for %d seconds...\n", ScanDuration)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), ScanDuration*time.Second))
	client, err := ble.Connect(ctx, filter)
	if err != nil {
		device.Stop()
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

func (w *Wave) Disconnect() error {
	_ = w.BLE.Client.ClearSubscriptions()
	err := w.BLE.Client.CancelConnection()
	device.Stop()
	return err
}

func (w *Wave) HandleNotifications() error {
	err := w.Subscribe(w.defaultNotificationHandler)
	return err
}

func (w *Wave) Subscribe(handler ble.NotificationHandler) error {
	err := w.BLE.Client.Subscribe(w.BLE.Characteristic, false, handler)
	return err
}

func (w *Wave) Unsubscribe() error {
	err := w.BLE.Client.Unsubscribe(w.BLE.Characteristic, false)
	return err
}

func (w *Wave) SendQuery(q Query) error {
	b := q.ToBytes()
	err := w.BLE.Client.WriteCharacteristic(w.BLE.Characteristic, b, true)
	return err
}

// Default handler for Notifications, updates w.WaveState
func (w *Wave) defaultNotificationHandler(data []byte) {
	// Parse the incoming data into a Query
	q, _ := NewQuery(data)
	// Handle the query
	switch q.ID {
	case "ButtonEvent":
		buttonEvent, _ := NewButtonEvent(q)
		w.State.Buttons.Set(buttonEvent)
	case "Datastream":
		dataStream, _ := NewDatastream(q)
		w.State.SetMotionData(dataStream.MotionData)
		w.State.SetSensorData(dataStream.Data)
	case "BatteryStatus":
		batteryStatus, _ := NewBatteryStatus(q)
		w.State.SetBatteryStatus(batteryStatus)
	case "DeviceInfo":
		fmt.Println(q.Payload)
	case "DeviceMode":
		fmt.Println(q.Payload)
	case "Identify":
		fmt.Println(q.Payload)
	case "Recenter":
		fmt.Println(q.Payload)
	case "DisplayFrame":
		fmt.Println(q.Payload)
	case "MAX_VAL":
		fmt.Println(q.Payload)
	default:
		fmt.Println(q.Payload)
	}
}
