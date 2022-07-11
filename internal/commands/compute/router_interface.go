package compute

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func RouterInterfaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interface",
		Short: "Manage compute router interfaces",
	}

	commands.Add(cmd, &routerInterfaceListCommand{}, &routerInterfaceCreateCommand{}, &routerInterfaceDeleteCommand{})

	return cmd
}

type routerInterfaceListCommand struct {
	filter string
}

func (r *routerInterfaceListCommand) Run(cmd *cobra.Command, args []string) error {
	router, err := findRouter(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	items, err := compute.NewRouterInterfaceService(commands.Config.Client, router.ID).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch routers: %w", err)
	}

	if len(r.filter) != 0 {
		items = filter.Find(items, r.filter)
	}

	return commands.PrintStdout(items)
}

func (r *routerInterfaceListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list ROUTER",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List router interfaces",
		Long:    "Lists all router interfaces of the selected router.",
		Args:    cobra.ExactArgs(1),
		RunE:    r.Run,
	}

	cmd.Flags().StringVar(&r.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type routerInterfaceCreateCommand struct {
	network   string
	privateIP net.IP
}

func (r *routerInterfaceCreateCommand) Run(cmd *cobra.Command, args []string) error {
	router, err := findRouter(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	network, err := findNetwork(cmd.Context(), r.network)
	if err != nil {
		return err
	}

	data := compute.RouterInterfaceCreate{
		NetworkID: network.ID,
	}

	if len(r.privateIP) != 0 {
		_, cidr, err := net.ParseCIDR(network.CIDR)
		if err != nil {
			return fmt.Errorf("parse network cidr: %w", err)
		}

		if !cidr.Contains(r.privateIP) {
			return fmt.Errorf("private ip %s is not in network %s", r.privateIP, network.CIDR)
		}

		data.PrivateIP = r.privateIP.String()
	}

	item, err := compute.NewRouterInterfaceService(commands.Config.Client, router.ID).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create router interface: %w", err)
	}

	return commands.PrintStdout(item)
}

func (r *routerInterfaceCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create ROUTER",
		Short: "Create a router interface",
		Long:  "Creates a new router interface.",
		Args:  cobra.ExactArgs(1),
		RunE:  r.Run,
	}

	cmd.Flags().StringVar(&r.network, "network", "", "the network to use")
	cmd.Flags().IPVar(&r.privateIP, "private-ip", nil, "the private IP to use")

	_ = cmd.MarkFlagRequired("network")

	return cmd
}

type routerInterfaceDeleteCommand struct {
	force bool
}

func (r *routerInterfaceDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	router, err := findRouter(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	routerInterfaces, err := compute.NewRouterInterfaceService(commands.Config.Client, router.ID).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch router interfaces: %w", err)
	}

	routerInterface, err := filter.FindOne(routerInterfaces, args[1])
	if err != nil {
		return fmt.Errorf("find router interface: %w", err)
	}

	if !r.force && !commands.ConfirmDeletion("router interface", routerInterface) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NewRouterInterfaceService(commands.Config.Client, router.ID).Delete(cmd.Context(), routerInterface.ID)
	if err != nil {
		return fmt.Errorf("delete router interface: %w", err)
	}

	return nil
}

func (r *routerInterfaceDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete ROUTER INTERFACE",
		Short: "Delete router interface",
		Long:  "Deletes a compute router interface.",
		Args:  cobra.ExactArgs(2),
		RunE:  r.Run,
	}

	cmd.Flags().BoolVar(&r.force, "force", false, "force the deletion of the router interface without asking for confirmation")

	return cmd
}
