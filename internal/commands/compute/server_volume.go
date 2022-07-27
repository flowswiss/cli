package compute

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ServerVolumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"volumes"},
		Short:   "Manage server volumes",
	}

	commands.Add(cmd,
		&serverVolumeListCommand{},
		&serverVolumeCreateCommand{},
		&serverVolumeAttachCommand{},
		&serverVolumeDetachCommand{},
		&serverVolumeDeleteCommand{},
	)

	return cmd
}

type serverVolumeListCommand struct {
	filter string
}

func (s *serverVolumeListCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	items, err := compute.NewVolumeService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch volumes: %w", err)
	}

	volumes := make([]compute.Volume, 0, len(items))
	for _, item := range items {
		if item.AttachedTo.ID == server.ID {
			volumes = append(volumes, item)
		}
	}

	if len(s.filter) != 0 {
		volumes = filter.Find(volumes, s.filter)
	}

	return commands.PrintStdout(volumes)
}

func (s *serverVolumeListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverVolumeListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list SERVER",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List volumes attached to server",
		Long:              "Lists all compute volumes attached to the selected server.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type serverVolumeCreateCommand struct {
	name     string
	size     int
	snapshot string
}

func (s *serverVolumeCreateCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	data := compute.VolumeCreate{
		Name:       s.name,
		Size:       s.size,
		LocationID: server.Location.ID,
		InstanceID: server.ID,
	}

	if len(s.snapshot) != 0 {
		snapshot, err := findSnapshot(cmd.Context(), s.snapshot)
		if err != nil {
			return err
		}

		data.SnapshotID = snapshot.ID

		if data.Size == 0 {
			data.Size = snapshot.Size
		}
	}

	volume, err := compute.NewVolumeService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create volume: %w", err)
	}

	return commands.PrintStdout(volume)
}

func (s *serverVolumeCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverVolumeCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create SERVER",
		Aliases:           []string{"add", "new"},
		Short:             "Create a new volume",
		Long:              "Creates a new compute volume.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.name, "name", "", "name of the volume")
	cmd.Flags().IntVar(&s.size, "size", 0, "size of the volume in GiB")
	cmd.Flags().StringVar(&s.snapshot, "restore-from", "", "snapshot to create the volume from")

	_ = cmd.MarkFlagRequired("name")

	return cmd
}

type serverVolumeAttachCommand struct {
}

func (s *serverVolumeAttachCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	volume, err := findVolume(cmd.Context(), args[1])
	if err != nil {
		return err
	}

	if volume.AttachedTo.ID != server.ID {
		// already attached to this server
		return nil
	}

	if volume.AttachedTo.ID != 0 {
		return fmt.Errorf("volume is already attached to server %q", compute.Server(volume.AttachedTo))
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

func (s *serverVolumeAttachCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		return completeVolume(cmd.Context(), toComplete, func(item compute.Volume) bool {
			return item.AttachedTo.ID == 0
		})
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverVolumeAttachCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "attach SERVER VOLUME",
		Short:             "Attach a volume to a server",
		Long:              "Attaches a volume to a server.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	return cmd
}

type serverVolumeDetachCommand struct {
	force bool
}

func (s *serverVolumeDetachCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	volume, err := findVolume(cmd.Context(), args[1])
	if err != nil {
		return err
	}

	if volume.AttachedTo.ID != server.ID {
		return fmt.Errorf("volume is not attached to the server")
	}

	if !s.force && !commands.Confirm(fmt.Sprintf("are you sure you want to detach volume %q from server %q?", volume, server)) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NewVolumeService(commands.Config.Client).Detach(cmd.Context(), volume.ID, server.ID)
	if err != nil {
		return fmt.Errorf("detach volume: %w", err)
	}

	return nil
}

func (s *serverVolumeDetachCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		server, err := findServer(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeVolume(cmd.Context(), toComplete, func(item compute.Volume) bool {
			return item.AttachedTo.ID == server.ID
		})
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverVolumeDetachCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "detach SERVER VOLUME",
		Short:             "Detach a volume from a server",
		Long:              "Detaches a volume from a server.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().BoolVar(&s.force, "force", false, "force detaching the volume without asking for confirmation")

	return cmd
}

type serverVolumeDeleteCommand struct {
	force bool
}

func (s *serverVolumeDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	volume, err := findVolume(cmd.Context(), args[1])
	if err != nil {
		return err
	}

	if volume.AttachedTo.ID != server.ID {
		return fmt.Errorf("volume is not attached to the server")
	}

	if !s.force && !commands.ConfirmDeletion("volume", volume) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	service := compute.NewVolumeService(commands.Config.Client)

	err = service.Detach(cmd.Context(), volume.ID, server.ID)
	if err != nil {
		return fmt.Errorf("detach volume: %w", err)
	}

	err = service.Delete(cmd.Context(), volume.ID)
	if err != nil {
		return fmt.Errorf("delete volume: %w", err)
	}

	return nil
}

func (s *serverVolumeDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		server, err := findServer(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeVolume(cmd.Context(), toComplete, func(item compute.Volume) bool {
			return item.AttachedTo.ID == server.ID
		})
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverVolumeDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete SERVER VOLUME",
		Short:             "Delete a volume from a server",
		Long:              "Deletes a volume currently attached to a server.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().BoolVar(&s.force, "force", false, "force the deletion of the volume without asking for confirmation")

	return cmd
}
