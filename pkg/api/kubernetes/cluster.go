package kubernetes

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/kubernetes"

	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
)

type Cluster kubernetes.Cluster

func (c Cluster) String() string {
	return c.Name
}

func (c Cluster) Keys() (identifiers []string) {
	return []string{fmt.Sprint(c.ID), c.Name, c.DNSName, c.PublicAddress}
}

func (c Cluster) Columns() []string {
	return []string{"id", "name", "status", "product", "location", "network", "address", "control plane", "worker"}
}

func (c Cluster) Values() map[string]interface{} {
	return map[string]interface{}{
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

type ClusterService struct {
	delegate kubernetes.ClusterService
}

func NewClusterService(client goclient.Client) ClusterService {
	return ClusterService{
		delegate: kubernetes.NewClusterService(client),
	}
}

func (c ClusterService) List(ctx context.Context) ([]Cluster, error) {
	res, err := c.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Cluster, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Cluster(item)
	}

	return items, nil
}

type ClusterCreate = kubernetes.ClusterCreate
type ClusterWorkerCreate = kubernetes.ClusterWorkerCreate

func (c ClusterService) Create(ctx context.Context, data ClusterCreate) (common.Ordering, error) {
	return c.delegate.Create(ctx, data)
}

type ClusterUpdate = kubernetes.ClusterUpdate

func (c ClusterService) Update(ctx context.Context, id int, data ClusterUpdate) (Cluster, error) {
	res, err := c.delegate.Update(ctx, id, data)
	if err != nil {
		return Cluster{}, err
	}

	return Cluster(res), nil
}

func (c ClusterService) Delete(ctx context.Context, id int) error {
	return c.delegate.Delete(ctx, id)
}

type ClusterKubeConfig = kubernetes.ClusterKubeConfig

func (c ClusterService) GetKubeConfig(ctx context.Context, id int) (ClusterKubeConfig, error) {
	return c.delegate.GetKubeConfig(ctx, id)
}

type ClusterConfiguration = kubernetes.ClusterConfiguration

func (c ClusterService) GetConfiguration(ctx context.Context, id int) (ClusterConfiguration, error) {
	return c.delegate.GetConfiguration(ctx, id)
}

func (c ClusterService) UpdateConfiguration(ctx context.Context, id int, data ClusterConfiguration) (config ClusterConfiguration, err error) {
	return c.delegate.UpdateConfiguration(ctx, id, data)
}

type ClusterUpdateFlavor = kubernetes.ClusterUpdateFlavor

func (c ClusterService) UpdateFlavor(ctx context.Context, id int, data ClusterUpdateFlavor) (Cluster, error) {
	cluster, err := c.delegate.UpdateFlavor(ctx, id, data)
	if err != nil {
		return Cluster{}, err
	}

	return Cluster(cluster), nil
}

type ClusterPerformAction = kubernetes.ClusterPerformAction

func (c ClusterService) PerformAction(ctx context.Context, id int, data ClusterPerformAction) (Cluster, error) {
	cluster, err := c.delegate.PerformAction(ctx, id, data)
	if err != nil {
		return Cluster{}, err
	}

	return Cluster(cluster), nil
}