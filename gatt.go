package ble

const (
	serviceInterface        = "org.bluez.GattService1"
	characteristicInterface = "org.bluez.GattCharacteristic1"
	descriptorInterface     = "org.bluez.GattDescriptor1"
)

type gattHandle interface {
	base

	UUID() string
}

func (handle *blob) UUID() string {
	return handle.properties["UUID"].Value().(string)
}

type Service interface {
	gattHandle
}

func (cache *ObjectCache) GetService(uuid string) (Service, error) {
	return cache.find(serviceInterface, func(serv *blob) bool {
		return serv.UUID() == uuid
	})
}

type readWriteHandle interface {
	gattHandle

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

type NotifyHandler func([]byte)

type Characteristic interface {
	readWriteHandle

	Notifying() bool

	StartNotify() error
	StopNotify() error

	HandleNotify(NotifyHandler) error
}

func (cache *ObjectCache) GetCharacteristic(uuid string) (Characteristic, error) {
	return cache.find(characteristicInterface, func(char *blob) bool {
		return char.UUID() == uuid
	})
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

type Descriptor interface {
	readWriteHandle
}

func (cache *ObjectCache) GetDescriptor(uuid string) (Descriptor, error) {
	return cache.find(descriptorInterface, func(desc *blob) bool {
		return desc.UUID() == uuid
	})
}
