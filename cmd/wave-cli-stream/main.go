package main

import (
	"flag"
	"fmt"
	"log"

	gowave "github.com/AMcPherran/go-wave"
)

var buttonStates = map[string]bool{
	"A": false,
	"B": false,
	"C": false,
}

func main() {
	flag.Parse()

	wave, err := gowave.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Wave: %s", err)
	}

	go handleInbound(wave)

	<-wave.BLE.Client.Disconnected()
	// Disconnect the connection
	wave.Disconnect()
}

func handleInbound(wave *gowave.Wave) {
	// Subscribe to incoming data
	if err := wave.Subscribe(handleNotifications); err != nil {
		log.Fatalf("subscribe failed: %s", err)
	}
	log.Printf("Receiving incoming data from Wave")
	<-wave.BLE.Client.Disconnected()
	log.Printf("Wave is disconnecting")
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
	default:
		fmt.Println(q.Payload)
	}
}

func handleDatastream(ds gowave.Datastream) {
	// Print out the Euler vector if the B button is held down
	if buttonStates["B"] {
		fmt.Println(ds.MotionData.Euler)
	}
	if buttonStates["A"] {
		fmt.Println(ds.MotionData.CurrentPos)
	}
	if buttonStates["C"] {
		fmt.Println(ds.Data.Accel)
	}
}

func handleButtonEvent(be gowave.ButtonEvent) {
	fmt.Println(be)
	if be.Action == "Up" || be.Action == "ExtraLongUp" || be.Action == "LongUp" {
		buttonStates[be.ID] = false
	} else if be.Action == "Down" {
		buttonStates[be.ID] = true
	}
}
