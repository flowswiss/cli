package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func findSnapshot(ctx context.Context, term string) (compute.Snapshot, error) {
	snapshots, err := compute.NewSnapshotService(commands.Config.Client).List(ctx)
	if err != nil {
		return compute.Snapshot{}, fmt.Errorf("fetch snapshots: %w", err)
	}

	snapshot, err := filter.FindOne(snapshots, term)
	if err != nil {
		return compute.Snapshot{}, fmt.Errorf("find snapshot: %w", err)
	}

	return snapshot, nil
}
