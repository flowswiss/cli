package compute

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/kubernetes"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func NodeActionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "action",
		Aliases: []string{"actions"},
		Short:   "Manage your kubernetes node actions",
	}

	commands.Add(cmd,
		&nodeActionListCommand{},
		&nodeActionRunCommand{},
	)

	return cmd
}

type nodeActionListCommand struct {
	filter string
}

func (n *nodeActionListCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	node, err := findNode(cmd.Context(), cluster.ID, args[1])
	if err != nil {
		return err
	}

	actions := make([]kubernetes.NodeAction, len(node.Status.Actions))
	for i, action := range node.Status.Actions {
		actions[i] = kubernetes.NodeAction(action)
	}

	if len(n.filter) != 0 {
		actions = filter.Find(actions, n.filter)
	}

	return commands.PrintStdout(actions)
}

func (n *nodeActionListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeCluster(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		cluster, err := findCluster(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeNode(cmd.Context(), cluster, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *nodeActionListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list CLUSTER NODE",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List available actions",
		Long:              "Prints a table of all available kubernetes node actions for the selected node.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().StringVar(&n.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type nodeActionRunCommand struct{}

func (n *nodeActionRunCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	node, err := findNode(cmd.Context(), cluster.ID, args[1])
	if err != nil {
		return err
	}

	actions := make([]kubernetes.NodeAction, len(node.Status.Actions))
	for i, action := range node.Status.Actions {
		actions[i] = kubernetes.NodeAction(action)
	}

	action, err := filter.FindOne(actions, args[2])
	if err != nil {
		return fmt.Errorf("find action: %w", err)
	}

	data := kubernetes.NodePerformAction{
		Action: action.Command,
	}

	node, err = kubernetes.NewNodeService(commands.Config.Client, cluster.ID).PerformAction(cmd.Context(), node.ID, data)
	if err != nil {
		return fmt.Errorf("run node action: %w", err)
	}

	return commands.PrintStdout(node)
}

func (n *nodeActionRunCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeCluster(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		cluster, err := findCluster(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeNode(cmd.Context(), cluster, toComplete)
	}

	if len(args) == 1 {
		cluster, err := findCluster(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		node, err := findNode(cmd.Context(), cluster.ID, args[1])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		actions := make([]kubernetes.NodeAction, len(node.Status.Actions))
		for i, action := range node.Status.Actions {
			actions[i] = kubernetes.NodeAction(action)
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

func (n *nodeActionRunCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "run CLUSTER NODE ACTION",
		Short:             "Run action",
		Long:              "Runs the given action on the selected kubernetes node.",
		Args:              cobra.ExactArgs(3),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	return cmd
}
