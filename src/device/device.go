package device

type device struct {
	name string
}

func NewDevice(name string) device {
	return device{name}
}

func (device *device) getName() string {
	return device.name
}
