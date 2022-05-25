package macbaremetal

import (
	"context"
	"fmt"
	"net"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/macbaremetal"
)

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

func (s SecurityGroupRule) Values() map[string]interface{} {
	protocolName := "unknown"
	if name, ok := ProtocolNames[s.Protocol]; ok {
		protocolName = name
	}

	return map[string]interface{}{
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

type SecurityGroupRuleService struct {
	delegate macbaremetal.SecurityGroupRuleService
}

func NewSecurityGroupRuleService(client goclient.Client, securityGroupID int) SecurityGroupRuleService {
	return SecurityGroupRuleService{
		delegate: macbaremetal.NewSecurityGroupRuleService(client, securityGroupID),
	}
}

func (s SecurityGroupRuleService) List(ctx context.Context) ([]SecurityGroupRule, error) {
	res, err := s.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]SecurityGroupRule, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = SecurityGroupRule(item)
	}

	return items, nil
}

type SecurityGroupRuleCreate = macbaremetal.SecurityGroupRuleOptions

func (s SecurityGroupRuleService) Create(ctx context.Context, data SecurityGroupRuleCreate) (SecurityGroupRule, error) {
	res, err := s.delegate.Create(ctx, data)
	if err != nil {
		return SecurityGroupRule{}, err
	}

	return SecurityGroupRule(res), nil
}

type SecurityGroupRuleUpdate = macbaremetal.SecurityGroupRuleOptions

func (s SecurityGroupRuleService) Update(ctx context.Context, id int, data SecurityGroupRuleUpdate) (SecurityGroupRule, error) {
	res, err := s.delegate.Update(ctx, id, data)
	if err != nil {
		return SecurityGroupRule{}, err
	}

	return SecurityGroupRule(res), nil
}

func (s SecurityGroupRuleService) Delete(ctx context.Context, id int) error {
	return s.delegate.Delete(ctx, id)
}
