/*
Package ble provides functions to discover, connect, pair,
and communicate with Bluetooth Low Energy peripheral devices.

This implementation uses the BlueZ D-Bus interface, rather than sockets.
It is similar to github.com/adafruit/Adafruit_Python_BluefruitLE
*/
package ble

import (
	"fmt"
	"io"

	"github.com/godbus/dbus"
)

const (
	objectManager = "org.freedesktop.DBus.ObjectManager"
)

type Connection struct {
	bus *dbus.Conn

	// It would be nice to factor out the subtypes here,
	// but then the reflection used by dbus.Store() wouldn't work.
	objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
}

func Open() (*Connection, error) {
	bus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	conn := Connection{bus: bus}
	err = conn.Update()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &conn, nil
}

func (conn *Connection) Close() {
	conn.bus.Close()
}

// Get all objects and properties.
// See http://dbus.freedesktop.org/doc/dbus-specification.html#standard-interfaces-objectmanager
func (conn *Connection) Update() error {
	call := conn.bus.Object("org.bluez", "/").Call(
		dot(objectManager, "GetManagedObjects"),
		0,
	)
	return call.Store(&conn.objects)
}

type dbusInterfaces *map[string]map[string]dbus.Variant

// The iterObjects function applies a function of type objectProc to
// each object in the cache.  It should return true if the iteration
// should stop, false if it should continue.
type objectProc func(dbus.ObjectPath, dbusInterfaces) bool

func (conn *Connection) iterObjects(proc objectProc) {
	for path, dict := range conn.objects {
		if proc(path, &dict) {
			return
		}
	}
}

// Print prints the objects om the cache.
func (conn *Connection) Print(w io.Writer) {
	printer := func(path dbus.ObjectPath, dict dbusInterfaces) bool {
		printObject(w, path, dict)
		return false
	}
	conn.iterObjects(printer)
}

func printObject(w io.Writer, path dbus.ObjectPath, dict dbusInterfaces) bool {
	fmt.Fprintln(w, path)
	for iface, props := range *dict {
		printProperties(w, iface, props)
	}
	fmt.Fprintln(w)
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
	Conn() *Connection
	Path() dbus.ObjectPath
	Interface() string
	Name() string
	Print(io.Writer)
}

type properties map[string]dbus.Variant

type blob struct {
	conn       *Connection
	path       dbus.ObjectPath
	iface      string
	properties properties
	object     dbus.BusObject
}

func (obj *blob) Conn() *Connection {
	return obj.conn
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

func (obj *blob) Print(w io.Writer) {
	fmt.Fprintf(w, "%s [%s]\n", obj.path, obj.iface)
	printProperties(w, "", obj.properties)
}

func printProperties(w io.Writer, iface string, props properties) {
	indent := "    "
	if iface != "" {
		fmt.Fprintf(w, "%s%s\n", indent, iface)
		indent += indent
	}
	for key, val := range props {
		fmt.Fprintf(w, "%s%s %s\n", indent, key, val.String())
	}
}

// The findObject function tests each object with functions of type predicate.
type predicate func(*blob) bool

// findObject finds an object satisfying the given tests.
// If returns an error if zero or more than one is found.
func (conn *Connection) findObject(iface string, tests ...predicate) (*blob, error) {
	var found []*blob
	conn.iterObjects(func(path dbus.ObjectPath, dict dbusInterfaces) bool {
		props := (*dict)[iface]
		if props == nil {
			return false
		}
		obj := &blob{
			conn:       conn,
			path:       path,
			iface:      iface,
			properties: props,
			object:     conn.bus.Object("org.bluez", path),
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
