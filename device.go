package ble

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	deviceInterface = "org.bluez.Device1"
	interfacesAdded = "org.freedesktop.DBus.ObjectManager.InterfacesAdded"
)

// The Device type corresponds to the org.bluez.Device1 interface.
// See bluez/doc/devicet-api.txt
type Device interface {
	BaseObject

	Address() Address
	AddressType() string
	UUIDs() []string
	Connected() bool
	Paired() bool

	Connect() error
	Disconnect() error
	Pair() error
}

func (conn *Connection) matchDevice(matching predicate) (Device, error) {
	return conn.findObject(deviceInterface, matching)
}

// ValidAddress checks whether addr is a valid MAC address.
func ValidAddress(addr string) bool {
	_, err := net.ParseMAC(addr)
	return err == nil
}

// GetDeviceByAddress finds a Device in the object cache with the given address.
func (conn *Connection) GetDeviceByAddress(address Address) (Device, error) {
	addr := Address(strings.ToUpper(string(address)))
	device, err := conn.matchDevice(func(device *blob) bool {
		return device.Address() == addr
	})
	if err != nil {
		err = fmt.Errorf("%w with address %s", err, addr)
	}
	return device, err
}

// GetDeviceByName finds a Device in the object cache with the given name.
func (conn *Connection) GetDeviceByName(name string) (Device, error) {
	device, err := conn.matchDevice(func(device *blob) bool {
		return device.Name() == name
	})
	if err != nil {
		err = fmt.Errorf("%w with name %q", err, name)
	}
	return device, err
}

// GetDeviceByUUID finds a Device in the object cache matching the given UUIDs.
func (conn *Connection) GetDeviceByUUID(uuids ...string) (Device, error) {
	device, err := conn.matchDevice(func(device *blob) bool {
		return UUIDsInclude(device.UUIDs(), uuids)
	})
	if err != nil {
		err = fmt.Errorf("%w with UUIDs %v", err, uuids)
	}
	return device, err
}

// UUIDsInclude tests whether the advertised UUIDs contain all of the ones in uuids.
func UUIDsInclude(advertised []string, uuids []string) bool {
	for _, u := range uuids {
		if !ValidUUID(u) {
			log.Printf("invalid UUID %s", u)
			return false
		}
		if !stringsContain(advertised, LongUUID(u)) {
			return false
		}
	}
	return true
}

func (device *blob) Address() Address {
	return Address(device.properties["Address"].Value().(string))
}

func (device *blob) AddressType() string {
	return device.properties["AddressType"].Value().(string)
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

func (device *blob) Disconnect() error {
	log.Printf("%s: disconnecting", device.Name())
	return device.call("Disconnect")
}

func (device *blob) Pair() error {
	log.Printf("%s: pairing", device.Name())
	return device.call("Pair")
}

func stringsContain(a []string, str string) bool {
	for _, s := range a {
		if s == str {
			return true
		}
	}
	return false
}
