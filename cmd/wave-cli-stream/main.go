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
			// If the Button was pressed down, display a dot, recenter, and request BatteryStatus
			if btnAction == "Down" {
				frame := [][]byte{
					{000, 000, 000, 000, 000, 000, 000, 000, 000},
					{000, 000, 000, 000, 255, 000, 000, 000, 000},
					{000, 000, 000, 255, 255, 255, 000, 000, 000},
					{000, 000, 000, 000, 255, 000, 000, 000, 000},
					{000, 000, 000, 000, 000, 000, 000, 000, 000},
				}
				err := wave.SetDisplay(frame)
				if err != nil {
					fmt.Println(err)
				}
				err = wave.Recenter()
				if err != nil {
					fmt.Println(err)
				}
				err = wave.SendBatteryStatusRequest()
				if err != nil {
					fmt.Println(err)
				}
				bs := wave.State.GetBatteryStatus()
				fmt.Println(bs)
			}
			lastState.Buttons.Set(middleButton)
		} else {
			// While the button is held down, print the current motion data
			if btnAction == "Long" || btnAction == "ExtraLong" {
				md := wave.State.GetMotionData()
				fmt.Println(md.Euler)
				frame := [][]byte{
					{000, 000, 000, 000, 000, 000, 000, 000, 000},
					{000, 000, 000, 255, 255, 255, 000, 000, 000},
					{000, 000, 255, 255, 255, 255, 255, 000, 000},
					{000, 000, 000, 255, 255, 255, 000, 000, 000},
					{000, 000, 000, 000, 000, 000, 000, 000, 000},
				}
				err := wave.SetDisplay(frame)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		time.Sleep(500 * time.Microsecond)
	}

	<-wave.BLE.Client.Disconnected()
	log.Printf("Wave is disconnecting")
	// Disconnect the connection
	wave.Disconnect()
}
