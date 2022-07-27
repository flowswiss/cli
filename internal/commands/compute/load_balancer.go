package compute

import (
	"context"
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/console"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func LoadBalancerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "load-balancer",
		Aliases: []string{"load-balancers", "loadbalancer", "loadbalancers"},
		Short:   "Manage compute load balancer",
	}

	commands.Add(cmd, &loadBalancerListCommand{}, &loadBalancerCreateCommand{}, &loadBalancerUpdateCommand{}, &loadBalancerDeleteCommand{})
	commands.Add(cmd, &loadBalancerProtocolListCommand{}, &loadBalancerAlgorithmListCommand{}, &loadBalancerHealthCheckTypeListCommand{})
	cmd.AddCommand(LoadBalancerPoolCommand(), LoadBalancerMemberCommand())

	return cmd
}

type loadBalancerListCommand struct {
	filter string
}

func (l *loadBalancerListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.NewLoadBalancerService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch loadBalancers: %w", err)
	}

	if len(l.filter) != 0 {
		items = filter.Find(items, l.filter)
	}

	return commands.PrintStdout(items)
}

func (l *loadBalancerListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List load balancers",
		Long:              "Lists all load balancers of the current tenant.",
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type loadBalancerCreateCommand struct {
	name      string
	internal  bool
	network   string
	privateIP net.IP
}

func (l *loadBalancerCreateCommand) Run(cmd *cobra.Command, args []string) error {
	network, err := findNetwork(cmd.Context(), l.network)
	if err != nil {
		return err
	}

	data := compute.LoadBalancerCreate{
		Name:             l.name,
		LocationID:       network.Location.ID,
		AttachExternalIP: !l.internal,
		NetworkID:        network.ID,
	}

	if len(l.privateIP) != 0 {
		data.PrivateIP = l.privateIP.String()
	}

	ordering, err := compute.NewLoadBalancerService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create load balancer: %w", err)
	}

	progress := console.NewProgress("Creating load balancer")
	go progress.Display(commands.Stderr)

	err = common.WaitForOrder(cmd.Context(), commands.Config.Client, ordering)
	if err != nil {
		progress.Complete("Order failed")

		return fmt.Errorf("wait for order: %w", err)
	}

	progress.Complete("Order completed")
	return nil
}

func (l *loadBalancerCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Short:             "Create a load balancer",
		Long:              "Creates a new load balancer",
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.name, "name", "", "name of the load balancer")
	cmd.Flags().BoolVar(&l.internal, "internal", false, "do not attach a public elastic ip to the load balancer")
	cmd.Flags().StringVar(&l.network, "network", "", "network to create the load balancer in")
	cmd.Flags().IPVar(&l.privateIP, "private-ip", net.IP{}, "private ip of the load balancer within the network")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

type loadBalancerUpdateCommand struct {
	name string
}

func (l *loadBalancerUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	data := compute.LoadBalancerUpdate{
		Name: l.name,
	}

	loadBalancer, err = compute.NewLoadBalancerService(commands.Config.Client).Update(cmd.Context(), loadBalancer.ID, data)
	if err != nil {
		return fmt.Errorf("update load balancer: %w", err)
	}

	return commands.PrintStdout(loadBalancer)
}

func (l *loadBalancerUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeLoadBalancer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerUpdateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update LOAD-BALANCER",
		Short:             "Update load balancer",
		Long:              "Updates a compute load balancer.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.name, "name", "", "name of the load balancer")

	return cmd
}

type loadBalancerDeleteCommand struct {
	force bool
}

func (l *loadBalancerDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !l.force && !commands.ConfirmDeletion("load balancer", loadBalancer) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NewLoadBalancerService(commands.Config.Client).Delete(cmd.Context(), loadBalancer.ID)
	if err != nil {
		return fmt.Errorf("delete load balancer: %w", err)
	}

	return nil
}

func (l *loadBalancerDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeLoadBalancer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete LOAD-BALANCER",
		Short:             "Delete load balancer",
		Long:              "Deletes a compute load balancer.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().BoolVar(&l.force, "force", false, "force the deletion of the load balancer without asking for confirmation")

	return cmd
}

type loadBalancerProtocolListCommand struct {
	filter string
}

func (l *loadBalancerProtocolListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.LoadBalancerProtocols(cmd.Context(), commands.Config.Client)
	if err != nil {
		return fmt.Errorf("fetch load balancer protocols: %w", err)
	}

	if len(l.filter) != 0 {
		items = filter.Find(items, l.filter)
	}

	return commands.PrintStdout(items)
}

func (l *loadBalancerProtocolListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerProtocolListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "protocol",
		Aliases:           []string{"protocols"},
		Short:             "List load balancer protocols",
		Long:              "Lists all load balancer protocols.",
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type loadBalancerAlgorithmListCommand struct {
	filter string
}

func (l *loadBalancerAlgorithmListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.LoadBalancerAlgorithms(cmd.Context(), commands.Config.Client)
	if err != nil {
		return fmt.Errorf("fetch load balancer algorithms: %w", err)
	}

	if len(l.filter) != 0 {
		items = filter.Find(items, l.filter)
	}

	return commands.PrintStdout(items)
}

func (l *loadBalancerAlgorithmListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerAlgorithmListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "algorithm",
		Aliases:           []string{"algorithms"},
		Short:             "List load balancer algorithms",
		Long:              "Lists all load balancer algorithms.",
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type loadBalancerHealthCheckTypeListCommand struct {
	filter string
}

func (l *loadBalancerHealthCheckTypeListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.LoadBalancerHealthCheckTypes(cmd.Context(), commands.Config.Client)
	if err != nil {
		return fmt.Errorf("fetch load balancer health check types: %w", err)
	}

	if len(l.filter) != 0 {
		items = filter.Find(items, l.filter)
	}

	return commands.PrintStdout(items)
}

func (l *loadBalancerHealthCheckTypeListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerHealthCheckTypeListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "health-check-type",
		Aliases:           []string{"health-check-types"},
		Short:             "List load balancer health check types",
		Long:              "Lists all load balancer health check types.",
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.filter, "filter", "", "custom term to filter the results")

	return cmd
}

func completeLoadBalancer(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	loadBalancers, err := compute.NewLoadBalancerService(commands.Config.Client).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(loadBalancers, term)

	names := make([]string, len(filtered))
	for i, loadBalancer := range filtered {
		names[i] = loadBalancer.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findLoadBalancer(ctx context.Context, term string) (compute.LoadBalancer, error) {
	loadBalancers, err := compute.NewLoadBalancerService(commands.Config.Client).List(ctx)
	if err != nil {
		return compute.LoadBalancer{}, fmt.Errorf("fetch load balancers: %w", err)
	}

	loadBalancer, err := filter.FindOne(loadBalancers, term)
	if err != nil {
		return compute.LoadBalancer{}, fmt.Errorf("find load balancer: %w", err)
	}

	return loadBalancer, nil
}
