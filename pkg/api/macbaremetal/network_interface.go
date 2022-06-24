package macbaremetal

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/macbaremetal"
)

type NetworkInterface macbaremetal.NetworkInterface

func (n NetworkInterface) Keys() []string {
	keys := []string{fmt.Sprint(n.ID), n.PrivateIP, n.MacAddress}

	if n.AttachedElasticIP.ID != 0 {
		keys = append(keys, n.AttachedElasticIP.PublicIP)
	}

	return keys
}

func (n NetworkInterface) Columns() []string {
	return []string{"id", "private ip", "mac address", "network", "security group", "attached elastic ip"}
}

func (n NetworkInterface) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":                  n.ID,
		"private ip":          n.PrivateIP,
		"mac address":         n.MacAddress,
		"network":             fmt.Sprint(n.Network.Name, " (", n.Network.Subnet, ")"),
		"security group":      n.SecurityGroup.Name,
		"attached elastic ip": n.AttachedElasticIP.PublicIP,
	}
}

type NetworkInterfaceService struct {
	delegate macbaremetal.NetworkInterfaceService
}

func NewNetworkInterfaceService(client goclient.Client, deviceID int) NetworkInterfaceService {
	return NetworkInterfaceService{
		delegate: macbaremetal.NewNetworkInterfaceService(client, deviceID),
	}
}

func (n NetworkInterfaceService) List(ctx context.Context) ([]NetworkInterface, error) {
	res, err := n.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]NetworkInterface, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = NetworkInterface(item)
	}

	return items, nil
}

type NetworkInterfaceSecurityGroupUpdate = macbaremetal.NetworkInterfaceSecurityGroupUpdate

func (n NetworkInterfaceService) UpdateSecurityGroup(ctx context.Context, id int, data NetworkInterfaceSecurityGroupUpdate) (NetworkInterface, error) {
	res, err := n.delegate.UpdateSecurityGroup(ctx, id, data)
	if err != nil {
		return NetworkInterface{}, err
	}

	return NetworkInterface(res), nil
}
