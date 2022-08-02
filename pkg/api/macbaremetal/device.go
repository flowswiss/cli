package macbaremetal

import (
	"bytes"
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/macbaremetal"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type Device macbaremetal.Device

func (d Device) String() string {
	return d.Name
}

func (d Device) Keys() []string {
	return []string{fmt.Sprint(d.ID), d.Name, d.Hostname}
}

func (d Device) Columns() []string {
	return []string{"id", "name", "location", "product", "operating system", "public ip", "network", "hostname", "status"}
}

func (d Device) Values() map[string]interface{} {
	networkBuffer := &bytes.Buffer{}
	publicIPBuffer := &bytes.Buffer{}

	networkBuffer.WriteString(fmt.Sprintf("%s (", d.Network.Name))
	for i, iface := range d.NetworkInterfaces {
		networkBuffer.WriteString(iface.PrivateIP)

		if iface.PublicIP != "" {
			publicIPBuffer.WriteString(fmt.Sprintf("%s, ", iface.PublicIP))
		}

		if i+1 < len(d.NetworkInterfaces) {
			networkBuffer.WriteString(", ")
		}
	}
	networkBuffer.WriteRune(')')

	publicIP := publicIPBuffer.String()
	if len(publicIP) > 0 {
		publicIP = publicIP[:len(publicIP)-2]
	}

	return map[string]interface{}{
		"id":               d.ID,
		"name":             d.Name,
		"location":         d.Location.Name,
		"product":          common.Product(d.Product),
		"operating system": fmt.Sprintf("%s %s", d.OperatingSystem.Name, d.OperatingSystem.Version),
		"public ip":        publicIP,
		"network":          networkBuffer.String(),
		"hostname":         d.Hostname,
		"status":           d.Status.Name,
	}
}

type DeviceService struct {
	delegate macbaremetal.DeviceService
}

func NewDeviceService(client goclient.Client) DeviceService {
	return DeviceService{
		delegate: macbaremetal.NewDeviceService(client),
	}
}

func (d DeviceService) List(ctx context.Context) ([]Device, error) {
	res, err := d.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Device, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Device(item)
	}

	return items, nil
}

func (d DeviceService) Get(ctx context.Context, id int) (Device, error) {
	device, err := d.delegate.Get(ctx, id)
	return Device(device), err
}

type DeviceVNCConnection = macbaremetal.DeviceVNCConnection

func (d DeviceService) GetVNC(ctx context.Context, id int) (DeviceVNCConnection, error) {
	return d.delegate.GetVNC(ctx, id)
}

type DeviceCreate = macbaremetal.DeviceCreate

func (d DeviceService) Create(ctx context.Context, data DeviceCreate) (common.Ordering, error) {
	res, err := d.delegate.Create(ctx, data)
	if err != nil {
		return common.Ordering{}, err
	}

	return res, nil
}

type DeviceUpdate = macbaremetal.DeviceUpdate

func (d DeviceService) Update(ctx context.Context, id int, data DeviceUpdate) (Device, error) {
	res, err := d.delegate.Update(ctx, id, data)
	if err != nil {
		return Device{}, err
	}

	return Device(res), nil
}

func (d DeviceService) Delete(ctx context.Context, id int) error {
	return d.delegate.Delete(ctx, id)
}
