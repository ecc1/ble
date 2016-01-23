package main

import (
	"log"

	"github.com/ecc1/ble"
)

const dexcomUUID = "f0aca0b1-ebfa-f96f-28da-076c35a521db"

func main() {
	objects, err := ble.ManagedObjects()
	if err != nil {
		log.Fatal(err)
	}

	device, err := objects.Discover(0, dexcomUUID)
	if err != nil {
		log.Fatal(err)
	}

	if !device.Connected() {
		err = device.Connect()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("%s: already connected\n", device.Name())
	}

	if !device.Paired() {
		err = device.Pair()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("%s: already paired\n", device.Name())
	}
}
