package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func VolumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"volumes"},
		Short:   "Manage compute volumes",
	}

	commands.Add(cmd,
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
	volumes, err := compute.NewVolumeService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch volumes: %w", err)
	}

	if len(v.filter) != 0 {
		volumes = filter.Find(volumes, v.filter)
	}

	return commands.PrintStdout(volumes)
}

func (v *volumeListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List volumes",
		Long:    "Lists all compute volumes.",
		RunE:    v.Run,
	}

	cmd.Flags().StringVar(&v.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type volumeCreateCommand struct {
	name     string
	size     int
	location string
	server   string
	snapshot string
}

func (v *volumeCreateCommand) Run(cmd *cobra.Command, args []string) error {
	data := compute.VolumeCreate{
		Name: v.name,
		Size: v.size,
	}

	if len(v.location) != 0 {
		location, err := common.FindLocation(cmd.Context(), commands.Config.Client, v.location)
		if err != nil {
			return err
		}

		data.LocationID = location.ID
	}

	if len(v.server) != 0 {
		server, err := findServer(cmd.Context(), v.server)
		if err != nil {
			return err
		}

		data.InstanceID = server.ID

		if data.LocationID == 0 {
			data.LocationID = server.Location.ID
		}
	}

	if len(v.snapshot) != 0 {
		snapshot, err := findSnapshot(cmd.Context(), v.snapshot)
		if err != nil {
			return err
		}

		data.SnapshotID = snapshot.ID

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

	volume, err := compute.NewVolumeService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create volume: %w", err)
	}

	return commands.PrintStdout(volume)
}

func (v *volumeCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "new"},
		Short:   "Create a new volume",
		Long:    "Creates a new compute volume.",
		RunE:    v.Run,
	}

	cmd.Flags().StringVar(&v.name, "name", "", "name of the volume")
	cmd.Flags().IntVar(&v.size, "size", 0, "size of the volume in GiB")
	cmd.Flags().StringVar(&v.location, "location", "", "location of the volume")
	cmd.Flags().StringVar(&v.server, "attach-to", "", "server to attach the volume to")
	cmd.Flags().StringVar(&v.snapshot, "restore-from", "", "snapshot to create the volume from")

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
		InstanceID: server.ID,
	}

	volume, err = compute.NewVolumeService(commands.Config.Client).Attach(cmd.Context(), volume.ID, data)
	if err != nil {
		return fmt.Errorf("attach volume: %w", err)
	}

	return commands.PrintStdout(volume)
}

func (v *volumeAttachCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach VOLUME SERVER",
		Short: "Attach a volume to a server",
		Long:  "Attaches a volume to a server.",
		Args:  cobra.ExactArgs(2),
		RunE:  v.Run,
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

	err = compute.NewVolumeService(commands.Config.Client).Detach(cmd.Context(), volume.ID, server.ID)
	if err != nil {
		return fmt.Errorf("detach volume: %w", err)
	}

	return nil
}

func (v *volumeDetachCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach VOLUME SERVER",
		Short: "Detach a volume from a server",
		Long:  "Detaches a volume from a server.",
		Args:  cobra.ExactArgs(2),
		RunE:  v.Run,
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

	data := compute.VolumeRevert{
		SnapshotID: snapshot.ID,
	}

	volume, err = compute.NewVolumeService(commands.Config.Client).Revert(cmd.Context(), volume.ID, data)
	if err != nil {
		return fmt.Errorf("revert volume: %w", err)
	}

	return commands.PrintStdout(volume)
}

func (v *volumeRevertCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revert VOLUME SNAPSHOT",
		Short: "Revert a volume to a snapshot",
		Long:  "Reverts a volume to a snapshot.",
		Args:  cobra.ExactArgs(2),
		RunE:  v.Run,
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

	data := compute.VolumeExpand{
		Size: v.size,
	}

	volume, err = compute.NewVolumeService(commands.Config.Client).Expand(cmd.Context(), volume.ID, data)
	if err != nil {
		return fmt.Errorf("expand volume: %w", err)
	}

	return commands.PrintStdout(volume)
}

func (v *volumeExpandCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expand VOLUME",
		Short: "Expand a volume",
		Long:  "Expands a volume.",
		Args:  cobra.ExactArgs(1),
		RunE:  v.Run,
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

	err = compute.NewVolumeService(commands.Config.Client).Delete(cmd.Context(), volume.ID)
	if err != nil {
		return fmt.Errorf("delete volume: %w", err)
	}

	return nil
}

func (v *volumeDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete VOLUME",
		Short: "Delete a volume",
		Long:  "Deletes a volume.",
		Args:  cobra.ExactArgs(1),
		RunE:  v.Run,
	}

	cmd.Flags().BoolVar(&v.force, "force", false, "force the deletion of the volume without asking for confirmation")

	return cmd
}

func findVolume(ctx context.Context, term string) (compute.Volume, error) {
	volumes, err := compute.NewVolumeService(commands.Config.Client).List(ctx)
	if err != nil {
		return compute.Volume{}, fmt.Errorf("fetch volumes: %w", err)
	}

	volume, err := filter.FindOne(volumes, term)
	if err != nil {
		return compute.Volume{}, fmt.Errorf("find volume: %w", err)
	}

	return volume, nil
}
