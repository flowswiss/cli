package macbaremetal

import (
	"bytes"
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/flowswiss/goclient/v2/macbaremetal"

	commonclt "github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	DeviceCreate    = macbaremetal.DeviceCreateReq
	DeviceGet       = macbaremetal.DeviceGetReq
	DeviceList      = core.Cursor
	DeviceUpdate    = macbaremetal.DeviceUpdateReq
	DeviceRunAction = macbaremetal.DevicePerformReq
	DeviceDelete    = macbaremetal.DeviceDeleteReq

	DeviceVNC     = macbaremetal.DeviceGetReq
	VNCConnection = macbaremetal.DeviceVNCConnection
	WorkflowList  = macbaremetal.DeviceWorkflowListReq
	WorkflowRun   = macbaremetal.DeviceWorkflowRunReq
)

type GenericDeviceService struct {
	generic.OrderedCreateService[DeviceCreate]
	generic.Read[macbaremetal.Device, Device, DeviceGet, DeviceList]
	generic.UpdateService[DeviceUpdate, macbaremetal.Device, Device]
	generic.PerformActionService[DeviceRunAction, macbaremetal.Device, Device]
	generic.DeleteService[DeviceDelete]

	client *macbaremetal.DeviceService
}

func (d GenericDeviceService) GetVNC(ctx context.Context, vnc DeviceVNC) (VNCConnection, error) {
	return d.client.GetVNC(ctx, vnc)
}

func (d GenericDeviceService) WorkflowList(ctx context.Context, req WorkflowList) ([]DeviceWorkflow, error) {
	listOutput, err := d.client.WorkflowList(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]DeviceWorkflow, len(listOutput.Items))
	for idx, item := range listOutput.Items {
		items[idx] = DeviceWorkflow(item)
	}

	return items, nil
}

func (d GenericDeviceService) WorkflowRun(ctx context.Context, req WorkflowRun) (Device, error) {
	device, err := d.client.WorkflowRun(ctx, req)
	if err != nil {
		return Device{}, err
	}

	return Device(device), err
}

func DeviceService() GenericDeviceService {
	client := commands.Client.MacBareMetal.Device
	cast := func(device macbaremetal.Device) Device {
		return Device(device)
	}

	return GenericDeviceService{
		generic.NewOrderedCreate[DeviceCreate](client),
		generic.NewRead[macbaremetal.Device, Device, DeviceGet, DeviceList](client, cast),
		generic.NewUpdate[DeviceUpdate, macbaremetal.Device, Device](client, cast),
		generic.NewPerformAction[DeviceRunAction, macbaremetal.Device, Device](client, cast),
		generic.NewDelete[DeviceDelete](client),
		client,
	}
}

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

func (d Device) Values() map[string]any {
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

	return map[string]any{
		"id":               d.ID,
		"name":             d.Name,
		"location":         d.Location.Name,
		"product":          commonclt.Product(d.Product),
		"operating system": fmt.Sprintf("%s %s", d.OperatingSystem.Name, d.OperatingSystem.Version),
		"public ip":        publicIP,
		"network":          networkBuffer.String(),
		"hostname":         d.Hostname,
		"status":           d.Status.Name,
	}
}

type DeviceAction macbaremetal.DeviceAction

func (d DeviceAction) Keys() []string {
	return []string{fmt.Sprint(d.ID), d.Name, d.Command}
}

func (d DeviceAction) Columns() []string {
	return []string{"id", "name", "command"}
}

func (d DeviceAction) Values() map[string]any {
	return map[string]any{
		"id":      d.ID,
		"name":    d.Name,
		"command": d.Command,
	}
}

type DeviceWorkflow macbaremetal.DeviceWorkflow

func (d DeviceWorkflow) Keys() []string {
	return []string{fmt.Sprint(d.ID), d.Name, d.Command}
}

func (d DeviceWorkflow) Columns() []string {
	return []string{"id", "name", "command"}
}

func (d DeviceWorkflow) Values() map[string]any {
	return map[string]any{
		"id":      d.ID,
		"name":    d.Name,
		"command": d.Command,
	}
}
