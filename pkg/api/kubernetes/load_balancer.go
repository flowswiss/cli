package kubernetes

import (
	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/kubernetes"

	"github.com/flowswiss/cli/v2/pkg/api/compute"
)

type (
	LoadBalancer = compute.LoadBalancer

	LoadBalancerList = kubernetes.LoadBalancerListReq
)

type GenericLoadBalancerService struct {
	generic.ListService[LoadBalancerList, kubernetes.LoadBalancer, LoadBalancer]
}

func LoadBalancerService() GenericLoadBalancerService {
	client := commands.Client.Kubernetes.LoadBalancer
	cast := func(loadBalancer kubernetes.LoadBalancer) LoadBalancer {
		return LoadBalancer(loadBalancer)
	}

	return GenericLoadBalancerService{
		generic.NewList[LoadBalancerList, kubernetes.LoadBalancer, LoadBalancer](client, cast),
	}
}
