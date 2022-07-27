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

func NetworkInterfaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network-interface",
		Aliases: []string{"network-interfaces"},
		Short:   "Manage network interfaces",
	}

	commands.Add(cmd, &networkInterfaceListCommand{}, &networkInterfaceCreateCommand{}, &networkInterfaceUpdateCommand{}, &networkInterfaceDeleteCommand{})

	return cmd
}

type networkInterfaceListCommand struct {
	filter string
}

func (n *networkInterfaceListCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	service := compute.NewNetworkInterfaceService(commands.Config.Client, server.ID)

	items, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch network interfaces: %w", err)
	}

	if len(n.filter) != 0 {
		items = filter.Find(items, n.filter)
	}

	return commands.PrintStdout(items)
}

func (n *networkInterfaceListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n networkInterfaceListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list SERVER",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List network interfaces",
		Long:              "Lists all network interfaces of the current server.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().StringVar(&n.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type networkInterfaceCreateCommand struct {
	network   string
	privateIP net.IP
}

func (n *networkInterfaceCreateCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	network, err := findNetwork(cmd.Context(), n.network)
	if err != nil {
		return err
	}

	privateIP := ""
	if len(n.privateIP) != 0 {
		_, cidr, err := net.ParseCIDR(network.CIDR)
		if err != nil {
			return fmt.Errorf("parse CIDR: %w", err)
		}

		if !cidr.Contains(n.privateIP) {
			return fmt.Errorf("private IP %s is not in the network %s", n.privateIP, network.CIDR)
		}

		privateIP = n.privateIP.String()
	}

	data := compute.NetworkInterfaceCreate{
		NetworkID: network.ID,
		PrivateIP: privateIP,
	}

	iface, err := compute.NewNetworkInterfaceService(commands.Config.Client, server.ID).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create network interface: %w", err)
	}

	return commands.PrintStdout(iface)
}

func (n *networkInterfaceCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkInterfaceCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create SERVER",
		Aliases:           []string{"add", "new"},
		Short:             "Create a network interface",
		Long:              "Creates a new network interface for the current server.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().StringVar(&n.network, "network", "", "the network to use")
	cmd.Flags().IPVar(&n.privateIP, "private-ip", nil, "the private IP to use")

	_ = cmd.MarkFlagRequired("network")

	return cmd
}

type networkInterfaceUpdateCommand struct {
	disableSecurity bool
	enableSecurity  bool
	securityGroups  []string
}

func (n *networkInterfaceUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	iface, err := findNetworkInterface(cmd.Context(), server.ID, args[1])
	if err != nil {
		return err
	}

	service := compute.NewNetworkInterfaceService(commands.Config.Client, server.ID)

	if n.disableSecurity || n.enableSecurity {
		data := compute.NetworkInterfaceSecurityUpdate{
			Security: n.enableSecurity,
		}

		iface, err = service.UpdateSecurity(cmd.Context(), iface.ID, data)
		if err != nil {
			return fmt.Errorf("update network interface security: %w", err)
		}
	}

	if len(n.securityGroups) != 0 {
		if !iface.Security {
			return fmt.Errorf("cannot update security groups of a non-security network interface")
		}

		securityGroupIDs := make([]int, len(n.securityGroups))
		for idx, group := range n.securityGroups {
			securityGroup, err := findSecurityGroup(cmd.Context(), group)
			if err != nil {
				return err
			}

			securityGroupIDs[idx] = securityGroup.ID
		}

		data := compute.NetworkInterfaceSecurityGroupUpdate{
			SecurityGroupIDs: securityGroupIDs,
		}

		iface, err = service.UpdateSecurityGroups(cmd.Context(), iface.ID, data)
		if err != nil {
			return fmt.Errorf("update network interface security groups: %w", err)
		}
	}

	return commands.PrintStdout(iface)
}

func (n *networkInterfaceUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		server, err := findServer(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeNetworkInterface(cmd.Context(), server, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkInterfaceUpdateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update SERVER NETWORK-INTERFACE",
		Short:             "Update a network interface",
		Long:              "Updates a network interface of the current server.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().BoolVar(&n.disableSecurity, "disable-security", false, "disable security for the network interface")
	cmd.Flags().BoolVar(&n.enableSecurity, "enable-security", false, "enable security for the network interface")
	cmd.Flags().StringSliceVar(&n.securityGroups, "security-group", nil, "the security groups to use")

	cmd.MarkFlagsMutuallyExclusive("disable-security", "enable-security")

	return cmd
}

type networkInterfaceDeleteCommand struct {
	force bool
}

func (n *networkInterfaceDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	iface, err := findNetworkInterface(cmd.Context(), server.ID, args[1])
	if err != nil {
		return err
	}

	if iface.AttachedElasticIP.ID != 0 {
		return fmt.Errorf("network interface still has an elastic ip attached to it")
	}

	if !n.force && !commands.ConfirmDeletion("network interface", iface) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NewNetworkInterfaceService(commands.Config.Client, server.ID).Delete(cmd.Context(), iface.ID)
	if err != nil {
		return fmt.Errorf("delete network interface: %w", err)
	}

	return nil
}

func (n *networkInterfaceDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		server, err := findServer(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeNetworkInterface(cmd.Context(), server, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkInterfaceDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete SERVER NETWORK-INTERFACE",
		Aliases:           []string{"remove", "rm", "delete", "del"},
		Short:             "Delete a network interface",
		Long:              "Deletes a network interface of the current server.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().BoolVar(&n.force, "force", false, "force the deletion of the network interface without asking for confirmation")

	return cmd
}

func completeNetworkInterface(ctx context.Context, server compute.Server, term string) ([]string, cobra.ShellCompDirective) {
	interfaces, err := compute.NewNetworkInterfaceService(commands.Config.Client, server.ID).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(interfaces, term)

	names := make([]string, len(filtered))
	for i, iface := range filtered {
		names[i] = iface.PrivateIP
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findNetworkInterface(ctx context.Context, serverID int, term string) (compute.NetworkInterface, error) {
	ifaces, err := compute.NewNetworkInterfaceService(commands.Config.Client, serverID).List(ctx)
	if err != nil {
		return compute.NetworkInterface{}, fmt.Errorf("fetch network interfaces: %w", err)
	}

	iface, err := filter.FindOne(ifaces, term)
	if err != nil {
		return compute.NetworkInterface{}, fmt.Errorf("find network interface: %w", err)
	}

	return iface, nil
}
