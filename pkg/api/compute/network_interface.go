package compute

import (
	"fmt"
	"strings"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
)

type (
	NetworkInterfaceCreate              = compute.NetworkInterfaceCreateReq
	NetworkInterfaceList                = compute.NetworkInterfaceListReq
	NetworkInterfaceDelete              = compute.NetworkInterfaceDeleteReq
	NetworkInterfaceSecurityUpdate      = compute.NetworkInterfaceSecurityUpdateReq
	NetworkInterfaceSecurityGroupUpdate = compute.NetworkInterfaceSecurityGroupUpdateReq
)

type GenericNetworkInterfaceService = generic.CFD[
	compute.NetworkInterface,
	NetworkInterface,
	NetworkInterfaceCreate,
	NetworkInterfaceList,
	NetworkInterfaceDelete,
]

func NetworkInterfaceService() GenericNetworkInterfaceService {
	return generic.NewCFD[
		compute.NetworkInterface,
		NetworkInterface,
		NetworkInterfaceCreate,
		NetworkInterfaceList,
		NetworkInterfaceDelete,
	](commands.Client.Compute.NetworkInterface, func(networkInterface compute.NetworkInterface) NetworkInterface {
		return NetworkInterface(networkInterface)
	})
}

type NetworkInterface compute.NetworkInterface

func (n NetworkInterface) String() string {
	return fmt.Sprintf("%s (%s)", n.MacAddress, n.PrivateIP)
}

func (n NetworkInterface) Keys() []string {
	return []string{fmt.Sprint(n.ID), n.MacAddress, n.PrivateIP}
}

func (n NetworkInterface) Columns() []string {
	return []string{"id", "mac address", "private ip", "network", "security groups", "attached elastic ip"}
}

func (n NetworkInterface) Values() map[string]any {
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

	return map[string]any{
		"id":                  n.ID,
		"mac address":         n.MacAddress,
		"private ip":          n.PrivateIP,
		"network":             Network(n.Network),
		"security groups":     securityGroupBuffer.String(),
		"attached elastic ip": n.AttachedElasticIP.PublicIP,
	}
}
