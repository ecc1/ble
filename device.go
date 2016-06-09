package ble

import (
	"log"
)

const (
	deviceInterface = "org.bluez.Device1"
	interfacesAdded = "org.freedesktop.DBus.ObjectManager.InterfacesAdded"
)

// The Device type corresponds to the org.bluez.Device1 interface.
// See bluez/doc/devicet-api.txt
type Device interface {
	BaseObject

	UUIDs() []string
	Connected() bool
	Paired() bool

	Connect() error
	Pair() error
}

// GetDevice finds a Device in the object cache with the given UUIDs.
func GetDevice(uuids ...string) (Device, error) {
	return findObject(deviceInterface, func(device *blob) bool {
		if uuids != nil {
			advertised := device.UUIDs()
			for _, u := range uuids {
				if !validUUID(u) {
					log.Fatalln("invalid UUID", u)
				}
				if !stringArrayContains(advertised, u) {
					return false
				}
			}
		}
		return true
	})
}

func (device *blob) UUIDs() []string {
	return device.properties["UUIDs"].Value().([]string)
}

func (device *blob) Connected() bool {
	return device.properties["Connected"].Value().(bool)
}

func (device *blob) Paired() bool {
	return device.properties["Paired"].Value().(bool)
}

func (device *blob) Connect() error {
	log.Printf("%s: connecting", device.Name())
	return device.call("Connect")
}

func (device *blob) Pair() error {
	log.Printf("%s: pairing", device.Name())
	return device.call("Pair")
}

func stringArrayContains(a []string, str string) bool {
	for _, s := range a {
		if s == str {
			return true
		}
	}
	return false
}
