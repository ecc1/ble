package ble

import (
	"fmt"
	"log"

	"github.com/godbus/dbus"
)

const (
	objectManager = "org.freedesktop.DBus.ObjectManager"
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
	// but then the reflection used by dbus.Store() wouldn't work.
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

func (cache *ObjectCache) update() error {
	updated, err := ManagedObjects()
	if err != nil {
		return err
	}
	cache.objects = updated.objects
	return nil
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
	cache.iter(printObject)
}

func printObject(path dbus.ObjectPath, dict propertiesDict) bool {
	fmt.Println(path)
	for iface, props := range *dict {
		printProperties(iface, props)
	}
	fmt.Println()
	return false
}

type base interface {
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
	v := obj.properties["Name"].Value()
	name, ok := v.(string)
	if ok {
		return name
	} else {
		return string(obj.path)
	}
}

func (obj *blob) call(method string, args ...interface{}) error {
	return obj.object.Call(dot(obj.iface, method), 0, args...).Err
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

type predicate func(dbus.ObjectPath, properties) bool

func (cache *ObjectCache) collect(iface string, tests ...predicate) []*blob {
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
			path:       path,
			iface:      iface,
			properties: props,
			object:     bus.Object("org.bluez", path),
		}
		objects = append(objects, obj)
		return false
	})
	return objects
}

func (cache *ObjectCache) find(iface string, tests ...predicate) (*blob, error) {
	objects := cache.collect(iface, tests...)
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

func dot(a, b string) string {
	return a + "." + b
}
