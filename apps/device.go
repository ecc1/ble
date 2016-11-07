package main

import (
	"log"
	"os"

	"github.com/ecc1/ble"
)

func main() {
	conn, err := ble.Open()
	if err != nil {
		log.Fatal(err)
	}
	device := ble.Device(nil)
	if len(os.Args) == 2 && !ble.ValidUUID(os.Args[1]) {
		device, err = conn.GetDeviceByName(os.Args[1])
	} else {
		uuids := os.Args[1:]
		device, err = conn.GetDevice(uuids...)
	}
	if err != nil {
		log.Fatal(err)
	}
	device.Print(os.Stdout)
}
