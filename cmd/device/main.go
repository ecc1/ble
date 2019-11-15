package main

import (
	"log"
	"os"

	"github.com/ecc1/ble"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s (address|name|UUID)", os.Args[0])
	}
	d := os.Args[1]
	var err error
	conn, err := ble.Open()
	if err != nil {
		log.Fatal(err)
	}
	var device ble.Device
	if ble.ValidAddress(d) {
		device, err = conn.GetDeviceByAddress(ble.Address(d))
	} else if ble.ValidUUID(d) {
		device, err = conn.GetDeviceByUUID(d)
	} else {
		device, err = conn.GetDeviceByName(d)
	}
	if err != nil {
		log.Fatal(err)
	}
	device.Print(os.Stdout)
}
