// +build !nofilter

package ble

import (
	"log"

	"github.com/godbus/dbus"
)

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
