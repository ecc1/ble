/*
Package ble provides functions to discover, connect, pair,
and communicate with Bluetooth Low Energy peripheral devices.

This implementation uses the BlueZ D-Bus interface, rather than sockets.
It is similar to github.com/adafruit/Adafruit_Python_BluefruitLE
*/
package ble

import (
	"fmt"

	"github.com/godbus/dbus"
)

const (
	objectManager = "org.freedesktop.DBus.ObjectManager"
)

var (
	bus *dbus.Conn

	// It would be nice to factor out the subtypes here,
	// but then the reflection used by dbus.Store() wouldn't work.
	objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
)

func init() {
	var err error
	bus, err = dbus.SystemBus()
	if err != nil {
		panic(err)
	}
	err = Update()
	if err != nil {
		panic(err)
	}
}

// Get all objects and properties.
// See http://dbus.freedesktop.org/doc/dbus-specification.html#standard-interfaces-objectmanager
func Update() error {
	call := bus.Object("org.bluez", "/").Call(
		dot(objectManager, "GetManagedObjects"),
		0,
	)
	return call.Store(&objects)
}

type dbusInterfaces *map[string]map[string]dbus.Variant

// The iterObjects function applies a function of type objectProc to
// each object in the cache.  It should return true if the iteration
// should stop, false if it should continue.
type objectProc func(dbus.ObjectPath, dbusInterfaces) bool

func iterObjects(proc objectProc) {
	for path, dict := range objects {
		if proc(path, &dict) {
			return
		}
	}
}

// Print prints the objects om the cache.
func Print() {
	iterObjects(printObject)
}

func printObject(path dbus.ObjectPath, dict dbusInterfaces) bool {
	fmt.Println(path)
	for iface, props := range *dict {
		printProperties(iface, props)
	}
	fmt.Println()
	return false
}

// The BaseObject interface wraps basic operations on a D-Bus object.
//
// Path returns the object's path.
//
// Interface returns the name of the D-Bus interface provided by the object.
//
// Name returns the object's name.
//
// Print prints the object.
type BaseObject interface {
	Path() dbus.ObjectPath
	Interface() string
	Name() string
	Print()
}

type properties map[string]dbus.Variant

type blob struct {
	path       dbus.ObjectPath
	iface      string
	properties properties
	object     dbus.BusObject
}

func (obj *blob) Path() dbus.ObjectPath {
	return obj.path
}

func (obj *blob) Interface() string {
	return obj.iface
}

func (obj *blob) Name() string {
	name, ok := obj.properties["Name"].Value().(string)
	if ok {
		return name
	} else {
		return string(obj.path)
	}
}

func (obj *blob) callv(method string, args ...interface{}) *dbus.Call {
	return obj.object.Call(dot(obj.iface, method), 0, args...)
}

func (obj *blob) call(method string, args ...interface{}) error {
	return obj.callv(method, args...).Err
}

func (obj *blob) Print() {
	fmt.Printf("%s [%s]\n", obj.path, obj.iface)
	printProperties("", obj.properties)
}

func printProperties(iface string, props properties) {
	indent := "    "
	if iface != "" {
		fmt.Printf("%s%s\n", indent, iface)
		indent += indent
	}
	for key, val := range props {
		fmt.Printf("%s%s %s\n", indent, key, val.String())
	}
}

// The findObject function tests each object with functions of type predicate.
type predicate func(*blob) bool

// findObject finds an object satisfying the given tests.
// If returns an error if zero or more than one is found.
func findObject(iface string, tests ...predicate) (*blob, error) {
	var found []*blob
	iterObjects(func(path dbus.ObjectPath, dict dbusInterfaces) bool {
		props := (*dict)[iface]
		if props == nil {
			return false
		}
		obj := &blob{
			path:       path,
			iface:      iface,
			properties: props,
			object:     bus.Object("org.bluez", path),
		}
		for _, test := range tests {
			if !test(obj) {
				return false
			}
		}
		found = append(found, obj)
		return false
	})
	switch len(found) {
	case 1:
		return found[0], nil
	default:
		return nil, fmt.Errorf("found %d instances of interface %s", len(found), iface)
	}
}

func dot(a, b string) string {
	return a + "." + b
}
