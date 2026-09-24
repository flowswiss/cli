package kubernetes

import (
	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/kubernetes"

	"github.com/flowswiss/cli/v2/pkg/api/compute"
)

type (
	Volume = compute.Volume

	VolumeList   = kubernetes.VolumeListReq
	VolumeDelete = kubernetes.VolumeDeleteReq
)

type GenericVolumeService struct {
	generic.ListService[VolumeList, kubernetes.Volume, Volume]
	generic.DeleteService[VolumeDelete]
}

func VolumeService() GenericVolumeService {
	client := commands.Client.Kubernetes.Volume
	cast := func(volume kubernetes.Volume) Volume {
		return Volume(volume)
	}

	return GenericVolumeService{
		generic.NewList[VolumeList, kubernetes.Volume, Volume](client, cast),
		generic.NewDelete[VolumeDelete](client),
	}
}
