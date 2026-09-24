package kubernetes

import (
	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/kubernetes"
)

type (
	Snapshot = compute.Snapshot

	SnapshotList = kubernetes.SnapshotListReq
)

type GenericSnapshotService struct {
	generic.ListService[SnapshotList, kubernetes.Snapshot, Snapshot]
}

func SnapshotService() GenericSnapshotService {
	client := commands.Client.Kubernetes.Snapshot
	cast := func(snapshot kubernetes.Snapshot) Snapshot {
		return Snapshot(snapshot)
	}

	return GenericSnapshotService{
		generic.NewList[SnapshotList, kubernetes.Snapshot, Snapshot](client, cast),
	}
}
