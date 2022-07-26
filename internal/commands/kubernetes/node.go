package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/kubernetes"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func NodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "node",
		Aliases: []string{"nodes"},
		Short:   "Manage your cluster nodes",
	}

	commands.Add(cmd,
		&nodeListCommand{},
		&nodeDeleteCommand{},
	)

	cmd.AddCommand(
		NodeActionCommand(),
	)

	return cmd
}

type nodeListCommand struct {
	filter string
}

func (n *nodeListCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	items, err := kubernetes.NewNodeService(commands.Config.Client, cluster.ID).List(cmd.Context())
	if err != nil {
		return err
	}

	if len(n.filter) != 0 {
		items = filter.Find(items, n.filter)
	}

	return commands.PrintStdout(items)
}

func (n *nodeListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list CLUSTER",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List all nodes",
		Long:    "Prints a table of all nodes belonging to the selected cluster.",
		Args:    cobra.ExactArgs(1),
		RunE:    n.Run,
	}

	cmd.Flags().StringVar(&n.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type nodeDeleteCommand struct {
	force bool
}

func (n *nodeDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	node, err := findNode(cmd.Context(), cluster.ID, args[1])
	if err != nil {
		return err
	}

	if !n.force && !commands.ConfirmDeletion("node", node) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = kubernetes.NewNodeService(commands.Config.Client, cluster.ID).Delete(cmd.Context(), node.ID)
	if err != nil {
		return fmt.Errorf("delete node: %w", err)
	}

	return nil
}

func (n *nodeDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete CLUSTER NODE",
		Short: "Delete node",
		Long:  "Deletes a kubernetes node.",
		Args:  cobra.ExactArgs(2),
		RunE:  n.Run,
	}

	cmd.Flags().BoolVar(&n.force, "force", false, "forces deletion of the node without asking for confirmation")

	return cmd
}

func findNode(ctx context.Context, clusterID int, term string) (kubernetes.Node, error) {
	nodes, err := kubernetes.NewNodeService(commands.Config.Client, clusterID).List(ctx)
	if err != nil {
		return kubernetes.Node{}, fmt.Errorf("fetch nodes: %w", err)
	}

	node, err := filter.FindOne(nodes, term)
	if err != nil {
		return kubernetes.Node{}, fmt.Errorf("find node: %w", err)
	}

	return node, nil
}
