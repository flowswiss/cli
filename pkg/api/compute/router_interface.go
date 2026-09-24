package compute

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
)

type (
	RouterInterfaceCreate = compute.RouterInterfaceCreateReq
	RouterInterfaceList   = compute.RouterInterfaceListReq
	RouterInterfaceDelete = compute.RouterInterfaceDeleteReq
)

type GenericRouterInterfaceService = generic.CFD[
	compute.RouterInterface,
	RouterInterface,
	RouterInterfaceCreate,
	RouterInterfaceList,
	RouterInterfaceDelete,
]

func RouterInterfaceService() GenericRouterInterfaceService {
	return generic.NewCFD[
		compute.RouterInterface,
		RouterInterface,
		RouterInterfaceCreate,
		RouterInterfaceList,
		RouterInterfaceDelete,
	](commands.Client.Compute.RouterInterface, func(routerInterface compute.RouterInterface) RouterInterface {
		return RouterInterface(routerInterface)
	})
}

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

func (r RouterInterface) Values() map[string]any {
	return map[string]any{
		"id":         r.ID,
		"network":    Network(r.Network),
		"private ip": r.PrivateIP,
	}
}
