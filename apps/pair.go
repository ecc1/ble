package main

import (
	"log"
	"os"

	"github.com/ecc1/ble"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s UUID\n", os.Args[0])
	}

	device, err := ble.Discover(0, os.Args[1])
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
