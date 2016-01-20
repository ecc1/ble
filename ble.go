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

type ObjectCache struct {
	// It would be nice to factor out the subtypes here,
	// but then the reflection used by Store() wouldn't work.
	objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
}

// Get all objects and properties.
// See http://dbus.freedesktop.org/doc/dbus-specification.html#standard-interfaces-objectmanager
func ManagedObjects() (*ObjectCache, error) {
	call := bus.Object("org.bluez", "/").Call(
		dot(objectManager, "GetManagedObjects"),
		0,
	)
	if call.Err != nil {
		return nil, call.Err
	}
	var objs ObjectCache
	err := call.Store(&objs.objects)
	return &objs, err
}

type propertiesDict *map[string]map[string]dbus.Variant

// A function of type objectProc is applied to each managed object.
// It should return true if the iteration should stop,
// false if it should continue.
type objectProc func(dbus.ObjectPath, propertiesDict) bool

func (cache *ObjectCache) iter(proc objectProc) {
	for path, dict := range cache.objects {
		if proc(path, &dict) {
			return
		}
	}
}

func (cache *ObjectCache) Print() {
	cache.iter(func(path dbus.ObjectPath, dict propertiesDict) bool {
		fmt.Println(path)
		for iface, props := range *dict {
			fmt.Println("   ", iface)
			for p, v := range props {
				fmt.Println("       ", p, v.String())
			}
		}
		return false
	})
}

type properties map[string]dbus.Variant

type blob struct {
	Path       dbus.ObjectPath
	Interface  string
	properties properties
	object     dbus.BusObject
}

func (obj *blob) call(method string, args ...interface{}) error {
	return obj.object.Call(dot(obj.Interface, method), 0, args...).Err
}

func (obj *blob) print() {
	fmt.Printf("%s [%s]\n", obj.Path, obj.Interface)
	for key, val := range obj.properties {
		fmt.Println("   ", key, val.String())
	}
}

type blobPredicate func(dbus.ObjectPath, properties) bool

func (cache *ObjectCache) find(iface string, tests ...blobPredicate) (*blob, error) {
	var objects []*blob
	cache.iter(func(path dbus.ObjectPath, dict propertiesDict) bool {
		props := (*dict)[iface]
		if props == nil {
			return false
		}
		for _, test := range tests {
			if !test(path, props) {
				return false
			}
		}
		obj := &blob{
			Path:       path,
			Interface:  iface,
			properties: props,
			object:     bus.Object("org.bluez", path),
		}
		objects = append(objects, obj)
		return false
	})
	switch len(objects) {
	case 1:
		return objects[0], nil
	case 0:
		return nil, fmt.Errorf("interface %s not found", iface)
	default:
		log.Printf("WARNING: found %d instances of interface %s\n", len(objects), iface)
		return objects[0], nil
	}
}

type Adapter blob

func (adapter *Adapter) Print() {
	(*blob)(adapter).print()
}

func (cache *ObjectCache) GetAdapter() (*Adapter, error) {
	p, err := cache.find(adapterInterface)
	return (*Adapter)(p), err
}

func (adapter *Adapter) SetDiscoveryFilter(uuids ...string) error {
	log.Println("setting discovery filter", uuids)
	return (*blob)(adapter).call(
		"SetDiscoveryFilter",
		map[string]dbus.Variant{
			"Transport": dbus.MakeVariant("le"),
			"UUIDs":     dbus.MakeVariant(uuids),
		},
	)
}

func (adapter *Adapter) StartDiscovery() error {
	log.Println("starting discovery")
	return (*blob)(adapter).call("StartDiscovery")
}

func (adapter *Adapter) StopDiscovery() error {
	log.Println("stopping discovery")
	return (*blob)(adapter).call("StopDiscovery")
}

type Device blob

func (device *Device) Print() {
	(*blob)(device).print()
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

func (adapter *Adapter) Discover(timeout time.Duration, uuids ...string) error {
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
	err = adapter.SetDiscoveryFilter(uuids...)
	if err != nil {
		return err
	}
	err = adapter.StartDiscovery()
	if err != nil {
		return err
	}
	defer adapter.StopDiscovery()
	var t <-chan time.Time
	if timeout != 0 {
		t = time.After(timeout)
	}
	for {
		select {
		case s := <-signals:
			switch s.Name {
			case interfacesAdded:
				log.Println("received signal", s.Name)
				bus.RemoveSignal(signals)
				close(signals)
				log.Println("discovery finished")
				return nil
			default:
				log.Println("unexpected signal", s.Name)
				log.Println(s.Body)
			}
		case <-t:
			return fmt.Errorf("discovery timeout")
		}
	}
	return nil
}

func (cache *ObjectCache) Discover(timeout time.Duration, uuids ...string) (*Device, error) {
	device, err := cache.GetDevice(uuids...)
	if err == nil {
		log.Printf("device %s already discovered\n", device.Path)
		return device, nil
	}
	adapter, err := cache.GetAdapter()
	if err != nil {
		return nil, err
	}
	err = adapter.Discover(timeout, uuids...)
	if err != nil {
		return nil, err
	}
	// update object cache
	updated, err := ManagedObjects()
	if err != nil {
		log.Fatal(err)
	}
	cache.objects = updated.objects
	return cache.GetDevice(uuids...)
}

func (device *Device) Connect() error {
	log.Println("connect")
	return (*blob)(device).call("Connect")
}

func (device *Device) Pair() error {
	log.Println("pair")
	return (*blob)(device).call("Pair")
}

func dot(a, b string) string {
	return a + "." + b
}

func stringArrayContains(a []string, str string) bool {
	for _, s := range a {
		if s == str {
			return true
		}
	}
	return false
}
