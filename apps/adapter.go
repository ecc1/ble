package main

import (
	"log"
	"os"

	"github.com/ecc1/ble"
)

func main() {
	adapter, err := ble.GetAdapter()
	if err != nil {
		log.Fatal(err)
	}
	adapter.Print(os.Stdout)
}
