package ble

import (
	"log"

	"github.com/godbus/dbus"
)

const (
	deviceInterface = "org.bluez.Device1"
	interfacesAdded = "org.freedesktop.DBus.ObjectManager.InterfacesAdded"
)

type Device blob

func (device *Device) blob() *blob {
	return (*blob)(device)
}

func (device *Device) Print() {
	device.blob().print()
}

func (cache *ObjectCache) GetDevice(uuids ...string) (*Device, error) {
	p, err := cache.find(deviceInterface, func(path dbus.ObjectPath, props properties) bool {
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
	return (*Device)(p), err
}

func (device *Device) Connect() error {
	log.Println("connect")
	return device.blob().call("Connect")
}

func (device *Device) Pair() error {
	log.Println("pair")
	return device.blob().call("Pair")
}

func stringArrayContains(a []string, str string) bool {
	for _, s := range a {
		if s == str {
			return true
		}
	}
	return false
}
