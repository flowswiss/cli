package kubernetes

import (
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/flowswiss/goclient/v2/kubernetes"

	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
)

type (
	ClusterCreate    = kubernetes.ClusterCreateReq
	ClusterGet       = kubernetes.ClusterGetReq
	ClusterList      = core.Cursor
	ClusterUpdate    = kubernetes.ClusterUpdateReq
	ClusterRunAction = kubernetes.ClusterPerformActionReq
	ClusterDelete    = kubernetes.ClusterDeleteReq

	ClusterKubeConfigGet       = kubernetes.ClusterGetReq
	ClusterConfigurationGet    = kubernetes.ClusterGetReq
	ClusterConfigurationUpdate = kubernetes.ClusterConfigurationReq
	ClusterFlavorUpdate        = kubernetes.ClusterUpdateFlavorReq

	// Wrapped types to allow direct usage

	ClusterWorkerCreate = kubernetes.ClusterWorkerCreateReq
	ClusterWorkerUpdate = kubernetes.ClusterWorkerUpdateReq

	KubeConfig    = kubernetes.ClusterKubeConfig
	Configuration = kubernetes.ClusterConfiguration
)

type GenericClusterService struct {
	generic.OrderedCreateService[ClusterCreate]
	generic.Read[kubernetes.Cluster, Cluster, ClusterGet, ClusterList]
	generic.UpdateService[ClusterUpdate, kubernetes.Cluster, Cluster]
	generic.PerformActionService[ClusterRunAction, kubernetes.Cluster, Cluster]
	generic.DeleteService[ClusterDelete]

	client *kubernetes.ClusterService
}

func (c GenericClusterService) GetKubeConfig(ctx context.Context, get ClusterKubeConfigGet) (KubeConfig, error) {
	config, err := c.client.GetKubeConfig(ctx, get)
	if err != nil {
		return KubeConfig{}, err
	}

	return config, nil
}

func (c GenericClusterService) GetConfiguration(
	ctx context.Context,
	get ClusterConfigurationGet,
) (Configuration, error) {
	configuration, err := c.client.GetConfiguration(ctx, get)
	if err != nil {
		return Configuration{}, err
	}

	return configuration, nil
}

func (c GenericClusterService) UpdateConfiguration(
	ctx context.Context,
	update ClusterConfigurationUpdate,
) (Configuration, error) {
	configuration, err := c.client.UpdateConfiguration(ctx, update)
	if err != nil {
		return Configuration{}, err
	}

	return configuration, nil
}

func (c GenericClusterService) UpdateFlavor(ctx context.Context, update ClusterFlavorUpdate) (Cluster, error) {
	cluster, err := c.client.UpdateFlavor(ctx, update)
	if err != nil {
		return Cluster{}, err
	}

	return Cluster(cluster), nil
}

func ClusterService() GenericClusterService {
	client := commands.Client.Kubernetes.Cluster
	cast := func(cluster kubernetes.Cluster) Cluster {
		return Cluster(cluster)
	}

	return GenericClusterService{
		generic.NewOrderedCreate[ClusterCreate](client),
		generic.NewRead[kubernetes.Cluster, Cluster, ClusterGet, ClusterList](client, cast),
		generic.NewUpdate[ClusterUpdate, kubernetes.Cluster, Cluster](client, cast),
		generic.NewPerformAction[ClusterRunAction, kubernetes.Cluster, Cluster](client, cast),
		generic.NewDelete[ClusterDelete](client),
		client,
	}
}

type Cluster kubernetes.Cluster

func (c Cluster) String() string {
	return c.Name
}

func (c Cluster) Keys() []string {
	return []string{fmt.Sprint(c.ID), c.Name, c.DNSName, c.PublicAddress}
}

func (c Cluster) Columns() []string {
	return []string{"id", "name", "status", "product", "location", "network", "address", "control plane", "worker"}
}

func (c Cluster) Values() map[string]any {
	return map[string]any{
		"id":            c.ID,
		"name":          c.Name,
		"status":        c.Status.Name,
		"product":       common.Product(c.Product),
		"location":      common.Location(c.Location),
		"network":       compute.Network(c.Network),
		"address":       fmt.Sprintf("%s (%s)", c.DNSName, c.PublicAddress),
		"control plane": fmt.Sprintf("%d/%d (%s)", c.NodeCount.Current.ControlPlane, c.NodeCount.Expected.ControlPlane, c.ExpectedPreset.ControlPlane.Name),
		"worker":        fmt.Sprintf("%d/%d (%s)", c.NodeCount.Current.Worker, c.NodeCount.Expected.Worker, c.ExpectedPreset.Worker.Name),
	}
}

type ClusterAction kubernetes.ClusterAction

func (c ClusterAction) Keys() []string {
	return []string{fmt.Sprint(c.ID), c.Name, c.Command}
}

func (c ClusterAction) Columns() []string {
	return []string{"id", "name", "command"}
}

func (c ClusterAction) Values() map[string]any {
	return map[string]any{
		"id":      c.ID,
		"name":    c.Name,
		"command": c.Command,
	}
}
