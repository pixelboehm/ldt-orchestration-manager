package database

type Device struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	MacAddress string `json:"macAddress"`
	Twin       string `json:"twin"`
	Version    string `json:"version"`
}

func NewDevice(name string, macAddress string, twin string, version string) *Device {
	return &Device{
		Name:       name,
		MacAddress: macAddress,
		Twin:       twin,
		Version:    version}
}
