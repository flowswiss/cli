package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	ServerCreate    = compute.ServerCreateReq
	ServerGet       = compute.ServerGetReq
	ServerList      = core.Cursor
	ServerUpdate    = compute.ServerUpdateReq
	ServerUpgrade   = compute.ServerUpgradeReq
	ServerRunAction = compute.ServerPerformReq
	ServerDelete    = compute.ServerDeleteReq
)

type GenericServerService struct {
	generic.OrderedCreateService[ServerCreate]
	generic.Read[compute.Server, Server, ServerGet, ServerList]
	generic.UpdateService[ServerUpdate, compute.Server, Server]
	generic.PerformActionService[ServerRunAction, compute.Server, Server]
	generic.DeleteService[ServerDelete]

	client *compute.ServerService
}

func (s GenericServerService) Upgrade(ctx context.Context, upgrade ServerUpgrade) (common.Ordering, error) {
	order, err := s.client.Upgrade(ctx, upgrade)
	if err != nil {
		return common.Ordering{}, err
	}

	return order, err
}

func ServerService() GenericServerService {
	client := commands.Client.Compute.Server
	cast := func(server compute.Server) Server {
		return Server(server)
	}

	return GenericServerService{
		generic.NewOrderedCreate[ServerCreate](client),
		generic.NewRead[compute.Server, Server, ServerGet, ServerList](client, cast),
		generic.NewUpdate[ServerUpdate, compute.Server, Server](client, cast),
		generic.NewPerformAction[ServerRunAction, compute.Server, Server](client, cast),
		generic.NewDelete[ServerDelete](client),
		client,
	}
}

type Server compute.Server

func (s Server) String() string {
	return s.Name
}

func (s Server) Keys() (identifiers []string) {
	for _, network := range s.Networks {
		identifiers = append(identifiers, network.Name, network.CIDR)

		for _, iface := range network.Interfaces {
			identifiers = append(identifiers, fmt.Sprint(iface.ID), iface.PrivateIP, iface.PublicIP)
		}
	}

	return append(identifiers, fmt.Sprint(s.ID), s.Name)
}

func (s Server) Columns() []string {
	return []string{"id", "name", "status", "product", "operating system", "location", "public ip", "network"}
}

func (s Server) Values() map[string]any {
	networkBuffer := &strings.Builder{}
	publicIPBuffer := &strings.Builder{}

	for i, network := range s.Networks {
		if i != 0 {
			networkBuffer.WriteString(", ")
		}

		networkBuffer.WriteString(fmt.Sprintf("%s (", network.Name))
		for j, iface := range network.Interfaces {
			if j != 0 {
				networkBuffer.WriteString(", ")
			}

			networkBuffer.WriteString(iface.PrivateIP)

			if iface.PublicIP != "" {
				publicIPBuffer.WriteString(fmt.Sprintf("%s, ", iface.PublicIP))
			}
		}
		networkBuffer.WriteRune(')')
	}

	publicIP := publicIPBuffer.String()
	if len(publicIP) > 0 {
		publicIP = publicIP[:len(publicIP)-2]
	}

	return map[string]any{
		"id":               s.ID,
		"name":             s.Name,
		"status":           s.Status.Name,
		"product":          common.Product(s.Product),
		"operating system": Image{Image: s.Image},
		"location":         common.Location(s.Location),
		"public ip":        publicIP,
		"network":          networkBuffer.String(),
	}
}
