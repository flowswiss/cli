package compute

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/kubernetes"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ClusterActionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "action",
		Aliases: []string{"actions"},
		Short:   "Manage your kubernetes cluster actions",
	}

	commands.Add(cmd,
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

func (c *clusterActionListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list CLUSTER",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List available actions",
		Long:    "Prints a table of all available kubernetes cluster actions for the selected cluster.",
		Args:    cobra.ExactArgs(1),
		RunE:    c.Run,
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

func (c *clusterActionRunCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run CLUSTER ACTION",
		Short: "Run action",
		Long:  "Runs the given action on the selected kubernetes cluster.",
		Args:  cobra.ExactArgs(2),
		RunE:  c.Run,
	}

	return cmd
}
