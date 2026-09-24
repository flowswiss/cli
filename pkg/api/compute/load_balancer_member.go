package compute

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
)

type (
	LoadBalancerMemberCreate = compute.LoadBalancerMemberCreateReq
	LoadBalancerMemberList   = compute.LoadBalancerMemberListReq
	LoadBalancerMemberDelete = compute.LoadBalancerMemberDeleteReq
)

type GenericLoadBalancerMemberService = generic.CFD[
	compute.LoadBalancerMember,
	LoadBalancerMember,
	LoadBalancerMemberCreate,
	LoadBalancerMemberList,
	LoadBalancerMemberDelete,
]

func LoadBalancerMemberService() GenericLoadBalancerMemberService {
	return generic.NewCFD[
		compute.LoadBalancerMember,
		LoadBalancerMember,
		LoadBalancerMemberCreate,
		LoadBalancerMemberList,
		LoadBalancerMemberDelete,
	](commands.Client.Compute.LoadBalancerMember, func(member compute.LoadBalancerMember) LoadBalancerMember {
		return LoadBalancerMember(member)
	})
}

type LoadBalancerMember compute.LoadBalancerMember

func (l LoadBalancerMember) Host() string {
	return fmt.Sprintf("%s:%d", l.Address, l.Port)
}

func (l LoadBalancerMember) String() string {
	return l.Name
}

func (l LoadBalancerMember) Keys() []string {
	return []string{fmt.Sprint(l.ID), l.Name, l.Host(), l.Status.Key, l.Status.Name}
}

func (l LoadBalancerMember) Columns() []string {
	return []string{"id", "name", "address", "status"}
}

func (l LoadBalancerMember) Values() map[string]any {
	return map[string]any{
		"id":      l.ID,
		"name":    l.Name,
		"address": l.Host(),
		"status":  l.Status.Name,
	}
}
