package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

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

func (l LoadBalancer) Values() map[string]interface{} {
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

	return map[string]interface{}{
		"id":        l.ID,
		"name":      l.Name,
		"location":  common.Location(l.Location),
		"product":   common.Product(l.Product),
		"status":    l.Status.Name,
		"public ip": publicIP,
		"network":   networkBuffer.String(),
	}
}

type LoadBalancerService struct {
	delegate compute.LoadBalancerService
}

func NewLoadBalancerService(client goclient.Client) LoadBalancerService {
	return LoadBalancerService{
		delegate: compute.NewLoadBalancerService(client),
	}
}

func (l LoadBalancerService) List(ctx context.Context) ([]LoadBalancer, error) {
	res, err := l.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]LoadBalancer, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = LoadBalancer(item)
	}

	return items, nil
}

func (l LoadBalancerService) Get(ctx context.Context, id int) (LoadBalancer, error) {
	loadBalancer, err := l.delegate.Get(ctx, id)
	return LoadBalancer(loadBalancer), err
}

type LoadBalancerCreate = compute.LoadBalancerCreate

func (l LoadBalancerService) Create(ctx context.Context, data LoadBalancerCreate) (common.Ordering, error) {
	res, err := l.delegate.Create(ctx, data)
	if err != nil {
		return common.Ordering{}, err
	}

	return res, nil
}

type LoadBalancerUpdate = compute.LoadBalancerUpdate

func (l LoadBalancerService) Update(ctx context.Context, id int, data LoadBalancerUpdate) (LoadBalancer, error) {
	res, err := l.delegate.Update(ctx, id, data)
	if err != nil {
		return LoadBalancer{}, err
	}

	return LoadBalancer(res), nil
}

func (l LoadBalancerService) Delete(ctx context.Context, id int) error {
	return l.delegate.Delete(ctx, id)
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

func (l LoadBalancerProtocol) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":   l.ID,
		"key":  l.Key,
		"name": l.Name,
	}
}

func LoadBalancerProtocols(ctx context.Context, client goclient.Client) ([]LoadBalancerProtocol, error) {
	res, err := compute.NewLoadBalancerEntityService(client).ListProtocols(ctx, goclient.Cursor{NoFilter: 1})
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

func (l LoadBalancerAlgorithm) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":   l.ID,
		"key":  l.Key,
		"name": l.Name,
	}
}

func LoadBalancerAlgorithms(ctx context.Context, client goclient.Client) ([]LoadBalancerAlgorithm, error) {
	res, err := compute.NewLoadBalancerEntityService(client).ListAlgorithms(ctx, goclient.Cursor{NoFilter: 1})
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

func (l LoadBalancerHealthCheckType) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":   l.ID,
		"key":  l.Key,
		"name": l.Name,
	}
}

func LoadBalancerHealthCheckTypes(ctx context.Context, client goclient.Client) ([]LoadBalancerHealthCheckType, error) {
	res, err := compute.NewLoadBalancerEntityService(client).ListHealthCheckTypes(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]LoadBalancerHealthCheckType, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = LoadBalancerHealthCheckType(item)
	}

	return items, nil
}
