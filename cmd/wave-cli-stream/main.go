package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	gowave "github.com/AMcPherran/go-wave"
)

func main() {
	flag.Parse()

	wave, err := gowave.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Wave: %s", err)
	}
	// Subscribe to incoming data
	if err := wave.Subscribe(handleNotifications); err != nil {
		log.Fatalf("subscribe failed: %s", err)
	}
	time.Sleep(500 * time.Second)
	// Unsubscribe
	if err := wave.Unsubscribe(); err != nil {
		log.Fatalf("unsubscribe failed: %s", err)
	}

	// Disconnect the connection
	wave.Disconnect()
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
	fmt.Println(ds)
}

func handleButtonEvent(be gowave.ButtonEvent) {
	fmt.Println(be)
}
