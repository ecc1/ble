package ble

import (
	"log"
	"time"

	"github.com/godbus/dbus"
)

const (
	adapterInterface = "org.bluez.Adapter1"
)

type Adapter interface {
	base

	StartDiscovery() error
	StopDiscovery() error
	RemoveDevice(Device) error
	SetDiscoveryFilter(uuids ...string) error

	Discover(timeout time.Duration, uuids ...string) error
}

func (cache *ObjectCache) GetAdapter() (Adapter, error) {
	return cache.find(adapterInterface)
}

func (adapter *blob) StartDiscovery() error {
	log.Printf("%s: starting discovery\n", adapter.Name())
	return adapter.call("StartDiscovery")
}

func (adapter *blob) StopDiscovery() error {
	log.Printf("%s: stopping discovery\n", adapter.Name())
	return adapter.call("StopDiscovery")
}

func (adapter *blob) RemoveDevice(device Device) error {
	log.Printf("%s: removing device", adapter.Name(), device.Name())
	return adapter.call("RemoveDevice", device.Path())
}

func (adapter *blob) SetDiscoveryFilter(uuids ...string) error {
	log.Printf("%s: setting discovery filter %v\n", adapter.Name(), uuids)
	return adapter.call(
		"SetDiscoveryFilter",
		map[string]dbus.Variant{
			"Transport": dbus.MakeVariant("le"),
			"UUIDs":     dbus.MakeVariant(uuids),
		},
	)
}
