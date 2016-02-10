package main

import (
	"log"
	"os"

	"github.com/ecc1/ble"
)

func main() {
	uuids := os.Args[1:]
	device, err := ble.Discover(0, uuids...)
	if err != nil {
		log.Fatal(err)
	}
	device.Print()
}
