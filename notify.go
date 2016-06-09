package ble

import (
	"fmt"
	"log"

	"github.com/godbus/dbus"
)

var (
	notifySignals = make(chan *dbus.Signal, 100)
	notifyHandler = make(map[dbus.ObjectPath]NotifyHandler)
)

func (char *blob) HandleNotify(handler NotifyHandler) error {
	if len(notifyHandler) == 0 {
		go notifyLoop()
		bus.Signal(notifySignals)
	}
	path := char.Path()
	notifyHandler[path] = handler
	rule := fmt.Sprintf(
		"type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged',path='%s'",
		path,
	)
	err := addMatch(rule)
	if err != nil {
		return err
	}
	return char.StartNotify()
}

func applyHandler(s *dbus.Signal) {
	handler := notifyHandler[s.Path]
	if handler == nil {
		log.Printf("%s: no notify handler", s.Path)
		return
	}
	// Reflection used by dbus.Store() requires explicit type here.
	var changed map[string]dbus.Variant
	dbus.Store(s.Body[1:2], &changed)
	keys := []string{}
	for k, _ := range changed {
		keys = append(keys, k)
	}
	data, ok := changed["Value"].Value().([]byte)
	if ok {
		go handler(data)
	}
}

func notifyLoop() {
	for s := range notifySignals {
		applyHandler(s)
	}
}

// HandleNotify enables notifications from the GATT characterisitc with
// the specified UUID and applies the given handler to them when they arrive.
func HandleNotify(uuid string, handler NotifyHandler) error {
	char, err := GetCharacteristic(uuid)
	if err != nil {
		return err
	}
	return char.HandleNotify(handler)
}
