package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

type LoadBalancerPool compute.LoadBalancerPool

func (l LoadBalancerPool) String() string {
	return l.Name
}

func (l LoadBalancerPool) Keys() []string {
	return []string{
		fmt.Sprint(l.ID), l.Name, fmt.Sprint(l.EntryPort),
		l.Status.Name, l.Status.Key,
		l.EntryProtocol.Name, l.EntryProtocol.Key,
		l.TargetProtocol.Name, l.TargetProtocol.Key,
		l.Algorithm.Name, l.Algorithm.Key,
	}
}

func (l LoadBalancerPool) Columns() []string {
	return []string{"id", "name", "status", "entry protocol", "entry port", "target protocol", "algorithm", "sticky session"}
}

func (l LoadBalancerPool) Values() map[string]interface{} {
	return map[string]interface{}{
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

type LoadBalancerPoolService struct {
	delegate compute.LoadBalancerPoolService
}

func NewLoadBalancerPoolService(client goclient.Client, loadBalancerID int) LoadBalancerPoolService {
	return LoadBalancerPoolService{
		delegate: compute.NewLoadBalancerPoolService(client, loadBalancerID),
	}
}

func (l LoadBalancerPoolService) List(ctx context.Context) ([]LoadBalancerPool, error) {
	res, err := l.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]LoadBalancerPool, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = LoadBalancerPool(item)
	}

	return items, nil
}

type LoadBalancerHealthCheckOptions = compute.LoadBalancerHealthCheckOptions
type LoadBalancerPoolCreate = compute.LoadBalancerPoolCreate

func (l LoadBalancerPoolService) Create(ctx context.Context, data LoadBalancerPoolCreate) (LoadBalancerPool, error) {
	res, err := l.delegate.Create(ctx, data)
	if err != nil {
		return LoadBalancerPool{}, err
	}

	return LoadBalancerPool(res), nil
}

type LoadBalancerPoolUpdate = compute.LoadBalancerPoolUpdate

func (l LoadBalancerPoolService) Update(ctx context.Context, id int, data LoadBalancerPoolUpdate) (LoadBalancerPool, error) {
	res, err := l.delegate.Update(ctx, id, data)
	if err != nil {
		return LoadBalancerPool{}, err
	}

	return LoadBalancerPool(res), nil
}

func (l LoadBalancerPoolService) Delete(ctx context.Context, id int) error {
	return l.delegate.Delete(ctx, id)
}
