package macbaremetal

import (
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/macbaremetal"
)

type (
	NetworkInterfaceList                = macbaremetal.NetworkInterfaceListReq
	NetworkInterfaceSecurityGroupUpdate = macbaremetal.NetworkInterfaceSecurityGroupUpdateReq
)

type GenericNetworkInterfaceService struct {
	generic.ListService[NetworkInterfaceList, macbaremetal.NetworkInterface, NetworkInterface]

	client *macbaremetal.NetworkInterfaceService
}

func (n GenericNetworkInterfaceService) UpdateSecurityGroup(
	ctx context.Context,
	update NetworkInterfaceSecurityGroupUpdate,
) (NetworkInterface, error) {
	item, err := n.client.UpdateSecurityGroup(ctx, update)
	if err != nil {
		return NetworkInterface{}, err
	}

	return NetworkInterface(item), err
}

func NetworkInterfaceService() GenericNetworkInterfaceService {
	client := commands.Client.MacBareMetal.NetworkInterface
	cast := func(networkInterface macbaremetal.NetworkInterface) NetworkInterface {
		return NetworkInterface(networkInterface)
	}

	return GenericNetworkInterfaceService{
		generic.NewList[NetworkInterfaceList, macbaremetal.NetworkInterface, NetworkInterface](client, cast),
		client,
	}
}

type NetworkInterface macbaremetal.NetworkInterface

func (n NetworkInterface) Keys() []string {
	keys := []string{fmt.Sprint(n.ID), n.PrivateIP, n.MacAddress}

	if n.AttachedElasticIP.ID != 0 {
		keys = append(keys, n.AttachedElasticIP.PublicIP)
	}

	return keys
}

func (n NetworkInterface) Columns() []string {
	return []string{"id", "mac address", "private ip", "network", "security group", "attached elastic ip"}
}

func (n NetworkInterface) Values() map[string]any {
	return map[string]any{
		"id":                  n.ID,
		"mac address":         n.MacAddress,
		"private ip":          n.PrivateIP,
		"network":             fmt.Sprint(n.Network.Name, " (", n.Network.Subnet, ")"),
		"security group":      n.SecurityGroup.Name,
		"attached elastic ip": n.AttachedElasticIP.PublicIP,
	}
}
