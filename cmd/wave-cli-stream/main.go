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
	"github.com/pkg/errors"
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
	q := decodeByteStream(data)
	// Handle the query
	if q.ID == gowave.ButtonEventID {
		fmt.Println(data)
		buttonEvent, _ := gowave.NewButtonEvent(data)
		handleButtonEvent(buttonEvent)
	}
}

func handleButtonEvent(be gowave.ButtonEvent) {
	fmt.Println(be)
}

func explore(cln ble.Client, p *ble.Profile) error {
	for _, s := range p.Services {
		fmt.Printf("    Service: %s %s, Handle (0x%02X)\n", s.UUID, ble.Name(s.UUID), s.Handle)

		for _, c := range s.Characteristics {
			fmt.Printf("      Characteristic: %s %s, Property: 0x%02X (%s), Handle(0x%02X), VHandle(0x%02X)\n",
				c.UUID, ble.Name(c.UUID), c.Property, propString(c.Property), c.Handle, c.ValueHandle)
			if (c.Property & ble.CharRead) != 0 {
				b, err := cln.ReadCharacteristic(c)
				if err != nil {
					fmt.Printf("Failed to read characteristic: %s\n", err)
					continue
				}
				fmt.Printf("        Value         %x | %q\n", b, b)
			}

			for _, d := range c.Descriptors {
				fmt.Printf("        Descriptor: %s %s, Handle(0x%02x)\n", d.UUID, ble.Name(d.UUID), d.Handle)
				b, err := cln.ReadDescriptor(d)
				if err != nil {
					fmt.Printf("Failed to read descriptor: %s\n", err)
					continue
				}
				fmt.Printf("        Value         %x | %q\n", b, b)
			}

			if *sub != 0 {
				// Don't bother to subscribe the Service Changed characteristics.
				if c.UUID.Equal(ble.ServiceChangedUUID) {
					continue
				}

				// Don't touch the Apple-specific Service/Characteristic.
				// Service: D0611E78BBB44591A5F8487910AE4366
				// Characteristic: 8667556C9A374C9184ED54EE27D90049, Property: 0x18 (WN),
				//   Descriptor: 2902, Client Characteristic Configuration
				//   Value         0000 | "\x00\x00"
				if c.UUID.Equal(ble.MustParse("8667556C9A374C9184ED54EE27D90049")) {
					continue
				}

				if (c.Property & ble.CharNotify) != 0 {
					fmt.Printf("\n-- Subscribe to notification for %s --\n", *sub)
					h := func(req []byte) { fmt.Printf("Notified: %q [ % X ]\n", string(req), req) }
					if err := cln.Subscribe(c, false, h); err != nil {
						log.Fatalf("subscribe failed: %s", err)
					}
					time.Sleep(*sub)
					if err := cln.Unsubscribe(c, false); err != nil {
						log.Fatalf("unsubscribe failed: %s", err)
					}
					fmt.Printf("-- Unsubscribe to notification --\n")
				}
				if (c.Property & ble.CharIndicate) != 0 {
					fmt.Printf("\n-- Subscribe to indication of %s --\n", *sub)
					h := func(req []byte) { fmt.Printf("Indicated: %q [ % X ]\n", string(req), req) }
					if err := cln.Subscribe(c, true, h); err != nil {
						log.Fatalf("subscribe failed: %s", err)
					}
					time.Sleep(*sub)
					if err := cln.Unsubscribe(c, true); err != nil {
						log.Fatalf("unsubscribe failed: %s", err)
					}
					fmt.Printf("-- Unsubscribe to indication --\n")
				}
			}
		}
		fmt.Printf("\n")
	}
	return nil
}

func propString(p ble.Property) string {
	var s string
	for k, v := range map[ble.Property]string{
		ble.CharBroadcast:   "B",
		ble.CharRead:        "R",
		ble.CharWriteNR:     "w",
		ble.CharWrite:       "W",
		ble.CharNotify:      "N",
		ble.CharIndicate:    "I",
		ble.CharSignedWrite: "S",
		ble.CharExtended:    "E",
	} {
		if p&k != 0 {
			s += v
		}
	}
	return s
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
