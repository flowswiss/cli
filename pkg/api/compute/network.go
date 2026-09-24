package compute

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	NetworkCreate = compute.NetworkCreateReq
	NetworkGet    = compute.NetworkGetReq
	NetworkList   = core.Cursor
	NetworkUpdate = compute.NetworkUpdateReq
	NetworkDelete = compute.NetworkDeleteReq
)

type GenericNetworkService = generic.CRUD[
	compute.Network,
	Network,
	NetworkCreate,
	NetworkGet,
	NetworkList,
	NetworkUpdate,
	NetworkDelete,
]

func NetworkService() GenericNetworkService {
	return generic.NewCRUD[
		compute.Network,
		Network,
		NetworkCreate,
		NetworkGet,
		NetworkList,
		NetworkUpdate,
		NetworkDelete,
	](commands.Client.Compute.Network, func(network compute.Network) Network {
		return Network(network)
	})
}

type Network compute.Network

func (n Network) Keys() []string {
	return []string{fmt.Sprint(n.ID), n.Name, n.CIDR}
}

func (n Network) Columns() []string {
	return []string{"id", "name", "location", "cidr", "usage"}
}

func (n Network) Values() map[string]any {
	return map[string]any{
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
