package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/kubernetes"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func VolumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"volumes"},
		Short:   "Manage your cluster volumes",
	}

	commands.Add(cmd,
		&volumeListCommand{},
		&volumeDeleteCommand{},
	)

	return cmd
}

type volumeListCommand struct {
	filter string
}

func (v *volumeListCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	items, err := kubernetes.NewVolumeService(commands.Config.Client, cluster.ID).List(cmd.Context())
	if err != nil {
		return err
	}

	if len(v.filter) != 0 {
		items = filter.Find(items, v.filter)
	}

	return commands.PrintStdout(items)
}

func (v *volumeListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list CLUSTER",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List all volumes",
		Long:    "Prints a table of all volumes belonging to the selected cluster.",
		Args:    cobra.ExactArgs(1),
		RunE:    v.Run,
	}

	cmd.Flags().StringVar(&v.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type volumeDeleteCommand struct {
	force bool
}

func (v *volumeDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	volume, err := findVolume(cmd.Context(), cluster.ID, args[1])
	if err != nil {
		return err
	}

	if !v.force && !commands.ConfirmDeletion("volume", volume) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = kubernetes.NewVolumeService(commands.Config.Client, cluster.ID).Delete(cmd.Context(), volume.ID)
	if err != nil {
		return fmt.Errorf("delete volume: %w", err)
	}

	return nil
}

func (v *volumeDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete CLUSTER VOLUME",
		Short: "Delete volume",
		Long:  "Deletes a volume.",
		Args:  cobra.ExactArgs(2),
		RunE:  v.Run,
	}

	cmd.Flags().BoolVar(&v.force, "force", false, "forces deletion of the volume without asking for confirmation")

	return cmd
}

func findVolume(ctx context.Context, clusterID int, term string) (kubernetes.Volume, error) {
	volumes, err := kubernetes.NewVolumeService(commands.Config.Client, clusterID).List(ctx)
	if err != nil {
		return kubernetes.Volume{}, fmt.Errorf("fetch volumes: %w", err)
	}

	volume, err := filter.FindOne(volumes, term)
	if err != nil {
		return kubernetes.Volume{}, fmt.Errorf("find volume: %w", err)
	}

	return volume, nil
}
