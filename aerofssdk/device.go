package aerofssdk

import (
	"encoding/json"
	"errors"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
)

// Device, client wrapper
type DeviceClient struct {
	APIClient *api.Client
	Desc      Device
}

// Device descriptors
type Device api.Device
type DeviceStatus api.DeviceStatus

// Retrieve a list of existing Device descriptors
func ListDevices(c *api.Client, email string) ([]Device, error) {
	body, _, err := c.ListDevices(email)
	if err != nil {
		return nil, err
	}

	devices := []Device{}
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return nil, errors.New("Unable to demarshal list of devices")
	}
	return devices, err
}

// Return an existing device client given a deviceId
func NewDeviceClient(c *api.Client, deviceId string) (*DeviceClient, error) {
	body, _, err := c.GetDeviceMetadata(deviceId)
	if err != nil {
		return nil, err
	}
	device := Device{}
	err = json.Unmarshal(body, &device)
	return &DeviceClient{c, device}, err
}

// Update the name of the device
func (c *DeviceClient) Update(name string) error {
	body, _, err := c.APIClient.UpdateDevice(name)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &c.Desc)
	if err != nil {
		return errors.New("Unable to demarshal updated device metadata")
	}

	return nil
}

// Retrieve the status of the current device
func (c *DeviceClient) Status() (*DeviceStatus, error) {
	body, _, err := c.APIClient.GetDeviceStatus(c.Desc.Id)
	if err != nil {
		return nil, err
	}

	deviceStatus := new(DeviceStatus)
	err = json.Unmarshal(body, deviceStatus)
	if err != nil {
		return nil, errors.New("Unable to demarshal current device status")
	}

	return deviceStatus, err
}
