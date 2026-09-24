package compute

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
)

type (
	LoadBalancerPoolCreate = compute.LoadBalancerPoolCreateReq
	LoadBalancerPoolGet    = compute.LoadBalancerPoolGetReq
	LoadBalancerPoolList   = compute.LoadBalancerPoolListReq
	LoadBalancerPoolUpdate = compute.LoadBalancerPoolUpdateReq
	LoadBalancerPoolDelete = compute.LoadBalancerPoolDeleteReq

	LoadBalancerHealthCheckOptions = compute.LoadBalancerHealthCheckOptions
)

type GenericLoadBalancerPoolService = generic.CRUD[
	compute.LoadBalancerPool,
	LoadBalancerPool,
	LoadBalancerPoolCreate,
	LoadBalancerPoolGet,
	LoadBalancerPoolList,
	LoadBalancerPoolUpdate,
	LoadBalancerPoolDelete,
]

func LoadBalancerPoolService() GenericLoadBalancerPoolService {
	return generic.NewCRUD[
		compute.LoadBalancerPool,
		LoadBalancerPool,
		LoadBalancerPoolCreate,
		LoadBalancerPoolGet,
		LoadBalancerPoolList,
		LoadBalancerPoolUpdate,
		LoadBalancerPoolDelete,
	](commands.Client.Compute.LoadBalancerPool, func(pool compute.LoadBalancerPool) LoadBalancerPool {
		return LoadBalancerPool(pool)
	})
}

type LoadBalancerPool compute.LoadBalancerPool

func (l LoadBalancerPool) NameWithoutSpaces() string {
	return fmt.Sprintf("%s-on-port-%d-to-%s", l.EntryProtocol.Key, l.EntryPort, l.TargetProtocol.Key)
}

func (l LoadBalancerPool) String() string {
	return l.Name
}

func (l LoadBalancerPool) Keys() []string {
	return []string{
		fmt.Sprint(l.ID), l.Name, l.NameWithoutSpaces(), fmt.Sprint(l.EntryPort),
		l.Status.Name, l.Status.Key,
		l.EntryProtocol.Name, l.EntryProtocol.Key,
		l.TargetProtocol.Name, l.TargetProtocol.Key,
		l.Algorithm.Name, l.Algorithm.Key,
	}
}

func (l LoadBalancerPool) Columns() []string {
	return []string{"id", "name", "status", "entry protocol", "entry port", "target protocol", "algorithm", "sticky session"}
}

func (l LoadBalancerPool) Values() map[string]any {
	return map[string]any{
		"id":              l.ID,
		"name":            l.Name,
		"status":          l.Status.Name,
		"entry protocol":  l.EntryProtocol.Name,
		"entry port":      l.EntryPort,
		"target protocol": l.TargetProtocol.Name,
		"algorithm":       l.Algorithm.Name,
		"sticky session":  l.StickySession,
	}
}
