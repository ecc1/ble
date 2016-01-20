package main

import (
	"log"
	"os"

	"github.com/ecc1/ble"
)

func main() {
	objects, err := ble.ManagedObjects()
	if err != nil {
		log.Fatal(err)
	}
	uuids := os.Args[1:]
	device, err := objects.GetDevice(uuids...)
	if err != nil {
		log.Fatal(err)
	}
	device.Print()
}
