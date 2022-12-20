package device

type Device struct {
	name string
}

func (device *Device) getName() string {
	return device.name
}
