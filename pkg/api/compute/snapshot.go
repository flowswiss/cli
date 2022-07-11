package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

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

func (s Snapshot) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":     s.ID,
		"name":   s.Name,
		"volume": Volume(s.Volume),
		"status": s.Status.Name,
		"size":   fmt.Sprint(s.Size, " GiB"),
	}
}

type SnapshotService struct {
	delegate compute.SnapshotService
}

func NewSnapshotService(client goclient.Client) SnapshotService {
	return SnapshotService{
		delegate: compute.NewSnapshotService(client),
	}
}

func (v SnapshotService) List(ctx context.Context) ([]Snapshot, error) {
	res, err := v.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Snapshot, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Snapshot(item)
	}

	return items, nil
}

type SnapshotCreate = compute.SnapshotCreate

func (v SnapshotService) Create(ctx context.Context, data SnapshotCreate) (Snapshot, error) {
	volume, err := v.delegate.Create(ctx, data)
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot(volume), nil
}

type SnapshotUpdate = compute.SnapshotUpdate

func (v SnapshotService) Update(ctx context.Context, volumeID int, data SnapshotUpdate) (Snapshot, error) {
	volume, err := v.delegate.Update(ctx, volumeID, data)
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot(volume), nil
}

func (v SnapshotService) Delete(ctx context.Context, volumeID int) error {
	return v.delegate.Delete(ctx, volumeID)
}
