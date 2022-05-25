package macbaremetal

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/macbaremetal"
)

type Router macbaremetal.Router

func (r Router) Keys() []string {
	return []string{fmt.Sprint(r.ID), r.Name, r.PublicIP}
}

func (r Router) Columns() []string {
	return []string{"id", "name", "location", "public ip"}
}

func (r Router) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":        r.ID,
		"name":      r.Name,
		"location":  r.Location.Name,
		"public ip": r.PublicIP,
	}
}

type RouterService struct {
	delegate macbaremetal.RouterService
}

func NewRouterService(client goclient.Client) RouterService {
	return RouterService{
		delegate: macbaremetal.NewRouterService(client),
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

type RouterUpdate = macbaremetal.RouterUpdate

func (r RouterService) Update(ctx context.Context, id int, data RouterUpdate) (Router, error) {
	res, err := r.delegate.Update(ctx, id, data)
	if err != nil {
		return Router{}, err
	}

	return Router(res), nil
}
