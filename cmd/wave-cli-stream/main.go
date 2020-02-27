package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	gowave "github.com/AMcPherran/go-wave"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

const waveName = "Wave"
const scanDuration = 30
const apiServiceUUID = "f3402bdcd01711e9bb652a2ae2dbcce4"
const apiCharacteristicUUID = "f3402ea2d01711e9bb652a2ae2dbcce4"

var (
	sub = flag.Duration("sub", 500*time.Second, "subscribe to notification and indication for a specified period")
	sd  = flag.Duration("sd", 15*time.Second, "scanning duration, 0 for indefinitely")
)

func main() {
	flag.Parse()

	d, err := dev.NewDevice("WaveClient")
	if err != nil {
		log.Fatalf("Can't instantiate BLE device : %s", err)
	}
	ble.SetDefaultDevice(d)

	// Search for devices named "Wave"
	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.LocalName()) == strings.ToUpper(waveName)
	}

	// Scan for devices
	fmt.Printf("Scanning for %d...\n", scanDuration)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *sd))
	client, err := ble.Connect(ctx, filter)
	if err != nil {
		log.Fatalf("Can't connect to Wave: %s", err)
	}
	fmt.Printf("Connected to Wave %s\n", client.Addr())

	// Set the MTU to 128 (Max used by Wave)
	_, err = client.ExchangeMTU(128)
	if err != nil {
		log.Fatalf("Failed to set MTU to 128: %s", err)
	}

	// Make sure we had the chance to print out the message.
	done := make(chan struct{})
	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func() {
		<-client.Disconnected()
		fmt.Printf("[ %s ] is disconnected \n", client.Addr())
		close(done)
	}()

	fmt.Printf("Discovering profile...\n")
	profile, err := client.DiscoverProfile(true)
	if err != nil {
		log.Fatalf("can't discover profile: %s", err)
	}

	// Get the BLE Service for interacting with the Wave API
	service := getService(client, profile)
	if &service == nil {
		log.Fatalf("Couldn't identify the Wave API service")
	}

	// Get the Wave API Characteristic
	var waveC *ble.Characteristic
	for _, c := range service.Characteristics {
		if c.UUID.Equal(ble.MustParse(apiCharacteristicUUID)) {
			waveC = c
		}
	}
	if &waveC == nil {
		log.Fatalf("Couldn't identify the Wave API Characteristic")
	}

	// Subscribe to incoming data

	if err := client.Subscribe(waveC, false, handleNotifications); err != nil {
		log.Fatalf("subscribe failed: %s", err)
	}
	time.Sleep(*sub)
	if err := client.Unsubscribe(waveC, false); err != nil {
		log.Fatalf("unsubscribe failed: %s", err)
	}

	// Disconnect the connection. (On OS X, this might take a while.)
	fmt.Printf("Disconnecting [ %s ]... \n", client.Addr())
	client.CancelConnection()

	<-done
}

func getService(cln ble.Client, p *ble.Profile) *ble.Service {
	service := p.FindService(&ble.Service{
		UUID: ble.MustParse(apiServiceUUID),
	})

	return service
}

func handleNotifications(data []byte) {
	// Parse the incoming data into a Query
	q, _ := gowave.NewQuery(data)
	// Handle the query
	switch q.ID {
	case "ButtonEvent":
		buttonEvent, _ := gowave.NewButtonEvent(q)
		handleButtonEvent(buttonEvent)
	case "Datastream":
		dataStream, _ := gowave.NewDatastream(q)
		handleDatastream(dataStream)
	case "BatteryStatus":
		fmt.Println(q.Payload)
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
	}
}

func handleDatastream(ds gowave.Datastream) {
	//fmt.Println(q.Payload)
}

func handleButtonEvent(be gowave.ButtonEvent) {
	fmt.Println(be)
}
