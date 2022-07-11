package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type Network compute.Network

func (n Network) Keys() []string {
	return []string{fmt.Sprint(n.ID), n.Name, n.CIDR}
}

func (n Network) Columns() []string {
	return []string{"id", "name", "location", "cidr", "usage"}
}

func (n Network) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":       n.ID,
		"name":     n.Name,
		"location": common.Location(n.Location),
		"cidr":     n.CIDR,
		"usage":    fmt.Sprintf("%d/%d", n.UsedIPs, n.TotalIPs),
	}
}

func (n Network) String() string {
	return fmt.Sprintf("%s (%s)", n.Name, n.CIDR)
}

type NetworkService struct {
	delegate compute.NetworkService
}

func NewNetworkService(client goclient.Client) NetworkService {
	return NetworkService{
		delegate: compute.NewNetworkService(client),
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

type NetworkCreate = compute.NetworkCreate

func (n NetworkService) Create(ctx context.Context, data NetworkCreate) (Network, error) {
	res, err := n.delegate.Create(ctx, data)
	if err != nil {
		return Network{}, err
	}

	return Network(res), nil
}

type NetworkUpdate = compute.NetworkUpdate

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
