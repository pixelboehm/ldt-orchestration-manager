package device

type device struct {
	name       string
	macAddress string
	twin       string
	version    string
}

func NewDevice(name string, macAddress string, twin string, version string) device {
	return device{name, macAddress, twin, version}
}

func (device *device) getName() string {
	return device.name
}

func (device *device) getMacAddress() string {
	return device.macAddress
}

func (device *device) getTwin() string {
	return device.twin
}

func (device *device) getVersion() string {
	return device.version
}
