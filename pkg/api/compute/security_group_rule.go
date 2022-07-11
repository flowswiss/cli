package compute

import (
	"context"
	"fmt"
	"net"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

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

func (s SecurityGroupRule) Values() map[string]interface{} {
	protocolName := fmt.Sprintf("unknown (%d)", s.Protocol)
	if name, ok := ProtocolNames[s.Protocol]; ok {
		protocolName = name
	}

	return map[string]interface{}{
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

type SecurityGroupRuleService struct {
	delegate compute.SecurityGroupRuleService
}

func NewSecurityGroupRuleService(client goclient.Client, securityGroupID int) SecurityGroupRuleService {
	return SecurityGroupRuleService{
		delegate: compute.NewSecurityGroupRuleService(client, securityGroupID),
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

type SecurityGroupRuleCreate = compute.SecurityGroupRuleOptions

func (s SecurityGroupRuleService) Create(ctx context.Context, data SecurityGroupRuleCreate) (SecurityGroupRule, error) {
	res, err := s.delegate.Create(ctx, data)
	if err != nil {
		return SecurityGroupRule{}, err
	}

	return SecurityGroupRule(res), nil
}

type SecurityGroupRuleUpdate = compute.SecurityGroupRuleOptions

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
