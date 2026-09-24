package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	VolumeCreate = compute.VolumeCreateReq
	VolumeGet    = compute.VolumeGetReq
	VolumeList   = core.Cursor
	VolumeUpdate = compute.VolumeUpdateReq
	VolumeAttach = compute.VolumeAttachReq
	VolumeDetach = compute.VolumeDetachReq
	VolumeRevert = compute.VolumeRevertReq
	VolumeExpand = compute.VolumeExpandReq
	VolumeDelete = compute.VolumeDeleteReq
)

type GenericVolumeService struct {
	generic.CreateService[VolumeCreate, compute.Volume, Volume]
	generic.Read[compute.Volume, Volume, VolumeGet, VolumeList]
	generic.UpdateService[VolumeUpdate, compute.Volume, Volume]
	generic.DeleteService[VolumeDelete]

	client *compute.VolumeService
}

func (v GenericVolumeService) Attach(ctx context.Context, attach VolumeAttach) (Volume, error) {
	volume, err := v.client.Attach(ctx, attach)
	if err != nil {
		return Volume{}, err
	}

	return Volume(volume), err
}

func (v GenericVolumeService) Detach(ctx context.Context, detach VolumeDetach) error {
	err := v.client.Detach(ctx, detach)
	if err != nil {
		return err
	}

	return err
}

func (v GenericVolumeService) Revert(ctx context.Context, revert VolumeRevert) (Volume, error) {
	volume, err := v.client.Revert(ctx, revert)
	if err != nil {
		return Volume{}, err
	}

	return Volume(volume), err
}

func (v GenericVolumeService) Expand(ctx context.Context, expand VolumeExpand) (Volume, error) {
	volume, err := v.client.Expand(ctx, expand)
	if err != nil {
		return Volume{}, err
	}

	return Volume(volume), err
}

func VolumeService() GenericVolumeService {
	client := commands.Client.Compute.Volume
	cast := func(volume compute.Volume) Volume {
		return Volume(volume)
	}

	return GenericVolumeService{
		generic.NewCreate[VolumeCreate, compute.Volume, Volume](client, cast),
		generic.NewRead[compute.Volume, Volume, VolumeGet, VolumeList](client, cast),
		generic.NewUpdate[VolumeUpdate, compute.Volume, Volume](client, cast),
		generic.NewDelete[VolumeDelete](client),
		client,
	}
}

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

func (v Volume) Values() map[string]any {
	return map[string]any{
		"id":          v.ID,
		"name":        v.Name,
		"location":    common.Location(v.Location),
		"status":      v.Status.Name,
		"size":        fmt.Sprint(v.Size, " GiB"),
		"attached to": Server(v.AttachedTo),
	}
}
