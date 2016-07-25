package main

import (
	"log"
	"os"

	"github.com/ecc1/ble"
)

func main() {
	uuids := os.Args[1:]
	device, err := ble.GetDevice(uuids...)
	if err != nil {
		log.Fatal(err)
	}
	device.Print(os.Stdout)
}
