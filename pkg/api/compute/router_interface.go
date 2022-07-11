package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

type RouterInterface compute.RouterInterface

func (r RouterInterface) String() string {
	return fmt.Sprintf("%s (%s)", r.Network.Name, r.PrivateIP)
}

func (r RouterInterface) Keys() []string {
	return []string{fmt.Sprint(r.ID), r.Network.Name, r.PrivateIP}
}

func (r RouterInterface) Columns() []string {
	return []string{"id", "network", "private ip"}
}

func (r RouterInterface) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":         r.ID,
		"network":    Network(r.Network),
		"private ip": r.PrivateIP,
	}
}

type RouterInterfaceService struct {
	delegate compute.RouterInterfaceService
}

func NewRouterInterfaceService(client goclient.Client, routerID int) RouterInterfaceService {
	return RouterInterfaceService{
		delegate: compute.NewRouterInterfaceService(client, routerID),
	}
}

func (r RouterInterfaceService) List(ctx context.Context) ([]RouterInterface, error) {
	res, err := r.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]RouterInterface, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = RouterInterface(item)
	}

	return items, nil
}

type RouterInterfaceCreate = compute.RouterInterfaceCreate

func (r RouterInterfaceService) Create(ctx context.Context, data RouterInterfaceCreate) (RouterInterface, error) {
	res, err := r.delegate.Create(ctx, data)
	if err != nil {
		return RouterInterface{}, err
	}

	return RouterInterface(res), nil
}

func (r RouterInterfaceService) Delete(ctx context.Context, id int) error {
	return r.delegate.Delete(ctx, id)
}
