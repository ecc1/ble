// +build !nofilter

package ble

import (
	"log"

	"github.com/godbus/dbus"
)

func (adapter *blob) SetDiscoveryFilter(uuids ...string) error {
	log.Printf("%s: setting discovery filter %v", adapter.Name(), uuids)
	return adapter.call(
		"SetDiscoveryFilter",
		Properties{
			"Transport": dbus.MakeVariant("le"),
			"UUIDs":     dbus.MakeVariant(uuids),
		},
	)
}
