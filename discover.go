package ble

import (
	"fmt"
	"log"
	"time"

	"github.com/godbus/dbus"
)

func (adapter *blob) Discover(timeout time.Duration, uuids ...string) error {
	signals := make(chan *dbus.Signal)
	defer close(signals)
	bus.Signal(signals)
	defer bus.RemoveSignal(signals)
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
				log.Printf("%s: discovery finished\n", adapter.Name())
				return nil
			default:
				log.Printf("%s: unexpected signal %s\n", adapter.Name(), s.Name)
			}
		case <-t:
			return fmt.Errorf("discovery timeout")
		}
	}
	return nil
}

func (cache *ObjectCache) Discover(timeout time.Duration, uuids ...string) (Device, error) {
	device, err := cache.GetDevice(uuids...)
	if err == nil {
		log.Printf("%s already discovered\n", device.Name())
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
	err = cache.Update()
	if err != nil {
		log.Fatal(err)
	}
	return cache.GetDevice(uuids...)
}
