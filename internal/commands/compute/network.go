package compute

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/flowswiss/cli/v2/pkg/optional"
	"github.com/spf13/cobra"

	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func NetworkCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage compute networks",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # List all networks
      %[1]s compute network list
      
      # Create a new network
      %[1]s compute network create --name my-network --location ALP1
		`, app.Name)),
	}

	commands.Add(app, cmd,
		&networkListCommand{},
		&networkCreateCommand{},
		&networkUpdateCommand{},
		&networkDeleteCommand{},
	)

	return cmd
}

type networkListCommand struct {
	filter string
}

func (n *networkListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.NetworkService().List(cmd.Context(), core.CursorAll)
	if err != nil {
		return fmt.Errorf("fetch networks: %w", err)
	}

	if len(n.filter) != 0 {
		items = filter.Find(items, n.filter)
	}

	return commands.PrintStdout(items)
}

func (n *networkListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List networks",
		Long:              "Lists all networks of the current tenant.",
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().StringVar(&n.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type networkCreateCommand struct {
	name                string
	description         optional.Optional[string]
	location            string
	domainNameServers   []net.IP
	cidr                net.IPNet
	allocationPoolStart optional.Optional[net.IP]
	allocationPoolEnd   optional.Optional[net.IP]
	gateway             optional.Optional[net.IP]
}

func (n *networkCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Client, n.location)
	if err != nil {
		return err
	}

	domainNameServers := make([]string, len(n.domainNameServers))
	for i, dns := range n.domainNameServers {
		domainNameServers[i] = dns.String()
	}

	var allocationPoolStart *string
	if poolStart := n.allocationPoolStart.Value(); poolStart != nil {
		if !n.cidr.Contains(*poolStart) {
			return fmt.Errorf("start address of the allocation pool is not within the network cidr")
		}

		allocationPoolStart = new(poolStart.String())
	}

	var allocationPoolEnd *string
	if poolEnd := n.allocationPoolEnd.Value(); poolEnd != nil {
		if !n.cidr.Contains(*poolEnd) {
			return fmt.Errorf("end address of the allocation pool is not within the network cidr")
		}

		if bytes.Compare(*n.allocationPoolStart.Value(), *poolEnd) > 0 {
			return fmt.Errorf("start address of the allocation pool is greater than the end address of the allocation pool")
		}

		allocationPoolEnd = new(poolEnd.String())
	}

	var gateway *string
	if g := n.gateway.Value(); g != nil {
		if !n.cidr.Contains(*g) {
			return fmt.Errorf("gateway address is not within the network cidr")
		}

		gateway = new(g.String())
	}

	data := compute.NetworkCreate{
		Name:                &n.name,
		Description:         n.description.Value(),
		LocationID:          &location.ID,
		DomainNameServers:   domainNameServers,
		CIDR:                new(n.cidr.String()),
		AllocationPoolStart: allocationPoolStart,
		AllocationPoolEnd:   allocationPoolEnd,
		GatewayIP:           gateway,
	}

	item, err := compute.NetworkService().Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create network: %w", err)
	}

	return commands.PrintStdout(item)
}

func (n *networkCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkCreateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a network",
		Long:  "Creates a new network",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # Create a new network using default cidr
      %[1]s compute network create --name my-network --location ALP1
      
      # Create a new network using custom cidr
      %[1]s compute network create --name my-network --location ALP1 --cidr 10.0.0.0/24
      
      # Create a new network using custom allocation pool
      %[1]s compute network create --name my-network --location ALP1 --cidr 10.0.0.0/16 --allocation-pool-start 10.0.1.0 --allocation-pool-end 10.0.1.255
		`, app.Name)),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	defaultNet := net.IPNet{
		IP:   net.IPv4(172, 16, 0, 0),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	}

	cmd.Flags().StringVar(&n.name, "name", "", "name of the new network")
	cmd.Flags().StringVar(n.description.Configure(cmd, "description", "description of the network"))
	cmd.Flags().StringVar(&n.location, "location", "", "location where the network will be created")
	cmd.Flags().IPSliceVar(&n.domainNameServers, "domain-name-server", []net.IP{net.IPv4(1, 1, 1, 1), net.IPv4(8, 8, 8, 8)}, "domain name servers of the network")
	cmd.Flags().IPNetVar(&n.cidr, "cidr", defaultNet, "subnet cidr for the network")
	cmd.Flags().IPVar(n.allocationPoolStart.Configure(cmd, "allocation-pool-start", "start address of the dhcp allocation pool"))
	cmd.Flags().IPVar(n.allocationPoolEnd.Configure(cmd, "allocation-pool-end", "end address of the dhcp allocation pool"))
	cmd.Flags().IPVar(n.gateway.Configure(cmd, "gateway", "gateway address of the network"))

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("location")
	cmd.MarkFlagsRequiredTogether("allocation-pool-start", "allocation-pool-end")

	return cmd
}

