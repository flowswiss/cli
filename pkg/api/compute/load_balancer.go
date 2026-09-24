package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2"
	"github.com/flowswiss/goclient/v2/compute"
	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	LoadBalancerCreate    = compute.LoadBalancerCreateReq
	LoadBalancerGet       = compute.LoadBalancerGetReq
	LoadBalancerList      = core.Cursor
	LoadBalancerUpdate    = compute.LoadBalancerUpdateReq
	LoadBalancerRunAction = compute.LoadBalancerPerformReq
	LoadBalancerDelete    = compute.LoadBalancerDeleteReq
)

type GenericLoadBalancerService struct {
	generic.OrderedCreateService[LoadBalancerCreate]
	generic.Read[compute.LoadBalancer, LoadBalancer, LoadBalancerGet, LoadBalancerList]
	generic.UpdateService[LoadBalancerUpdate, compute.LoadBalancer, LoadBalancer]
	generic.PerformActionService[LoadBalancerRunAction, compute.LoadBalancer, LoadBalancer]
	generic.DeleteService[LoadBalancerDelete]
}

func LoadBalancerService() GenericLoadBalancerService {
	client := commands.Client.Compute.LoadBalancer
	cast := func(loadBalancer compute.LoadBalancer) LoadBalancer {
		return LoadBalancer(loadBalancer)
	}

	return GenericLoadBalancerService{
		generic.NewOrderedCreate[LoadBalancerCreate](client),
		generic.NewRead[compute.LoadBalancer, LoadBalancer, LoadBalancerGet, LoadBalancerList](client, cast),
		generic.NewUpdate[LoadBalancerUpdate, compute.LoadBalancer, LoadBalancer](client, cast),
		generic.NewPerformAction[LoadBalancerRunAction, compute.LoadBalancer, LoadBalancer](client, cast),
		generic.NewDelete[LoadBalancerDelete](client),
	}
}

type LoadBalancer compute.LoadBalancer

func (l LoadBalancer) String() string {
	return l.Name
}

func (l LoadBalancer) Keys() []string {
	return []string{fmt.Sprint(l.ID), l.Name}
}

func (l LoadBalancer) Columns() []string {
	return []string{"id", "name", "location", "product", "status", "public ip", "network"}
}

func (l LoadBalancer) Values() map[string]any {
	networkBuffer := &strings.Builder{}
	publicIPBuffer := &strings.Builder{}

	for i, network := range l.Networks {
		if i != 0 {
			networkBuffer.WriteString(", ")
		}

		networkBuffer.WriteString(fmt.Sprintf("%s (", network.Name))
		for j, iface := range network.Interfaces {
			if j != 0 {
				networkBuffer.WriteString(", ")
			}

			networkBuffer.WriteString(iface.PrivateIP)

			if iface.PublicIP != "" {
				publicIPBuffer.WriteString(fmt.Sprintf("%s, ", iface.PublicIP))
			}
		}
		networkBuffer.WriteRune(')')
	}

	publicIP := publicIPBuffer.String()
	if len(publicIP) > 0 {
		publicIP = publicIP[:len(publicIP)-2]
	}

	return map[string]any{
		"id":        l.ID,
		"name":      l.Name,
		"location":  common.Location(l.Location),
		"product":   common.Product(l.Product),
		"status":    l.Status.Name,
		"public ip": publicIP,
		"network":   networkBuffer.String(),
	}
}

type LoadBalancerProtocol compute.LoadBalancerProtocol

func (l LoadBalancerProtocol) String() string {
	return l.Name
}

func (l LoadBalancerProtocol) Keys() []string {
	return []string{fmt.Sprint(l.ID), l.Key, l.Name}
}

func (l LoadBalancerProtocol) Columns() []string {
	return []string{"id", "key", "name"}
}

func (l LoadBalancerProtocol) Values() map[string]any {
	return map[string]any{
		"id":   l.ID,
		"key":  l.Key,
		"name": l.Name,
	}
}

func LoadBalancerProtocols(ctx context.Context, client *goclient.Client) ([]LoadBalancerProtocol, error) {
	res, err := client.Compute.LoadBalancerEntity.ListProtocols(ctx, core.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]LoadBalancerProtocol, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = LoadBalancerProtocol(item)
	}

	return items, nil
}

type LoadBalancerAlgorithm compute.LoadBalancerAlgorithm

func (l LoadBalancerAlgorithm) String() string {
	return l.Name
}

func (l LoadBalancerAlgorithm) Keys() []string {
	return []string{fmt.Sprint(l.ID), l.Key, l.Name}
}

func (l LoadBalancerAlgorithm) Columns() []string {
	return []string{"id", "key", "name"}
}

func (l LoadBalancerAlgorithm) Values() map[string]any {
	return map[string]any{
		"id":   l.ID,
		"key":  l.Key,
		"name": l.Name,
	}
}

func LoadBalancerAlgorithms(ctx context.Context, client *goclient.Client) ([]LoadBalancerAlgorithm, error) {
	res, err := client.Compute.LoadBalancerEntity.ListAlgorithms(ctx, core.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]LoadBalancerAlgorithm, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = LoadBalancerAlgorithm(item)
	}

	return items, nil
}

type LoadBalancerHealthCheckType compute.LoadBalancerHealthCheckType

func (l LoadBalancerHealthCheckType) String() string {
	return l.Name
}

func (l LoadBalancerHealthCheckType) Keys() []string {
	return []string{fmt.Sprint(l.ID), l.Key, l.Name}
}

func (l LoadBalancerHealthCheckType) Columns() []string {
	return []string{"id", "key", "name"}
}

func (l LoadBalancerHealthCheckType) Values() map[string]any {
	return map[string]any{
		"id":   l.ID,
		"key":  l.Key,
		"name": l.Name,
	}
}

func LoadBalancerHealthCheckTypes(ctx context.Context, client *goclient.Client) ([]LoadBalancerHealthCheckType, error) {
	res, err := client.Compute.LoadBalancerEntity.ListHealthCheckTypes(ctx, core.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]LoadBalancerHealthCheckType, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = LoadBalancerHealthCheckType(item)
	}

	return items, nil
}
