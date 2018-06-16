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
	"time"

	"github.com/godbus/dbus"
)

const (
	objectManager = "org.freedesktop.DBus.ObjectManager"
)

// Use type aliases for these subtypes so the reflection
// used by dbus.Store() works correctly.

// Object represents a managed D-Bus object.
type Object = map[string]properties

type properties = map[string]dbus.Variant

// Connection represents a D-Bus connection.
type Connection struct {
	bus     *dbus.Conn
	objects map[dbus.ObjectPath]Object
}

// Open opens a connection to the system D-Bus
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

// Close closes the D-Bus connection.
func (conn *Connection) Close() {
	_ = conn.bus.Close()
}

// Update gets all objects and properties.
// See http://dbus.freedesktop.org/doc/dbus-specification.html#standard-interfaces-objectmanager
func (conn *Connection) Update() error {
	obj := conn.bus.Object("org.bluez", "/")
	call := obj.Call(dot(objectManager, "GetManagedObjects"), 0)
	return call.Store(&conn.objects)
}

// The iterObjects function applies a function of type objectProc to
// each object in the cache.  It should return true if the iteration
// should stop, false if it should continue.
type objectProc func(dbus.ObjectPath, Object) bool

func (conn *Connection) iterObjects(proc objectProc) {
	for path, dict := range conn.objects {
		if proc(path, dict) {
			return
		}
	}
}

// Print prints the objects in the cache.
func (conn *Connection) Print(w io.Writer) {
	printer := func(path dbus.ObjectPath, dict Object) bool {
		return printObject(w, path, dict)
	}
	conn.iterObjects(printer)
}

func printObject(w io.Writer, path dbus.ObjectPath, dict Object) bool {
	fmt.Fprintln(w, path)
	for iface, props := range dict {
		printProperties(w, iface, props)
	}
	fmt.Fprintln(w)
	return false
}

// BaseObject is the interface satisfied by bluez D-Bus objects.
type BaseObject interface {
	Conn() *Connection
	Path() dbus.ObjectPath
	Interface() string
	Name() string
	Print(io.Writer)
}

type blob struct {
	conn       *Connection
	path       dbus.ObjectPath
	iface      string
	properties properties
	object     dbus.BusObject
}

// Conn returns the object's D-Bus connection.
func (obj *blob) Conn() *Connection {
	return obj.conn
}

// Path returns the object's D-Bus path.
func (obj *blob) Path() dbus.ObjectPath {
	return obj.path
}

// Interface returns the object's D-Bus interface name.
func (obj *blob) Interface() string {
	return obj.iface
}

// Name returns the object's name.
func (obj *blob) Name() string {
	name, ok := obj.properties["Name"].Value().(string)
	if !ok {
		return string(obj.path)
	}
	return name
}

func (obj *blob) callv(method string, args ...interface{}) *dbus.Call {
	const callTimeout = 5 * time.Second
	c := obj.object.Go(dot(obj.iface, method), 0, nil, args...)
	if c.Err == nil {
		select {
		case <-c.Done:
		case <-time.After(callTimeout):
			c.Err = fmt.Errorf("BLE call timeout")
		}
	}
	return c
}

func (obj *blob) call(method string, args ...interface{}) error {
	return obj.callv(method, args...).Err
}

// Print prints the object.
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

// findObject finds an object satisfying the given predicate.
// If returns an error if zero or more than one is found.
func (conn *Connection) findObject(iface string, matching predicate) (*blob, error) {
	var found []*blob
	conn.iterObjects(func(path dbus.ObjectPath, dict Object) bool {
		props := dict[iface]
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
		if matching(obj) {
			found = append(found, obj)
		}
		return false
	})
	switch len(found) {
	case 1:
		return found[0], nil
	case 0:
		return nil, fmt.Errorf("interface %s not found", iface)
	default:
		return nil, fmt.Errorf("found %d instances of interface %s", len(found), iface)
	}
}

func dot(a, b string) string {
	return a + "." + b
}
