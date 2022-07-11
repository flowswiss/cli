package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

type Route compute.Route

func (r Route) String() string {
	return fmt.Sprint(r.Destination, " via ", r.NextHop)
}

func (r Route) Keys() []string {
	return []string{fmt.Sprint(r.ID), r.Destination, r.NextHop}
}

func (r Route) Columns() []string {
	return []string{"id", "destination", "next hop"}
}

func (r Route) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":          r.ID,
		"destination": r.Destination,
		"next hop":    r.NextHop,
	}
}

type RouteService struct {
	delegate compute.RouteService
}

func NewRouteService(client goclient.Client, routerID int) RouteService {
	return RouteService{
		delegate: compute.NewRouteService(client, routerID),
	}
}

func (r RouteService) List(ctx context.Context) ([]Route, error) {
	res, err := r.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Route, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Route(item)
	}

	return items, nil
}

type RouteCreate = compute.RouteCreate

func (r RouteService) Create(ctx context.Context, data RouteCreate) (Route, error) {
	res, err := r.delegate.Create(ctx, data)
	if err != nil {
		return Route{}, err
	}

	return Route(res), nil
}

func (r RouteService) Delete(ctx context.Context, id int) error {
	return r.delegate.Delete(ctx, id)
}
