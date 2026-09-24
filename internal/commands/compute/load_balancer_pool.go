package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/flowswiss/cli/v2/pkg/optional"
	"github.com/spf13/cobra"

	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func LoadBalancerPoolCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Aliases: []string{"pools"},
		Short:   "Manage load balancer pools",
	}

	commands.Add(app, cmd,
		&loadBalancerPoolListCommand{},
		&loadBalancerPoolCreateCommand{},
		&loadBalancerPoolUpdateCommand{},
		&loadBalancerPoolDeleteCommand{},
	)

	return cmd
}

type loadBalancerPoolListCommand struct {
	filter string
}

func (l *loadBalancerPoolListCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	items, err := compute.LoadBalancerPoolService().List(cmd.Context(), compute.LoadBalancerPoolList{LoadBalancerID: uint(loadBalancer.ID), Cursor: core.CursorAll})
	if err != nil {
		return fmt.Errorf("fetch load balancer pools: %w", err)
	}

	if len(l.filter) != 0 {
		items = filter.Find(items, l.filter)
	}

	return commands.PrintStdout(items)
}

func (l *loadBalancerPoolListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeLoadBalancer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerPoolListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list LOAD-BALANCER",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List load balancer pools",
		Long:              "Lists all pools of the selected load balancer.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type loadBalancerPoolCreateCommand struct {
	entryProtocol  string
	targetProtocol string
	certificate    optional.Optional[string]
	entryPort      int
	algorithm      string
	stickySession  bool

	healthCheckType               string
	healthCheckHTTPMethod         optional.Optional[string]
	healthCheckHTTPPath           optional.Optional[string]
	healthCheckInterval           optional.Optional[time.Duration]
	healthCheckTimeout            optional.Optional[time.Duration]
	healthCheckHealthyThreshold   optional.Optional[int]
	healthCheckUnhealthyThreshold optional.Optional[int]
}

func (l *loadBalancerPoolCreateCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	protocols, err := compute.LoadBalancerProtocols(cmd.Context(), commands.Client)
	if err != nil {
		return fmt.Errorf("fetch load balancer protocols: %w", err)
	}

	algorithms, err := compute.LoadBalancerAlgorithms(cmd.Context(), commands.Client)
	if err != nil {
		return fmt.Errorf("fetch load balancer algorithms: %w", err)
	}

	healthCheckTypes, err := compute.LoadBalancerHealthCheckTypes(cmd.Context(), commands.Client)
	if err != nil {
		return fmt.Errorf("fetch load balancer health check types: %w", err)
	}

	entryProtocol, err := filter.FindOne(protocols, l.entryProtocol)
	if err != nil {
		return fmt.Errorf("find entry protocol: %w", err)
	}

	targetProtocol, err := filter.FindOne(protocols, l.targetProtocol)
	if err != nil {
		return fmt.Errorf("find target protocol: %w", err)
	}

	algorithm, err := filter.FindOne(algorithms, l.algorithm)
	if err != nil {
		return fmt.Errorf("find balancing algorithm: %w", err)
	}

	healthCheckType, err := filter.FindOne(healthCheckTypes, l.healthCheckType)
	if err != nil {
		return fmt.Errorf("find health check type: %w", err)
	}

	var hcInterval *int
	if interval := l.healthCheckInterval.Value(); interval != nil {
		hcInterval = new(int(interval.Seconds()))
	}

	var hcTimeout *int
	if timeout := l.healthCheckTimeout.Value(); timeout != nil {
		hcTimeout = new(int(timeout.Seconds()))
	}

	data := compute.LoadBalancerPoolCreate{LoadBalancerID: uint(loadBalancer.ID),
		EntryProtocolID:      entryProtocol.ID,
		TargetProtocolID:     targetProtocol.ID,
		EntryPort:            l.entryPort,
		BalancingAlgorithmID: algorithm.ID,
		StickySession:        l.stickySession,

		HealthCheck: &compute.LoadBalancerHealthCheckOptions{
			TypeID:             healthCheckType.ID,
			HTTPMethod:         l.healthCheckHTTPMethod.Value(),
			HTTPPath:           l.healthCheckHTTPPath.Value(),
			Interval:           hcInterval,
			Timeout:            hcTimeout,
			HealthyThreshold:   l.healthCheckHealthyThreshold.Value(),
			UnhealthyThreshold: l.healthCheckUnhealthyThreshold.Value(),
		},
	}

	if certIdentifier := l.certificate.Value(); certIdentifier != nil {
		certificate, err := findCertificate(cmd.Context(), *certIdentifier)
		if err != nil {
			return err
		}

		data.CertificateID = &certificate.ID
	}

	item, err := compute.LoadBalancerPoolService().Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create load balancer pool: %w", err)
	}

	return commands.PrintStdout(item)
}

