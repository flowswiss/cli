package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type Volume compute.Volume

func (v Volume) String() string {
	return v.Name
}

func (v Volume) Keys() []string {
	return []string{fmt.Sprint(v.ID), v.Name, v.SerialNumber, v.Status.Key, v.Status.Name}
}

func (v Volume) Columns() []string {
	return []string{"id", "name", "location", "status", "size", "attached to"}
}

func (v Volume) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":          v.ID,
		"name":        v.Name,
		"location":    common.Location(v.Location),
		"status":      v.Status.Name,
		"size":        fmt.Sprint(v.Size, " GiB"),
		"attached to": Server(v.AttachedTo),
	}
}

type VolumeService struct {
	delegate compute.VolumeService
}

func NewVolumeService(client goclient.Client) VolumeService {
	return VolumeService{
		delegate: compute.NewVolumeService(client),
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

type VolumeCreate = compute.VolumeCreate

func (v VolumeService) Create(ctx context.Context, data VolumeCreate) (Volume, error) {
	volume, err := v.delegate.Create(ctx, data)
	if err != nil {
		return Volume{}, err
	}

	return Volume(volume), nil
}

type VolumeUpdate = compute.VolumeUpdate

func (v VolumeService) Update(ctx context.Context, volumeID int, data VolumeUpdate) (Volume, error) {
	volume, err := v.delegate.Update(ctx, volumeID, data)
	if err != nil {
		return Volume{}, err
	}

	return Volume(volume), nil
}

type VolumeAttach = compute.VolumeAttach

func (v VolumeService) Attach(ctx context.Context, volumeID int, data VolumeAttach) (Volume, error) {
	volume, err := v.delegate.Attach(ctx, volumeID, data)
	if err != nil {
		return Volume{}, err
	}

	return Volume(volume), nil
}

func (v VolumeService) Detach(ctx context.Context, volumeID, serverID int) error {
	return v.delegate.Detach(ctx, volumeID, serverID)
}

type VolumeRevert = compute.VolumeRevert

func (v VolumeService) Revert(ctx context.Context, volumeID int, data VolumeRevert) (Volume, error) {
	volume, err := v.delegate.Revert(ctx, volumeID, data)
	if err != nil {
		return Volume{}, err
	}

	return Volume(volume), nil
}

type VolumeExpand = compute.VolumeExpand

func (v VolumeService) Expand(ctx context.Context, volumeID int, data VolumeExpand) (Volume, error) {
	volume, err := v.delegate.Expand(ctx, volumeID, data)
	if err != nil {
		return Volume{}, err
	}

	return Volume(volume), nil
}

func (v VolumeService) Delete(ctx context.Context, volumeID int) error {
	return v.delegate.Delete(ctx, volumeID)
}
