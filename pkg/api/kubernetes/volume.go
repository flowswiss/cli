package kubernetes

import (
	"context"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/kubernetes"

	"github.com/flowswiss/cli/v2/pkg/api/compute"
)

type Volume = compute.Volume

type VolumeService struct {
	delegate kubernetes.VolumeService
}

func NewVolumeService(client goclient.Client, clusterID int) VolumeService {
	return VolumeService{
		delegate: kubernetes.NewVolumeService(client, clusterID),
	}
}

func (v VolumeService) List(ctx context.Context) ([]Volume, error) {
	res, err := v.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Volume, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Volume(item)
	}

	return items, nil
}

func (v VolumeService) Delete(ctx context.Context, id int) error {
	return v.delegate.Delete(ctx, id)
}