func (l *loadBalancerPoolCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeLoadBalancer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (l *loadBalancerPoolCreateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create LOAD-BALANCER",
		Short:             "Create a load balancer pool",
		Long:              "Creates a new load balancer pool",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().StringVar(&l.entryProtocol, "entry-protocol", "", "name of the entry protocol to use")
	cmd.Flags().StringVar(&l.targetProtocol, "target-protocol", "", "name of the target protocol to use")
	cmd.Flags().StringVar(l.certificate.Configure(cmd, "certificate", "name of the certificate to use"))
	cmd.Flags().IntVar(&l.entryPort, "entry-port", 0, "port of the entry protocol")
	cmd.Flags().StringVar(&l.algorithm, "algorithm", "", "name of the balancing algorithm to use")
	cmd.Flags().BoolVar(&l.stickySession, "sticky-session", false, "enable sticky session")

	cmd.Flags().StringVar(&l.healthCheckType, "health-check-type", "", "type of the health check")
	cmd.Flags().StringVar(l.healthCheckHTTPMethod.Configure(cmd, "health-check-http-method", "HTTP method to use for the health check"))
	cmd.Flags().StringVar(l.healthCheckHTTPPath.Configure(cmd, "health-check-http-path", "HTTP path to use for the health check"))
	cmd.Flags().DurationVar(l.healthCheckInterval.Configure(cmd, "health-check-interval", "interval of the health check"))
	cmd.Flags().DurationVar(l.healthCheckTimeout.Configure(cmd, "health-check-timeout", "timeout of the health check"))
	cmd.Flags().IntVar(l.healthCheckHealthyThreshold.Configure(cmd, "health-check-healthy-threshold", "healthy threshold of the health check"))
	cmd.Flags().IntVar(l.healthCheckUnhealthyThreshold.Configure(cmd, "health-check-unhealthy-threshold", "unhealthy threshold of the health check"))

	_ = cmd.MarkFlagRequired("entry-protocol")
	_ = cmd.MarkFlagRequired("target-protocol")
	_ = cmd.MarkFlagRequired("entry-port")
	_ = cmd.MarkFlagRequired("algorithm")

	return cmd
}

type loadBalancerPoolUpdateCommand struct {
	certificate   optional.Optional[string]
	algorithm     optional.Optional[string]
	stickySession optional.Optional[bool]

	healthCheckType               optional.Optional[string]
	healthCheckHTTPMethod         optional.Optional[string]
	healthCheckHTTPPath           optional.Optional[string]
	healthCheckInterval           optional.Optional[time.Duration]
	healthCheckTimeout            optional.Optional[time.Duration]
	healthCheckHealthyThreshold   optional.Optional[int]
	healthCheckUnhealthyThreshold optional.Optional[int]
}

func (l *loadBalancerPoolUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	loadBalancerPool, err := findLoadBalancerPool(cmd.Context(), loadBalancer.ID, args[1])
	if err != nil {
		return err
	}

	var hcInterval *int
	if interval := l.healthCheckInterval.Value(); interval != nil {
		hcInterval = new(int(interval.Seconds()))
	}

	var hcTimeout *int
	if timeout := l.healthCheckTimeout.Value(); timeout != nil {
		hcTimeout = new(int(timeout.Seconds()))
	}

	data := compute.LoadBalancerPoolUpdate{LoadBalancerID: uint(loadBalancer.ID), LoadBalancerPoolID: uint(loadBalancerPool.ID),
		StickySession: l.stickySession.Value(),
		HealthCheck: &compute.LoadBalancerHealthCheckOptions{
			HTTPMethod:         l.healthCheckHTTPMethod.Value(),
			HTTPPath:           l.healthCheckHTTPPath.Value(),
			Interval:           hcInterval,
			Timeout:            hcTimeout,
			HealthyThreshold:   l.healthCheckHealthyThreshold.Value(),
			UnhealthyThreshold: l.healthCheckUnhealthyThreshold.Value(),
		},
	}

	if certIdentifier := l.certificate.Value(); certIdentifier != nil {
		certificate, err := findCertificate(cmd.Context(), *certIdentifier)
		if err != nil {
			return err
		}

		data.CertificateID = &certificate.ID
	}

	if algIdentifier := l.algorithm.Value(); algIdentifier != nil {
		algorithms, err := compute.LoadBalancerAlgorithms(cmd.Context(), commands.Client)
		if err != nil {
			return fmt.Errorf("fetch load balancer algorithms: %w", err)
		}

		algorithm, err := filter.FindOne(algorithms, *algIdentifier)
		if err != nil {
			return fmt.Errorf("find balancing algorithm: %w", err)
		}

		data.BalancingAlgorithmID = &algorithm.ID
	}

	if hcTypeIdentifier := l.healthCheckType.Value(); hcTypeIdentifier != nil {
		healthCheckTypes, err := compute.LoadBalancerHealthCheckTypes(cmd.Context(), commands.Client)
		if err != nil {
			return fmt.Errorf("fetch load balancer health check types: %w", err)
		}

		healthCheckType, err := filter.FindOne(healthCheckTypes, *hcTypeIdentifier)
		if err != nil {
			return fmt.Errorf("find health check type: %w", err)
		}

		data.HealthCheck.TypeID = healthCheckType.ID
	}

	loadBalancerPool, err = compute.LoadBalancerPoolService().Update(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("update load balancer pool: %w", err)
	}

	return commands.PrintStdout(loadBalancerPool)
}

func (l *loadBalancerPoolUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
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

func (l *loadBalancerPoolUpdateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update LOAD-BALANCER POOL",
		Short: "Update load balancer pool",
		Long:  "Updates a compute load balancer pool.",
		Args:  cobra.ExactArgs(2),
		RunE:  l.Run,
	}

	cmd.Flags().StringVar(l.certificate.Configure(cmd, "certificate", "name of the certificate to use"))
	cmd.Flags().StringVar(l.algorithm.Configure(cmd, "algorithm", "name of the balancing algorithm to use"))
	cmd.Flags().BoolVar(l.stickySession.Configure(cmd, "sticky-session", "enable sticky session"))

	cmd.Flags().StringVar(l.healthCheckType.Configure(cmd, "health-check-type", "type of the health check"))
	cmd.Flags().StringVar(l.healthCheckHTTPMethod.Configure(cmd, "health-check-http-method", "HTTP method to use for the health check"))
	cmd.Flags().StringVar(l.healthCheckHTTPPath.Configure(cmd, "health-check-http-path", "HTTP path to use for the health check"))
	cmd.Flags().DurationVar(l.healthCheckInterval.Configure(cmd, "health-check-interval", "interval of the health check"))
	cmd.Flags().DurationVar(l.healthCheckTimeout.Configure(cmd, "health-check-timeout", "timeout of the health check"))
	cmd.Flags().IntVar(l.healthCheckHealthyThreshold.Configure(cmd, "health-check-healthy-threshold", "healthy threshold of the health check"))
	cmd.Flags().IntVar(l.healthCheckUnhealthyThreshold.Configure(cmd, "health-check-unhealthy-threshold", "unhealthy threshold of the health check"))

	return cmd
}

type loadBalancerPoolDeleteCommand struct {
	force bool
}

func (l *loadBalancerPoolDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	loadBalancerPool, err := findLoadBalancerPool(cmd.Context(), loadBalancer.ID, args[1])
	if err != nil {
		return err
	}

	if !l.force && !commands.ConfirmDeletion("load balancer pool", loadBalancerPool) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.LoadBalancerPoolService().Delete(cmd.Context(), compute.LoadBalancerPoolDelete{LoadBalancerID: uint(loadBalancer.ID), LoadBalancerPoolID: uint(loadBalancerPool.ID)})
	if err != nil {
		return fmt.Errorf("delete load balancer pool: %w", err)
	}

	return nil
}

func (l *loadBalancerPoolDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
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

func (l *loadBalancerPoolDeleteCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete LOAD-BALANCER POOL",
		Short:             "Delete load balancer pool",
		Long:              "Deletes a compute load balancer pool.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: l.CompleteArg,
		RunE:              l.Run,
	}

	cmd.Flags().BoolVar(&l.force, "force", false, "force the deletion of the load balancer pool without asking for confirmation")

	return cmd
}

func completeLoadBalancerPool(ctx context.Context, loadBalancer compute.LoadBalancer, term string) (
	[]string,
	cobra.ShellCompDirective,
) {
	loadBalancerPools, err := compute.LoadBalancerPoolService().List(ctx, compute.LoadBalancerPoolList{LoadBalancerID: uint(loadBalancer.ID), Cursor: core.CursorAll})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(loadBalancerPools, term)

	names := make([]string, len(filtered))
	for i, loadBalancerPool := range filtered {
		names[i] = loadBalancerPool.NameWithoutSpaces()
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findLoadBalancerPool(ctx context.Context, loadBalancerID int, term string) (compute.LoadBalancerPool, error) {
	loadBalancerPools, err := compute.LoadBalancerPoolService().List(ctx, compute.LoadBalancerPoolList{LoadBalancerID: uint(loadBalancerID), Cursor: core.CursorAll})
	if err != nil {
		return compute.LoadBalancerPool{}, fmt.Errorf("fetch load balancer pools: %w", err)
	}

	loadBalancerPool, err := filter.FindOne(loadBalancerPools, term)
	if err != nil {
		return compute.LoadBalancerPool{}, fmt.Errorf("find load balancer pool: %w", err)
	}

	return loadBalancerPool, nil
}
