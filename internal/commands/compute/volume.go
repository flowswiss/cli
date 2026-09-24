package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/pkg/optional"
	"github.com/spf13/cobra"

	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func VolumeCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"volumes"},
		Short:   "Manage compute volumes",
	}

	commands.Add(app, cmd,
		&volumeListCommand{},
		&volumeCreateCommand{},
		&volumeAttachCommand{},
		&volumeDetachCommand{},
		&volumeRevertCommand{},
		&volumeExpandCommand{},
		&volumeDeleteCommand{},
	)

	return cmd
}

type volumeListCommand struct {
	filter string
}

func (v *volumeListCommand) Run(cmd *cobra.Command, args []string) error {
	volumes, err := compute.VolumeService().List(cmd.Context(), core.CursorAll)
	if err != nil {
		return fmt.Errorf("fetch volumes: %w", err)
	}

	if len(v.filter) != 0 {
		volumes = filter.Find(volumes, v.filter)
	}

	return commands.PrintStdout(volumes)
}

func (v *volumeListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (v *volumeListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List volumes",
		Long:              "Lists all compute volumes.",
		ValidArgsFunction: v.CompleteArg,
		RunE:              v.Run,
	}

	cmd.Flags().StringVar(&v.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type volumeCreateCommand struct {
	name     string
	size     int
	location string
	server   optional.Optional[string]
	snapshot optional.Optional[string]
}

func (v *volumeCreateCommand) Run(cmd *cobra.Command, args []string) error {
	data := compute.VolumeCreate{
		Name: v.name,
		Size: v.size,
	}

	if len(v.location) != 0 {
		location, err := common.FindLocation(cmd.Context(), commands.Client, v.location)
		if err != nil {
			return err
		}

		data.LocationID = location.ID
	}

	if srv := v.server.Value(); srv != nil {
		server, err := findServer(cmd.Context(), *srv)
		if err != nil {
			return err
		}

		data.InstanceID = &server.ID

		if data.LocationID == 0 {
			data.LocationID = server.Location.ID
		}
	}

	if snap := v.snapshot.Value(); snap != nil {
		snapshot, err := findSnapshot(cmd.Context(), *snap)
		if err != nil {
			return err
		}

		data.SnapshotID = &snapshot.ID

		if data.Size == 0 {
			data.Size = snapshot.Size
		}

		if data.LocationID == 0 {
			data.LocationID = snapshot.Volume.Location.ID
		}
	}

	if data.LocationID == 0 {
		return fmt.Errorf("unable to determine location for the volume")
	}

	volume, err := compute.VolumeService().Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create volume: %w", err)
	}

	return commands.PrintStdout(volume)
}

func (v *volumeCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (v *volumeCreateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Aliases:           []string{"add", "new"},
		Short:             "Create a new volume",
		Long:              "Creates a new compute volume.",
		ValidArgsFunction: v.CompleteArg,
		RunE:              v.Run,
	}

	cmd.Flags().StringVar(&v.name, "name", "", "name of the volume")
	cmd.Flags().IntVar(&v.size, "size", 0, "size of the volume in GiB")
	cmd.Flags().StringVar(&v.location, "location", "", "location of the volume")
	cmd.Flags().StringVar(v.server.Configure(cmd, "attach-to", "server to attach the volume to"))
	cmd.Flags().StringVar(v.snapshot.Configure(cmd, "restore-from", "snapshot to create the volume from"))

	_ = cmd.MarkFlagRequired("name")

	return cmd
}

type volumeAttachCommand struct {
}

func (v *volumeAttachCommand) Run(cmd *cobra.Command, args []string) error {
	volume, err := findVolume(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	server, err := findServer(cmd.Context(), args[1])
	if err != nil {
		return err
	}

	data := compute.VolumeAttach{
		VolumeID:   uint(volume.ID),
		InstanceID: server.ID,
	}

	volume, err = compute.VolumeService().Attach(cmd.Context(), data)

	if err != nil {
		return fmt.Errorf("attach volume: %w", err)
	}

	return commands.PrintStdout(volume)
}

func (v *volumeAttachCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeVolume(cmd.Context(), toComplete, func(volume compute.Volume) bool {
			return volume.AttachedTo.ID == 0
		})
	}

	if len(args) == 1 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (v *volumeAttachCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "attach VOLUME SERVER",
		Short:             "Attach a volume to a server",
		Long:              "Attaches a volume to a server.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: v.CompleteArg,
		RunE:              v.Run,
	}

	return cmd
}

type volumeDetachCommand struct {
	force bool
}

func (v *volumeDetachCommand) Run(cmd *cobra.Command, args []string) error {
	volume, err := findVolume(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	server, err := findServer(cmd.Context(), args[1])
	if err != nil {
		return err
	}

	if volume.AttachedTo.ID != server.ID {
		return fmt.Errorf("volume is not attached to the server")
	}

	if !v.force && !commands.Confirm(fmt.Sprintf("are you sure you want to detach volume %q from server %q?", volume, server)) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.VolumeService().Detach(cmd.Context(), compute.VolumeDetach{VolumeID: uint(volume.ID), InstanceID: uint(server.ID)})
	if err != nil {
		return fmt.Errorf("detach volume: %w", err)
	}

	return nil
}

func (v *volumeDetachCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeVolume(cmd.Context(), toComplete, func(volume compute.Volume) bool {
			return volume.AttachedTo.ID != 0
		})
	}

	if len(args) == 1 {
		volume, err := findVolume(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return []string{volume.AttachedTo.Name}, cobra.ShellCompDirectiveNoFileComp
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (v *volumeDetachCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "detach VOLUME SERVER",
		Short:             "Detach a volume from a server",
		Long:              "Detaches a volume from a server.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: v.CompleteArg,
		RunE:              v.Run,
	}

	cmd.Flags().BoolVar(&v.force, "force", false, "force detaching the volume without asking for confirmation")

	return cmd
}

type volumeRevertCommand struct {
}

func (v *volumeRevertCommand) Run(cmd *cobra.Command, args []string) error {
	volume, err := findVolume(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	snapshot, err := findSnapshot(cmd.Context(), args[1])
	if err != nil {
		return err
	}

	data := compute.VolumeRevert{VolumeID: uint(volume.ID),
		SnapshotID: snapshot.ID,
	}

	volume, err = compute.VolumeService().Revert(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("revert volume: %w", err)
	}

	return commands.PrintStdout(volume)
}

func (v *volumeRevertCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeVolume(cmd.Context(), toComplete, nil)
	}

	if len(args) == 1 {
		volume, err := findVolume(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeVolumeSnapshot(cmd.Context(), volume, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (v *volumeRevertCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "revert VOLUME SNAPSHOT",
		Short:             "Revert a volume to a snapshot",
		Long:              "Reverts a volume to a snapshot.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: v.CompleteArg,
		RunE:              v.Run,
	}

	return cmd
}

type volumeExpandCommand struct {
	size int
}

func (v *volumeExpandCommand) Run(cmd *cobra.Command, args []string) error {
	volume, err := findVolume(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	data := compute.VolumeExpand{VolumeID: uint(volume.ID),
		Size: v.size,
	}

	volume, err = compute.VolumeService().Expand(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("expand volume: %w", err)
	}

	return commands.PrintStdout(volume)
}

func (v *volumeExpandCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeVolume(cmd.Context(), toComplete, nil)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (v *volumeExpandCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "expand VOLUME",
		Short:             "Expand a volume",
		Long:              "Expands a volume.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: v.CompleteArg,
		RunE:              v.Run,
	}

	cmd.Flags().IntVar(&v.size, "size", 0, "size of the volume in GiB")

	_ = cmd.MarkFlagRequired("size")

	return cmd
}

type volumeDeleteCommand struct {
	force bool
}

func (v *volumeDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	volume, err := findVolume(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !v.force && !commands.ConfirmDeletion("volume", volume) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.VolumeService().Delete(cmd.Context(), compute.VolumeDelete{ID: uint(volume.ID)})
	if err != nil {
		return fmt.Errorf("delete volume: %w", err)
	}

	return nil
}

func (v *volumeDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeVolume(cmd.Context(), toComplete, nil)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (v *volumeDeleteCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete VOLUME",
		Short:             "Delete a volume",
		Long:              "Deletes a volume.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: v.CompleteArg,
		RunE:              v.Run,
	}

	cmd.Flags().BoolVar(&v.force, "force", false, "force the deletion of the volume without asking for confirmation")

	return cmd
}

func completeVolume(ctx context.Context, term string, itemFilter func(volume compute.Volume) bool) (
	[]string,
	cobra.ShellCompDirective,
) {
	volumes, err := compute.VolumeService().List(ctx, core.CursorAll)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.FindWithCustomFilter(volumes, term, itemFilter)

	names := make([]string, len(filtered))
	for i, volume := range filtered {
		names[i] = volume.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func completeVolumeSnapshot(ctx context.Context, volume compute.Volume, term string) (
	[]string,
	cobra.ShellCompDirective,
) {
	snapshots, err := compute.SnapshotService().List(ctx, core.CursorAll)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.FindWithCustomFilter(snapshots, term, func(snapshot compute.Snapshot) bool {
		return snapshot.Volume.ID == volume.ID
	})

	names := make([]string, len(filtered))
	for i, snapshot := range filtered {
		names[i] = snapshot.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findVolume(ctx context.Context, term string) (compute.Volume, error) {
	volumes, err := compute.VolumeService().List(ctx, core.CursorAll)
	if err != nil {
		return compute.Volume{}, fmt.Errorf("fetch volumes: %w", err)
	}

	volume, err := filter.FindOne(volumes, term)
	if err != nil {
		return compute.Volume{}, fmt.Errorf("find volume: %w", err)
	}

	return volume, nil
}
