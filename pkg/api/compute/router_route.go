package compute

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
)

type (
	RouteCreate = compute.RouteCreateReq
	RouteList   = compute.RouteListReq
	RouteDelete = compute.RouteDeleteReq
)

type GenericRouteService = generic.CFD[
	compute.Route,
	Route,
	RouteCreate,
	RouteList,
	RouteDelete,
]

func RouteService() GenericRouteService {
	return generic.NewCFD[
		compute.Route,
		Route,
		RouteCreate,
		RouteList,
		RouteDelete,
	](commands.Client.Compute.Route, func(route compute.Route) Route {
		return Route(route)
	})
}

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

func (r Route) Values() map[string]any {
	return map[string]any{
		"id":          r.ID,
		"destination": r.Destination,
		"next hop":    r.NextHop,
	}
}
