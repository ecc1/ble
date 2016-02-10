package ble

const (
	serviceInterface        = "org.bluez.GattService1"
	characteristicInterface = "org.bluez.GattCharacteristic1"
	descriptorInterface     = "org.bluez.GattDescriptor1"
)

func (cache *ObjectCache) getObject(iface string, uuid string) (*blob, error) {
	return cache.find(iface, func(desc *blob) bool {
		return desc.UUID() == uuid
	})
}

// The GattHandle interface wraps common operations on GATT objects.
type GattHandle interface {
	BaseObject

	UUID() string
}

// UUID returns the handle's UUID
func (handle *blob) UUID() string {
	return handle.properties["UUID"].Value().(string)
}

// The Service type corresponds to the org.bluez.GattService1 interface.
// See bluez/doc/gatt-api.txt
type Service interface {
	GattHandle
}

// GetService finds a Service in the object cache with the given UUID.
func (cache *ObjectCache) GetService(uuid string) (Service, error) {
	return cache.getObject(serviceInterface, uuid)
}

// The ReadWriteHandle interface describes GATT objects that provide
// ReadValue and WriteValue operations.
type ReadWriteHandle interface {
	GattHandle

	ReadValue() ([]byte, error)
	WriteValue([]byte) error
}

func (handle *blob) ReadValue() ([]byte, error) {
	var data []byte
	err := handle.callv("ReadValue").Store(&data)
	return data, err
}

func (handle *blob) WriteValue(data []byte) error {
	return handle.call("WriteValue", data)
}

// A function of type NotifyHandler is used to handle notifications.
type NotifyHandler func([]byte)

// The Characteristic type corresponds to the org.bluez.GattCharacteristic1 interface.
// See bluez/doc/gatt-api.txt
type Characteristic interface {
	ReadWriteHandle

	Notifying() bool

	StartNotify() error
	StopNotify() error

	HandleNotify(NotifyHandler) error
}

// GetCharacteristic finds a Characteristic in the object cache with the given UUID.
func (cache *ObjectCache) GetCharacteristic(uuid string) (Characteristic, error) {
	return cache.getObject(characteristicInterface, uuid)
}

func (char *blob) Notifying() bool {
	return char.properties["Notifying"].Value().(bool)
}

func (char *blob) StartNotify() error {
	return char.call("StartNotify")
}

func (char *blob) StopNotify() error {
	return char.call("StopNotify")
}

// The Descriptor type corresponds to the org.bluez.GattDescriptor1 interface.
// See bluez/doc/gatt-api.txt
type Descriptor interface {
	ReadWriteHandle
}

// GetDescriptor finds a Descriptor in the object cache with the given UUID.
func (cache *ObjectCache) GetDescriptor(uuid string) (Descriptor, error) {
	return cache.getObject(descriptorInterface, uuid)
}
