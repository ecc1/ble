package ble

import (
	"log"

	"github.com/godbus/dbus"
)

const (
	deviceInterface = "org.bluez.Device1"
	interfacesAdded = "org.freedesktop.DBus.ObjectManager.InterfacesAdded"
)

type Device interface {
	base

	Connect() error
	Pair() error
}

func (cache *ObjectCache) GetDevice(uuids ...string) (Device, error) {
	return cache.find(deviceInterface, func(path dbus.ObjectPath, props properties) bool {
		if uuids != nil {
			v := props["UUIDs"].Value()
			advertised, ok := v.([]string)
			if !ok {
				log.Fatalln("unexpected UUIDs property:", v)
			}
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

func (device *blob) Connect() error {
	log.Println("connect")
	return device.call("Connect")
}

func (device *blob) Pair() error {
	log.Println("pair")
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
