package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

type ServerAction compute.ServerAction

func (d ServerAction) Keys() []string {
	return []string{fmt.Sprint(d.ID), d.Name, d.Command}
}

func (d ServerAction) Columns() []string {
	return []string{"id", "name", "command"}
}

func (d ServerAction) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":      d.ID,
		"name":    d.Name,
		"command": d.Command,
	}
}

type ServerActionService struct {
	delegate compute.ServerService
}

func NewServerActionService(client goclient.Client) ServerActionService {
	return ServerActionService{
		delegate: compute.NewServerService(client),
	}
}

type ServerRunAction = compute.ServerPerform

func (d ServerActionService) Run(ctx context.Context, serverID int, data ServerRunAction) (Server, error) {
	res, err := d.delegate.Perform(ctx, serverID, data)
	if err != nil {
		return Server{}, err
	}

	return Server(res), nil
}
