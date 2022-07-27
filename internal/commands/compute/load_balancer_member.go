package compute

import (
	"context"
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func LoadBalancerMemberCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "member",
		Aliases: []string{"members"},
		Short:   "Manage load balancer members",
	}

	commands.Add(cmd,
		&loadBalancerMemberListCommand{},
		&loadBalancerMemberCreateCommand{},
		&loadBalancerMemberDeleteCommand{},
	)

	return cmd
}

type loadBalancerMemberListCommand struct {
	filter string
}

func (l *loadBalancerMemberListCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	pool, err := findLoadBalancerPool(cmd.Context(), loadBalancer.ID, args[1])
	if err != nil {
		return err
	}

	service := compute.NewLoadBalancerMemberService(commands.Config.Client, loadBalancer.ID, pool.ID)

	items, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch loadBalancerMembers: %w", err)
	}

	if len(l.filter) != 0 {
		items = filter.Find(items, l.filter)
	}

	return commands.PrintStdout(items)
}

func (l *loadBalancerMemberListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeLoadBalancer(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeLoadBalancerPool(cmd.Context(), loadBalancer, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerMemberListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list LOAD-BALANCER POOL",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List load balancer members",
		Long:              "Lists all load balancer member of the selected load balancer pool.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type loadBalancerMemberCreateCommand struct {
	name    string
	address net.IP
	port    int
}

func (l *loadBalancerMemberCreateCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	pool, err := findLoadBalancerPool(cmd.Context(), loadBalancer.ID, args[1])
	if err != nil {
		return err
	}

	service := compute.NewLoadBalancerMemberService(commands.Config.Client, loadBalancer.ID, pool.ID)

	data := compute.LoadBalancerMemberCreate{
		Name:    l.name,
		Address: l.address.String(),
		Port:    l.port,
	}

	item, err := service.Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create load balancer member: %w", err)
	}

	return commands.PrintStdout(item)
}

func (l *loadBalancerMemberCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeLoadBalancer(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeLoadBalancerPool(cmd.Context(), loadBalancer, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerMemberCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create LOAD-BALANCER POOL",
		Short:             "Create a load balancer member",
		Long:              "Creates a new load balancer member.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.name, "name", "", "name of the load balancer member")
	cmd.Flags().IPVar(&l.address, "address", net.IP{}, "ip address of the load balancer member")
	cmd.Flags().IntVar(&l.port, "port", 0, "port of the load balancer member")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("address")
	_ = cmd.MarkFlagRequired("port")

	return cmd
}

type loadBalancerMemberDeleteCommand struct {
	force bool
}

func (l *loadBalancerMemberDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	pool, err := findLoadBalancerPool(cmd.Context(), loadBalancer.ID, args[1])
	if err != nil {
		return err
	}

	service := compute.NewLoadBalancerMemberService(commands.Config.Client, loadBalancer.ID, pool.ID)

	members, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch load balancer members: %w", err)
	}

	member, err := filter.FindOne(members, args[2])
	if err != nil {
		return fmt.Errorf("find load balancer member: %w", err)
	}

	if !l.force && !commands.ConfirmDeletion("load balancer member", member) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = service.Delete(cmd.Context(), member.ID)
	if err != nil {
		return fmt.Errorf("delete load balancer member: %w", err)
	}

	return nil
}

func (l *loadBalancerMemberDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeLoadBalancer(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeLoadBalancerPool(cmd.Context(), loadBalancer, toComplete)
	}

	if len(args) == 2 {
		loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		pool, err := findLoadBalancerPool(cmd.Context(), loadBalancer.ID, args[1])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeLoadBalancerMember(cmd.Context(), loadBalancer, pool, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerMemberDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete LOAD-BALANCER POOL MEMBER",
		Short:             "Delete load balancer member",
		Long:              "Deletes a load balancer member.",
		Args:              cobra.ExactArgs(3),
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().BoolVar(&l.force, "force", false, "force the deletion of the loadBalancerMember without asking for confirmation")

	return cmd
}

func completeLoadBalancerMember(ctx context.Context, loadBalancer compute.LoadBalancer, pool compute.LoadBalancerPool, term string) ([]string, cobra.ShellCompDirective) {
	members, err := compute.NewLoadBalancerMemberService(commands.Config.Client, loadBalancer.ID, pool.ID).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(members, term)

	names := make([]string, len(filtered))
	for i, member := range filtered {
		names[i] = member.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
