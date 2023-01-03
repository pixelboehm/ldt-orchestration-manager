package model

type Device struct {
	Name       string
	MacAddress string
	Twin       string `default:"none"`
	Version    string `default:"none"`
}

func NewDevice(name string, macAddress string, twin string, version string) Device {
	return Device{name, macAddress, twin, version}
}

func (device *Device) getName() string {
	return device.Name
}

func (device *Device) getMacAddress() string {
	return device.MacAddress
}

func (device *Device) getTwin() string {
	return device.Twin
}

func (device *Device) getVersion() string {
	return device.Version
}
