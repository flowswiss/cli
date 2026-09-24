package macbaremetal

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/flowswiss/goclient/v2/macbaremetal"
)

type (
	NetworkCreate = macbaremetal.NetworkCreateReq
	NetworkGet    = macbaremetal.NetworkGetReq
	NetworkList   = core.Cursor
	NetworkUpdate = macbaremetal.NetworkUpdateReq
	NetworkDelete = macbaremetal.NetworkDeleteReq
)

type GenericNetworkService = generic.CRUD[macbaremetal.Network, Network, NetworkCreate, NetworkGet, NetworkList, NetworkUpdate, NetworkDelete]

func NetworkService() GenericNetworkService {
	return generic.NewCRUD[
		macbaremetal.Network,
		Network,
		NetworkCreate,
		NetworkGet,
		NetworkList,
		NetworkUpdate,
		NetworkDelete,
	](commands.Client.MacBareMetal.Network, func(network macbaremetal.Network) Network {
		return Network(network)
	})
}

type Network macbaremetal.Network

func (n Network) String() string {
	return n.Name
}

func (n Network) Keys() []string {
	return []string{fmt.Sprint(n.ID), n.Name, n.Subnet}
}

func (n Network) Columns() []string {
	return []string{"id", "name", "location", "subnet", "usage"}
}

func (n Network) Values() map[string]any {
	return map[string]any{
		"id":       n.ID,
		"name":     n.Name,
		"location": n.Location.Name,
		"subnet":   n.Subnet,
		"usage":    fmt.Sprintf("%d/%d", n.UsedIPs, n.TotalIPs),
	}
}
