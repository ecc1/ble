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
	uuids := os.Args[1:]
	device, err := conn.Discover(0, "", uuids...)
	if err != nil {
		log.Fatal(err)
	}
	device.Print(os.Stdout)
}
