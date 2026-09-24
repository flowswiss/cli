package compute

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
	"github.com/flowswiss/goclient/v2/core"
)

type (
	SnapshotCreate = compute.SnapshotCreateReq
	SnapshotGet    = compute.SnapshotGetReq
	SnapshotList   = core.Cursor
	SnapshotUpdate = compute.SnapshotUpdateReq
	SnapshotDelete = compute.SnapshotDeleteReq
)

type GenericSnapshotService = generic.CRUD[
	compute.Snapshot,
	Snapshot,
	SnapshotCreate,
	SnapshotGet,
	SnapshotList,
	SnapshotUpdate,
	SnapshotDelete,
]

func SnapshotService() GenericSnapshotService {
	return generic.NewCRUD[
		compute.Snapshot,
		Snapshot,
		SnapshotCreate,
		SnapshotGet,
		SnapshotList,
		SnapshotUpdate,
		SnapshotDelete,
	](commands.Client.Compute.Snapshot, func(snapshot compute.Snapshot) Snapshot {
		return Snapshot(snapshot)
	})
}

type Snapshot compute.Snapshot

func (s Snapshot) String() string {
	return s.Name
}

func (s Snapshot) Keys() []string {
	return []string{fmt.Sprint(s.ID), s.Name, s.Status.Key, s.Status.Name, s.Volume.Name}
}

func (s Snapshot) Columns() []string {
	return []string{"id", "name", "volume", "status", "size"}
}

func (s Snapshot) Values() map[string]any {
	return map[string]any{
		"id":     s.ID,
		"name":   s.Name,
		"volume": Volume(s.Volume),
		"status": s.Status.Name,
		"size":   fmt.Sprint(s.Size, " GiB"),
	}
}
