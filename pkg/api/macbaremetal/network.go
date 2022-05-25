package macbaremetal

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/macbaremetal"
)

type Network macbaremetal.Network

func (n Network) Keys() []string {
	return []string{fmt.Sprint(n.ID), n.Name, n.Subnet}
}

func (n Network) Columns() []string {
	return []string{"id", "name", "location", "subnet", "usage"}
}

func (n Network) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":       n.ID,
		"name":     n.Name,
		"location": n.Location.Name,
		"subnet":   n.Subnet,
		"usage":    fmt.Sprintf("%d/%d", n.UsedIPs, n.TotalIPs),
	}
}

type NetworkService struct {
	delegate macbaremetal.NetworkService
}

func NewNetworkService(client goclient.Client) NetworkService {
	return NetworkService{
		delegate: macbaremetal.NewNetworkService(client),
	}
}

func (n NetworkService) List(ctx context.Context) ([]Network, error) {
	res, err := n.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Network, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Network(item)
	}

	return items, nil
}

type NetworkCreate = macbaremetal.NetworkCreate
type NetworkUpdate = macbaremetal.NetworkUpdate

func (n NetworkService) Create(ctx context.Context, data NetworkCreate) (Network, error) {
	res, err := n.delegate.Create(ctx, data)
	if err != nil {
		return Network{}, err
	}

	return Network(res), nil
}

func (n NetworkService) Update(ctx context.Context, id int, data NetworkUpdate) (Network, error) {
	res, err := n.delegate.Update(ctx, id, data)
	if err != nil {
		return Network{}, err
	}

	return Network(res), nil
}

func (n NetworkService) Delete(ctx context.Context, id int) error {
	return n.delegate.Delete(ctx, id)
}
