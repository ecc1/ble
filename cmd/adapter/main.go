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
	adapter, err := conn.GetAdapter()
	if err != nil {
		log.Fatal(err)
	}
	adapter.Print(os.Stdout)
}
