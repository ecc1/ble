package main

import (
	"log"
	"os"

	"github.com/ecc1/ble"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s (name|UUID)", os.Args[0])
	}
	d := os.Args[1]
	var err error
	conn, err := ble.Open()
	if err != nil {
		log.Fatal(err)
	}
	var device ble.Device
	if ble.ValidUUID(d) {
		device, err = conn.GetDevice(d)
	} else {
		device, err = conn.GetDeviceByName(d)
	}
	if err != nil {
		log.Fatal(err)
	}
	if !device.Connected() {
		log.Printf("%s: not connected", device.Name())
		return
	}
	err = device.Disconnect()
	if err != nil {
		log.Fatal(err)
	}
}
