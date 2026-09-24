package kubernetes

import (
	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/kubernetes"
	"github.com/flowswiss/cli/v2/pkg/filter"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/spf13/cobra"
)

func SnapshotCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "snapshot",
		Aliases: []string{"snapshots"},
		Short:   "Manage your cluster snapshots",
	}

	commands.Add(app, cmd,
		&snapshotListCommand{},
	)

	return cmd
}

type snapshotListCommand struct {
	filter string
}

func (s *snapshotListCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	items, err := kubernetes.SnapshotService().List(cmd.Context(), kubernetes.SnapshotList{ClusterID: uint(cluster.ID), Cursor: core.CursorAll})
	if err != nil {
		return err
	}

	if len(s.filter) != 0 {
		items = filter.Find(items, s.filter)
	}

	return commands.PrintStdout(items)
}

func (s *snapshotListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeCluster(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *snapshotListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list CLUSTER",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List all snapshots",
		Long:              "Prints a table of all snapshots belonging to the selected cluster.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.filter, "filter", "", "custom term to filter the results")

	return cmd
}
