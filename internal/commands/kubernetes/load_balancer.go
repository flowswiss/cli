package compute

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/kubernetes"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func LoadBalancerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "load-balancer",
		Aliases: []string{"load-balancers"},
		Short:   "Manage your cluster load-balancer",
	}

	commands.Add(cmd,
		&loadBalancerListCommand{},
	)

	return cmd
}

type loadBalancerListCommand struct {
	filter string
}

func (l *loadBalancerListCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	items, err := kubernetes.NewLoadBalancerService(commands.Config.Client, cluster.ID).List(cmd.Context())
	if err != nil {
		return err
	}

	if len(l.filter) != 0 {
		items = filter.Find(items, l.filter)
	}

	return commands.PrintStdout(items)
}

func (l *loadBalancerListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeCluster(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list CLUSTER",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List all load balancer",
		Long:              "Prints a table of all load balancer belonging to the selected cluster.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.filter, "filter", "", "custom term to filter the results")

	return cmd
}
