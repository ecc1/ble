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
	conn.Print(os.Stdout)
}