type networkUpdateCommand struct {
	name                optional.Optional[string]
	description         optional.Optional[string]
	domainNameServers   []net.IP
	allocationPoolStart optional.Optional[net.IP]
	allocationPoolEnd   optional.Optional[net.IP]
	gateway             optional.Optional[net.IP]
}

func (n *networkUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	network, err := findNetwork(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	_, cidr, err := net.ParseCIDR(network.CIDR)
	if err != nil {
		return fmt.Errorf("parse network cidr: %w", err)
	}

	domainNameServers := make([]string, len(n.domainNameServers))
	for i, dns := range n.domainNameServers {
		domainNameServers[i] = dns.String()
	}

	var allocationPoolStart *string
	if poolStart := n.allocationPoolStart.Value(); poolStart != nil {
		if !cidr.Contains(*poolStart) {
			return fmt.Errorf("start address of the allocation pool is not within the network cidr")
		}

		allocationPoolStart = new(poolStart.String())
	}

	var allocationPoolEnd *string
	if poolEnd := n.allocationPoolEnd.Value(); poolEnd != nil {
		if !cidr.Contains(*poolEnd) {
			return fmt.Errorf("end address of the allocation pool is not within the network cidr")
		}

		if poolStart := n.allocationPoolStart.Value(); poolStart != nil && bytes.Compare(*poolStart, *poolEnd) > 0 {
			return fmt.Errorf("start address of the allocation pool is greater than the end address of the allocation pool")
		}

		allocationPoolEnd = new(poolEnd.String())
	}

	var gateway *string
	if g := n.gateway.Value(); g != nil {
		if !cidr.Contains(*g) {
			return fmt.Errorf("gateway address is not within the network cidr")
		}

		gateway = new(g.String())
	}

	data := compute.NetworkUpdate{ID: uint(network.ID),
		Name:                n.name.Value(),
		Description:         n.description.Value(),
		DomainNameServers:   domainNameServers,
		AllocationPoolStart: allocationPoolStart,
		AllocationPoolEnd:   allocationPoolEnd,
		GatewayIP:           gateway,
	}

	network, err = compute.NetworkService().Update(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("update network: %w", err)
	}

	return commands.PrintStdout(network)
}

func (n *networkUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeNetwork(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkUpdateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update NETWORK",
		Short:             "Update network",
		Long:              "Updates a compute network.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().StringVar(n.name.Configure(cmd, "name", "name of the network"))
	cmd.Flags().StringVar(n.description.Configure(cmd, "description", "description of the network"))
	cmd.Flags().IPSliceVar(&n.domainNameServers, "domain-name-server", nil, "domain name servers of the network")
	cmd.Flags().IPVar(n.allocationPoolStart.Configure(cmd, "allocation-pool-start", "start address of the dhcp allocation pool"))
	cmd.Flags().IPVar(n.allocationPoolEnd.Configure(cmd, "allocation-pool-end", "end address of the dhcp allocation pool"))
	cmd.Flags().IPVar(n.gateway.Configure(cmd, "gateway", "gateway address of the network"))

	return cmd
}

type networkDeleteCommand struct {
	force bool
}

func (n *networkDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	network, err := findNetwork(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !n.force && !commands.ConfirmDeletion("network", network) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NetworkService().Delete(cmd.Context(), compute.NetworkDelete{ID: uint(network.ID)})
	if err != nil {
		return fmt.Errorf("delete network: %w", err)
	}

	return nil
}

func (n *networkDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeNetwork(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkDeleteCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete NETWORK",
		Short:             "Delete network",
		Long:              "Deletes a compute network.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().BoolVar(&n.force, "force", false, "force the deletion of the network without asking for confirmation")

	return cmd
}

func completeNetwork(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	networks, err := compute.NetworkService().List(ctx, core.CursorAll)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(networks, term)

	names := make([]string, len(filtered))
	for i, network := range filtered {
		names[i] = network.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findNetwork(ctx context.Context, term string) (compute.Network, error) {
	networks, err := compute.NetworkService().List(ctx, core.CursorAll)
	if err != nil {
		return compute.Network{}, fmt.Errorf("fetch networks: %w", err)
	}

	network, err := filter.FindOne(networks, term)
	if err != nil {
		return compute.Network{}, fmt.Errorf("find network: %w", err)
	}

	return network, nil
}
