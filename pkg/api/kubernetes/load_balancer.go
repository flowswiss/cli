package kubernetes

import (
	"context"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/kubernetes"

	"github.com/flowswiss/cli/v2/pkg/api/compute"
)

type LoadBalancer = compute.LoadBalancer

type LoadBalancerService struct {
	delegate kubernetes.LoadBalancerService
}

func NewLoadBalancerService(client goclient.Client, clusterID int) LoadBalancerService {
	return LoadBalancerService{
		delegate: kubernetes.NewLoadBalancerService(client, clusterID),
	}
}

func (v LoadBalancerService) List(ctx context.Context) ([]LoadBalancer, error) {
	res, err := v.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]LoadBalancer, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = LoadBalancer(item)
	}

	return items, nil
}
