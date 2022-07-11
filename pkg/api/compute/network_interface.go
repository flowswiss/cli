package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

type NetworkInterface compute.NetworkInterface

func (n NetworkInterface) Keys() []string {
	return []string{fmt.Sprint(n.ID), n.MacAddress, n.PrivateIP}
}

func (n NetworkInterface) Columns() []string {
	return []string{"id", "mac address", "private ip", "network", "security groups", "attached elastic ip"}
}

func (n NetworkInterface) Values() map[string]interface{} {
	securityGroupBuffer := strings.Builder{}
	if !n.Security {
		securityGroupBuffer.WriteString("- (security disabled)")
	}

	for idx, securityGroup := range n.SecurityGroups {
		if idx != 0 {
			securityGroupBuffer.WriteString(", ")
		}

		securityGroupBuffer.WriteString(SecurityGroup(securityGroup).String())
	}

	return map[string]interface{}{
		"id":                  n.ID,
		"mac address":         n.MacAddress,
		"private ip":          n.PrivateIP,
		"network":             Network(n.Network),
		"security groups":     securityGroupBuffer.String(),
		"attached elastic ip": n.AttachedElasticIP.PublicIP,
	}
}

type NetworkInterfaceService struct {
	delegate compute.NetworkInterfaceService
}

func NewNetworkInterfaceService(client goclient.Client, serverID int) NetworkInterfaceService {
	return NetworkInterfaceService{
		delegate: compute.NewNetworkInterfaceService(client, serverID),
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

type NetworkInterfaceCreate = compute.NetworkInterfaceCreate

func (n NetworkInterfaceService) Create(ctx context.Context, data NetworkInterfaceCreate) (NetworkInterface, error) {
	res, err := n.delegate.Create(ctx, data)
	if err != nil {
		return NetworkInterface{}, err
	}

	return NetworkInterface(res), nil
}

type NetworkInterfaceSecurityUpdate = compute.NetworkInterfaceSecurityUpdate

func (n NetworkInterfaceService) UpdateSecurity(ctx context.Context, id int, data NetworkInterfaceSecurityUpdate) (NetworkInterface, error) {
	res, err := n.delegate.UpdateSecurity(ctx, id, data)
	if err != nil {
		return NetworkInterface{}, err
	}

	return NetworkInterface(res), nil
}

type NetworkInterfaceSecurityGroupUpdate = compute.NetworkInterfaceSecurityGroupUpdate

func (n NetworkInterfaceService) UpdateSecurityGroups(ctx context.Context, id int, data NetworkInterfaceSecurityGroupUpdate) (NetworkInterface, error) {
	res, err := n.delegate.UpdateSecurityGroups(ctx, id, data)
	if err != nil {
		return NetworkInterface{}, err
	}

	return NetworkInterface(res), nil
}

func (n NetworkInterfaceService) Delete(ctx context.Context, id int) error {
	return n.delegate.Delete(ctx, id)
}
