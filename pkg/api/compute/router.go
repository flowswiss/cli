package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type Router compute.Router

func (r Router) String() string {
	return r.Name
}

func (r Router) Keys() []string {
	return []string{fmt.Sprint(r.ID), r.Name, r.PublicIP}
}

func (r Router) Columns() []string {
	return []string{"id", "name", "location", "public ip", "snat"}
}

func (r Router) Values() map[string]interface{} {
	snat := "disabled"
	if r.SourceNAT {
		snat = "enabled"
	}

	return map[string]interface{}{
		"id":        r.ID,
		"name":      r.Name,
		"location":  common.Location(r.Location),
		"public ip": r.PublicIP,
		"snat":      snat,
	}
}

type RouterService struct {
	delegate compute.RouterService
}

func NewRouterService(client goclient.Client) RouterService {
	return RouterService{
		delegate: compute.NewRouterService(client),
	}
}

func (r RouterService) List(ctx context.Context) ([]Router, error) {
	res, err := r.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Router, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Router(item)
	}

	return items, nil
}

type RouterCreate = compute.RouterCreate

func (r RouterService) Create(ctx context.Context, data RouterCreate) (Router, error) {
	res, err := r.delegate.Create(ctx, data)
	if err != nil {
		return Router{}, err
	}

	return Router(res), nil
}

type RouterUpdate = compute.RouterUpdate

func (r RouterService) Update(ctx context.Context, id int, data RouterUpdate) (Router, error) {
	res, err := r.delegate.Update(ctx, id, data)
	if err != nil {
		return Router{}, err
	}

	return Router(res), nil
}

func (r RouterService) Delete(ctx context.Context, id int) error {
	return r.delegate.Delete(ctx, id)
}
