package macbaremetal

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/macbaremetal"
)

type DeviceAction macbaremetal.DeviceAction

func (d DeviceAction) Keys() []string {
	return []string{fmt.Sprint(d.ID), d.Name, d.Command}
}

func (d DeviceAction) Columns() []string {
	return []string{"id", "name", "command"}
}

func (d DeviceAction) Values() map[string]interface{} {
	return map[string]interface{}{
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

func (d DeviceWorkflow) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":      d.ID,
		"name":    d.Name,
		"command": d.Command,
	}
}

type DeviceActionService struct {
	delegate macbaremetal.DeviceActionService
}

func NewDeviceActionService(client goclient.Client, deviceID int) DeviceActionService {
	return DeviceActionService{
		delegate: macbaremetal.NewDeviceActionService(client, deviceID),
	}
}

type DeviceRunAction = macbaremetal.DeviceRunAction

func (d DeviceActionService) Run(ctx context.Context, data DeviceRunAction) (Device, error) {
	res, err := d.delegate.Run(ctx, data)
	if err != nil {
		return Device{}, err
	}

	return Device(res), nil
}

type DeviceWorkflowService struct {
	delegate macbaremetal.DeviceWorkflowService
}

func NewDeviceWorkflowService(client goclient.Client, deviceID int) DeviceWorkflowService {
	return DeviceWorkflowService{
		delegate: macbaremetal.NewDeviceWorkflowService(client, deviceID),
	}
}

func (d DeviceWorkflowService) List(ctx context.Context) ([]DeviceWorkflow, error) {
	res, err := d.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]DeviceWorkflow, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = DeviceWorkflow(item)
	}

	return items, nil
}

type DeviceRunWorkflow = macbaremetal.DeviceRunWorkflow

func (d DeviceWorkflowService) Run(ctx context.Context, data DeviceRunWorkflow) (Device, error) {
	res, err := d.delegate.Run(ctx, data)
	if err != nil {
		return Device{}, err
	}

	return Device(res), nil
}
