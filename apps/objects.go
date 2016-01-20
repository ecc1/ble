package main

import (
	"log"

	"github.com/ecc1/ble"
)

func main() {
	objects, err := ble.ManagedObjects()
	if err != nil {
		log.Fatal(err)
	}
	objects.Print()
}
