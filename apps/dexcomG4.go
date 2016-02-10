package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ecc1/ble"
)

var (
	// service
	receiverService = dexcomUUID(0xa0b1)

	// characteristics
	heartbeat      = dexcomUUID(0x2b18)
	authentication = dexcomUUID(0xacac)
	sendData       = dexcomUUID(0xb20a)
	receiveData    = dexcomUUID(0xb20b)
)

func dexcomUUID(id uint16) string {
	return "f0ac" + fmt.Sprintf("%04x", id) + "-ebfa-f96f-28da-076c35a521db"
}

var done chan interface{}

var readID = []byte{0x01, 0x01, 0x01, 0x06, 0x00, 0x19, 0x0C, 0x47}

func main() {
	objects, err := ble.ManagedObjects()
	if err != nil {
		log.Fatal(err)
	}

	device, err := objects.Discover(0, receiverService)
	if err != nil {
		log.Fatal(err)
	}

	if !device.Connected() {
		err = device.Connect()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("%s: already connected\n", device.Name())
	}

	if !device.Paired() {
		err = device.Pair()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("%s: already paired\n", device.Name())
	}

	err = objects.Update()
	if err != nil {
		log.Fatal(err)
	}

	err = authenticate(objects)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s: authenticated\n", device.Name())

	// We need to enable heartbeat notifications
	// or else we won't get any receiveData responses.
	err = objects.HandleNotify(heartbeat, handleHeartbeat)
	if err != nil {
		log.Fatal(err)
	}

	err = objects.HandleNotify(receiveData, incomingData)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := objects.GetCharacteristic(sendData)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s: sending readID command\n", device.Name())
	err = tx.WriteValue(readID)
	if err != nil {
		log.Fatal(err)
	}

	// Wait indefinitely to receive heartbeats.
	<-done
}

func incomingData(data []byte) {
	log.Printf("incoming data %q\n", data)
}

func handleHeartbeat(data []byte) {
	log.Printf("heartbeat %v\n", data)
}

const (
	serialNumber = "SMxxxxxxxx"
)

var (
	authCode = []byte(serialNumber + "000000")
)

func authenticate(objects *ble.ObjectCache) error {
	auth, err := objects.GetCharacteristic(authentication)
	if err != nil {
		return err
	}
	data, err := auth.ReadValue()
	if err != nil {
		return err
	}
	if bytes.Equal(data, authCode) {
		log.Println("already authenticated")
		return nil
	}
	return auth.WriteValue(authCode)
}
