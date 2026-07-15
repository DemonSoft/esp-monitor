package database

import "remoteesp/internal/model"

type Database interface {
	Close()
	Driver() string
	DevicesUploaded(devices []model.Device) error
	Devices(filter string) ([]model.Device, error)
	UpdateDevice(device model.Device) error
	UpdateDeviceAction(ssdp string, action string) error
	DeleteDevice(ssdp string) error
}
