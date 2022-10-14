package kubernetes

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/kubernetes"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ClusterActionCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "action",
		Aliases: []string{"actions"},
		Short:   "Manage your kubernetes cluster actions",
	}

	commands.Add(app, cmd,
		&clusterActionListCommand{},
		&clusterActionRunCommand{},
	)

	return cmd
}

type clusterActionListCommand struct {
	filter string
}

func (c *clusterActionListCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	actions := make([]kubernetes.ClusterAction, len(cluster.Status.Actions))
	for i, action := range cluster.Status.Actions {
		actions[i] = kubernetes.ClusterAction(action)
	}

	if len(c.filter) != 0 {
		actions = filter.Find(actions, c.filter)
	}

	return commands.PrintStdout(actions)
}

func (c *clusterActionListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeCluster(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *clusterActionListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list CLUSTER",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List available actions",
		Long:              "Prints a table of all available kubernetes cluster actions for the selected cluster.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().StringVar(&c.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type clusterActionRunCommand struct{}

func (c *clusterActionRunCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	actions := make([]kubernetes.ClusterAction, len(cluster.Status.Actions))
	for i, action := range cluster.Status.Actions {
		actions[i] = kubernetes.ClusterAction(action)
	}

	action, err := filter.FindOne(actions, args[1])
	if err != nil {
		return fmt.Errorf("find action: %w", err)
	}

	data := kubernetes.ClusterPerformAction{
		Action: action.Command,
	}

	cluster, err = kubernetes.NewClusterService(commands.Config.Client).PerformAction(cmd.Context(), cluster.ID, data)
	if err != nil {
		return fmt.Errorf("run cluster action: %w", err)
	}

	return commands.PrintStdout(cluster)
}

func (c *clusterActionRunCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeCluster(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		cluster, err := findCluster(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		actions := make([]kubernetes.ClusterAction, len(cluster.Status.Actions))
		for i, action := range cluster.Status.Actions {
			actions[i] = kubernetes.ClusterAction(action)
		}

		filtered := filter.Find(actions, toComplete)

		names := make([]string, len(filtered))
		for i, action := range filtered {
			names[i] = action.Command
		}

		return names, cobra.ShellCompDirectiveNoFileComp
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *clusterActionRunCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "run CLUSTER ACTION",
		Short:             "Run action",
		Long:              "Runs the given action on the selected kubernetes cluster.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	return cmd
}
