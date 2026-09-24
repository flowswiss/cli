package macbaremetal

import (
	"fmt"
	"net"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/macbaremetal"
)

type (
	SecurityGroupRuleCreate = macbaremetal.SecurityGroupRuleCreateReq
	SecurityGroupRuleList   = macbaremetal.SecurityGroupRuleListReq
	SecurityGroupRuleUpdate = macbaremetal.SecurityGroupRuleUpdateReq
	SecurityGroupRuleDelete = macbaremetal.SecurityGroupRuleDeleteReq
)

type GenericSecurityGroupRuleService struct {
	generic.CreateService[SecurityGroupRuleCreate, macbaremetal.SecurityGroupRule, SecurityGroupRule]
	generic.ListService[SecurityGroupRuleList, macbaremetal.SecurityGroupRule, SecurityGroupRule]
	generic.UpdateService[SecurityGroupRuleUpdate, macbaremetal.SecurityGroupRule, SecurityGroupRule]
	generic.DeleteService[SecurityGroupRuleDelete]
}

func SecurityGroupRuleService() GenericSecurityGroupRuleService {
	client := commands.Client.MacBareMetal.SecurityGroupRule
	cast := func(rule macbaremetal.SecurityGroupRule) SecurityGroupRule {
		return SecurityGroupRule(rule)
	}

	return GenericSecurityGroupRuleService{
		generic.NewCreate[SecurityGroupRuleCreate, macbaremetal.SecurityGroupRule, SecurityGroupRule](client, cast),
		generic.NewList[SecurityGroupRuleList, macbaremetal.SecurityGroupRule, SecurityGroupRule](client, cast),
		generic.NewUpdate[SecurityGroupRuleUpdate, macbaremetal.SecurityGroupRule, SecurityGroupRule](client, cast),
		generic.NewDelete[SecurityGroupRuleDelete](client),
	}
}

var IPRangeAny = net.IPNet{
	IP:   net.IPv4(0, 0, 0, 0),
	Mask: net.IPv4Mask(0, 0, 0, 0),
}

var ProtocolIDs = map[string]int{
	"icmp": macbaremetal.ProtocolICMP,
	"tcp":  macbaremetal.ProtocolTCP,
	"udp":  macbaremetal.ProtocolUDP,
}

var ProtocolNames = map[int]string{
	macbaremetal.ProtocolICMP: "icmp",
	macbaremetal.ProtocolTCP:  "tcp",
	macbaremetal.ProtocolUDP:  "udp",
}

type SecurityGroupRule macbaremetal.SecurityGroupRule

func (s SecurityGroupRule) Keys() []string {
	return []string{fmt.Sprint(s.ID)} // TODO
}

func (s SecurityGroupRule) Columns() []string {
	return []string{"id", "direction", "protocol", "from port", "to port", "icmp type", "icmp code", "ip range"}
}

func (s SecurityGroupRule) Values() map[string]any {
	protocolName := "unknown"
	if name, ok := ProtocolNames[s.Protocol]; ok {
		protocolName = name
	}

	return map[string]any{
		"id":        s.ID,
		"direction": s.Direction,
		"protocol":  fmt.Sprintf("%d (%s)", s.Protocol, protocolName),
		"from port": s.FromPort,
		"to port":   s.ToPort,
		"icmp type": s.ICMPType,
		"icmp code": s.ICMPCode,
		"ip range":  s.IPRange,
	}
}
