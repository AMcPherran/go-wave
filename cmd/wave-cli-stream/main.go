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

	if err := wave.HandleNotifications(); err != nil {
		log.Fatalf("subscribe failed: %s", err)
	}
	log.Printf("Receiving incoming data from Wave")

	// Main loop for reading and acting on wave.State
	var lastState gowave.WaveState
	for true {
		middleButton := wave.State.Buttons.Middle()
		btnAction := middleButton.Action
		// Only update the display and recenter if the button state changed
		if middleButton != lastState.Buttons.Middle() {
			// If the Button has been released, blank the display
			if btnAction == "Up" || btnAction == "LongUp" || btnAction == "ExtraLongUp" {
				frame := gowave.BlankDisplayFrame()
				wave.SetDisplay(frame)
			}
			// If the Button was pressed down, display a dot and recenter
			if btnAction == "Down" {
				frame := [][]byte{
					{000, 000, 000, 000, 000, 000, 000, 000, 000},
					{000, 000, 000, 000, 255, 000, 000, 000, 000},
					{000, 000, 000, 255, 255, 255, 000, 000, 000},
					{000, 000, 000, 000, 255, 000, 000, 000, 000},
					{000, 000, 000, 000, 000, 000, 000, 000, 000},
				}
				wave.SetDisplay(frame)
				wave.Recenter()
			}
			lastState.Buttons.Set(middleButton)
		} else {
			// While the button is held down, print the current motion data
			if btnAction == "Down" || btnAction == "Long" || btnAction == "ExtraLong" {
				md := wave.State.GetMotionData()
				fmt.Println(md.Euler)
			}
		}

		time.Sleep(500 * time.Microsecond)
	}

	<-wave.BLE.Client.Disconnected()
	log.Printf("Wave is disconnecting")
	// Disconnect the connection
	wave.Disconnect()
}
