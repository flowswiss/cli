package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func SnapshotCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "snapshot",
		Aliases: []string{"snapshots"},
		Short:   "Manage compute snapshots",
	}

	commands.Add(app, cmd,
		&snapshotListCommand{},
		&snapshotCreateCommand{},
		&snapshotUpdateCommand{},
		&snapshotDeleteCommand{},
	)

	return cmd
}

type snapshotListCommand struct {
	filter string
}

func (s *snapshotListCommand) Run(cmd *cobra.Command, args []string) error {
	snapshots, err := compute.NewSnapshotService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch snapshots: %w", err)
	}

	if len(s.filter) != 0 {
		snapshots = filter.Find(snapshots, s.filter)
	}

	return commands.PrintStdout(snapshots)
}

func (s *snapshotListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *snapshotListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List snapshots",
		Long:              "Lists all compute snapshots.",
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type snapshotCreateCommand struct {
	name   string
	volume string
}

func (s *snapshotCreateCommand) Run(cmd *cobra.Command, args []string) error {
	volume, err := findVolume(cmd.Context(), s.volume)
	if err != nil {
		return err
	}

	data := compute.SnapshotCreate{
		Name:     s.name,
		VolumeID: volume.ID,
	}

	snapshot, err := compute.NewSnapshotService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}

	return commands.PrintStdout(snapshot)
}

func (s *snapshotCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *snapshotCreateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Aliases:           []string{"add", "new"},
		Short:             "Create a new snapshot",
		Long:              "Creates a new compute snapshot.",
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.name, "name", "", "name of the snapshot")
	cmd.Flags().StringVar(&s.volume, "volume", "", "volume to create the snapshot from")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("volume")

	return cmd
}

type snapshotUpdateCommand struct {
	name string
}

func (s *snapshotUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	snapshot, err := findSnapshot(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	data := compute.SnapshotUpdate{
		Name: s.name,
	}

	snapshot, err = compute.NewSnapshotService(commands.Config.Client).Update(cmd.Context(), snapshot.ID, data)
	if err != nil {
		return fmt.Errorf("update snapshot: %w", err)
	}

	return commands.PrintStdout(snapshot)
}

func (s *snapshotUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeSnapshot(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *snapshotUpdateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update SNAPSHOT",
		Short:             "Update snapshot",
		Long:              "Updates a snapshot.",
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.name, "name", "", "name of the snapshot")

	return cmd
}

type snapshotDeleteCommand struct {
	force bool
}

func (s *snapshotDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	snapshot, err := findSnapshot(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !s.force && !commands.ConfirmDeletion("snapshot", snapshot) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NewSnapshotService(commands.Config.Client).Delete(cmd.Context(), snapshot.ID)
	if err != nil {
		return fmt.Errorf("delete snapshot: %w", err)
	}

	return nil
}

func (s *snapshotDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeSnapshot(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *snapshotDeleteCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete SNAPSHOT",
		Short:             "Delete a snapshot",
		Long:              "Deletes a snapshot.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().BoolVar(&s.force, "force", false, "force the deletion of the snapshot without asking for confirmation")

	return cmd
}

func completeSnapshot(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	snapshots, err := compute.NewSnapshotService(commands.Config.Client).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(snapshots, term)

	names := make([]string, len(filtered))
	for i, snapshot := range filtered {
		names[i] = snapshot.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

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
