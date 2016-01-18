package ble

import (
	"fmt"
	"log"
	"time"

	"github.com/godbus/dbus"
)

const (
	adapterInterface = "org.bluez.Adapter1"
	deviceInterface  = "org.bluez.Device1"

	objectManager   = "org.freedesktop.DBus.ObjectManager"
	interfacesAdded = "org.freedesktop.DBus.ObjectManager.InterfacesAdded"
)

var bus *dbus.Conn

func init() {
	var err error
	bus, err = dbus.SystemBus()
	if err != nil {
		panic(err)
	}
}

func dot(a, b string) string {
	return a + "." + b
}

type objects *map[dbus.ObjectPath]map[string]map[string]dbus.Variant

// Get all objects and properties.
// See http://dbus.freedesktop.org/doc/dbus-specification.html#standard-interfaces-objectmanager
func managedObjects() (objects, error) {
	call := bus.Object("org.bluez", "/").Call(
		dot(objectManager, "GetManagedObjects"),
		0,
	)
	if call.Err != nil {
		return nil, call.Err
	}
	var objs map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	err := call.Store(&objs)
	return &objs, err
}

type Properties map[string]dbus.Variant

type Blob struct {
	Path       dbus.ObjectPath
	Properties Properties
	Object     dbus.BusObject
	Interface  string
}

func GetObject(
	path dbus.ObjectPath,
	info map[string]map[string]dbus.Variant,
	iface string,
) *Blob {
	props := info[iface]
	if props == nil {
		return nil
	}
	return &Blob{
		Path:       path,
		Properties: props,
		Object:     bus.Object("org.bluez", path),
		Interface:  iface,
	}
}

func notFound(iface string) error {
	return fmt.Errorf("interface %s not found", iface)
}

func (obj *Blob) call(method string, args ...interface{}) error {
	return obj.Object.Call(dot(obj.Interface, method), 0, args...).Err
}

func (obj *Blob) Print() {
	fmt.Printf("%s [%s]\n", obj.Path, obj.Interface)
	for key, val := range obj.Properties {
		fmt.Println("   ", key, val.String())
	}
}

func GetAdapter() (*Blob, error) {
	objects, err := managedObjects()
	if err != nil {
		return nil, err
	}
	for path, info := range *objects {
		obj := GetObject(path, info, adapterInterface)
		if obj != nil {
			return obj, nil
		}
	}
	return nil, notFound(adapterInterface)
}

func (obj *Blob) SetDiscoveryFilter(uuids ...string) error {
	log.Println("Setting discovery filter", uuids)
	return obj.call(
		"SetDiscoveryFilter",
		map[string]dbus.Variant{
			"Transport": dbus.MakeVariant("le"),
			"UUIDs":     dbus.MakeVariant(uuids),
		},
	)
}

func (obj *Blob) StartDiscovery() error {
	log.Println("Starting discovery")
	return obj.call("StartDiscovery")
}

func (obj *Blob) StopDiscovery() error {
	log.Println("Stopping discovery")
	return obj.call("StopDiscovery")
}

// A function of type DeviceHandler is called when a device is
// discovered. It should return true if discovery should stop,
// false if it should continue.
type DeviceHandler func(*Blob) bool

func (obj *Blob) Discover(handler DeviceHandler, timeout time.Duration) error {
	signals := make(chan *dbus.Signal)
	bus.Signal(signals)
	err := bus.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		"type='signal',interface='org.freedesktop.DBus.ObjectManager',member='InterfacesAdded'",
	).Err
	if err != nil {
		return err
	}
	err = obj.StartDiscovery()
	if err != nil {
		return err
	}
	defer obj.StopDiscovery()
	var t <-chan time.Time
	if timeout != 0 {
		t = time.After(timeout)
	}
	for {
		select {
		case s := <-signals:
			switch s.Name {
			case interfacesAdded:
				log.Println("Handling", s.Name)
				obj, err := getDevice(s)
				if err != nil {
					return err
				}
				if handler(obj) {
					close(signals)
					log.Println("Discovery finished")
					return nil
				}
			default:
				log.Println("Unexpected signal", s.Name)
				log.Println(s.Body)
			}
		case <-t:
			log.Println("Discovery timeout")
			return fmt.Errorf("discovery timeout")
		}
	}
	return nil
}

// Extract the object from the InterfacesAdded signal.
// See http://dbus.freedesktop.org/doc/dbus-specification.html#standard-interfaces-objectmanager
func getDevice(signal *dbus.Signal) (*Blob, error) {
	var path dbus.ObjectPath
	var info map[string]map[string]dbus.Variant
	err := dbus.Store(signal.Body, &path, &info)
	if err != nil {
		return nil, err
	}
	obj := GetObject(path, info, deviceInterface)
	if obj == nil {
		return nil, notFound(deviceInterface)
	}
	return obj, nil
}
