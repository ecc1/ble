package ble

import (
	"log"

	"github.com/godbus/dbus"
)

const (
	adapterInterface = "org.bluez.Adapter1"
)

type Adapter blob

func (adapter *Adapter) blob() *blob {
	return (*blob)(adapter)
}

func (adapter *Adapter) Print() {
	adapter.blob().print()
}

func (cache *ObjectCache) GetAdapter() (*Adapter, error) {
	p, err := cache.find(adapterInterface)
	return (*Adapter)(p), err
}

func (adapter *Adapter) SetDiscoveryFilter(uuids ...string) error {
	log.Println("setting discovery filter", uuids)
	return adapter.blob().call(
		"SetDiscoveryFilter",
		map[string]dbus.Variant{
			"Transport": dbus.MakeVariant("le"),
			"UUIDs":     dbus.MakeVariant(uuids),
		},
	)
}

func (adapter *Adapter) StartDiscovery() error {
	log.Println("starting discovery")
	return adapter.blob().call("StartDiscovery")
}

func (adapter *Adapter) StopDiscovery() error {
	log.Println("stopping discovery")
	return adapter.blob().call("StopDiscovery")
}
