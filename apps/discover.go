package main

import (
	"log"
	"time"

	"github.com/ecc1/ble"
)

const dexcomUUID = "f0aca0b1-ebfa-f96f-28da-076c35a521db"

func main() {
	adapter, err := ble.Adapter()
	if err != nil {
		log.Fatal(err)
	}
	err = adapter.SetDiscoveryFilter(dexcomUUID)
	if err != nil {
		log.Fatal(err)
	}
	handler := func(device *ble.Blob) bool {
		device.Print()
		return true
	}
	err = adapter.Discover(handler, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
