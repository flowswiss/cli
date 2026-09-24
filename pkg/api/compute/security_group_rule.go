package compute

import (
	"fmt"
	"net"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
)

type (
	SecurityGroupRuleCreate = compute.SecurityGroupRuleCreateReq
	SecurityGroupRuleList   = compute.SecurityGroupRuleListReq
	SecurityGroupRuleUpdate = compute.SecurityGroupRuleUpdateReq
	SecurityGroupRuleDelete = compute.SecurityGroupRuleDeleteReq
)

type GenericSecurityGroupRuleService struct {
	generic.CreateService[SecurityGroupRuleCreate, compute.SecurityGroupRule, SecurityGroupRule]
	generic.FilterService[SecurityGroupRuleList, compute.SecurityGroupRule, SecurityGroupRule]
	generic.UpdateService[SecurityGroupRuleUpdate, compute.SecurityGroupRule, SecurityGroupRule]
	generic.DeleteService[SecurityGroupRuleDelete]
}

func SecurityGroupRuleService() GenericSecurityGroupRuleService {
	client := commands.Client.Compute.SecurityGroupRule
	cast := func(rule compute.SecurityGroupRule) SecurityGroupRule {
		return SecurityGroupRule(rule)
	}

	return GenericSecurityGroupRuleService{
		generic.NewCreate[SecurityGroupRuleCreate, compute.SecurityGroupRule, SecurityGroupRule](client, cast),
		generic.NewFilter[SecurityGroupRuleList, compute.SecurityGroupRule, SecurityGroupRule](client, cast),
		generic.NewUpdate[SecurityGroupRuleUpdate, compute.SecurityGroupRule, SecurityGroupRule](client, cast),
		generic.NewDelete[SecurityGroupRuleDelete](client),
	}
}

var IPRangeAny = net.IPNet{
	IP:   net.IPv4zero,
	Mask: net.IPv4Mask(0, 0, 0, 0),
}

var ProtocolIDs = map[string]int{
	"any":  compute.ProtocolAny,
	"icmp": compute.ProtocolICMP,
	"tcp":  compute.ProtocolTCP,
	"udp":  compute.ProtocolUDP,
}

var ProtocolNames = map[int]string{
	compute.ProtocolAny:  "any",
	compute.ProtocolICMP: "icmp",
	compute.ProtocolTCP:  "tcp",
	compute.ProtocolUDP:  "udp",
}

type SecurityGroupRule compute.SecurityGroupRule

func (s SecurityGroupRule) String() string {
	remote := s.IPRange
	if len(remote) == 0 {
		remote = SecurityGroup(s.RemoteSecurityGroup).String()
		if len(remote) == 0 {
			remote = "any"
		}
	}

	if s.Protocol == compute.ProtocolAny {
		return fmt.Sprintf("%s any %s", s.Direction, remote)
	}

	if s.Protocol == compute.ProtocolICMP {
		return fmt.Sprintf("%s icmp %d %d %s", s.Direction, s.ICMPType, s.ICMPCode, remote)
	}

	return fmt.Sprintf("%s %s %d %d %s", s.Direction, ProtocolNames[s.Protocol], s.FromPort, s.ToPort, remote)
}

func (s SecurityGroupRule) Keys() []string {
	return []string{fmt.Sprint(s.ID)} // TODO
}

func (s SecurityGroupRule) Columns() []string {
	return []string{"id", "direction", "protocol", "from port", "to port", "icmp type", "icmp code", "ip range", "remote security group"}
}

func (s SecurityGroupRule) Values() map[string]any {
	protocolName := fmt.Sprintf("unknown (%d)", s.Protocol)
	if name, ok := ProtocolNames[s.Protocol]; ok {
		protocolName = name
	}

	return map[string]any{
		"id":                    s.ID,
		"direction":             s.Direction,
		"protocol":              protocolName,
		"from port":             s.FromPort,
		"to port":               s.ToPort,
		"icmp type":             s.ICMPType,
		"icmp code":             s.ICMPCode,
		"ip range":              s.IPRange,
		"remote security group": SecurityGroup(s.RemoteSecurityGroup),
	}
}
