package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func NetworkInterfaceCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "interface",
		Short:   "Manage network interfaces",
		Example: "", // TODO
	}

	commands.Add(cmd, &networkInterfaceListCommand{}, &networkInterfaceUpdateCommand{})

	return cmd
}

type networkInterfaceListCommand struct {
}

func (n *networkInterfaceListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	device, err := findDevice(ctx, config, args[0])
	if err != nil {
		return err
	}

	interfaces, err := macbaremetal.NewNetworkInterfaceService(config.Client, device.ID).List(ctx)
	if err != nil {
		return fmt.Errorf("fetch network interfaces: %w", err)
	}

	return commands.PrintStdout(interfaces)
}

func (n *networkInterfaceListCommand) Desc() *cobra.Command {
	return &cobra.Command{
		Use:     "list DEVICE",
		Short:   "List network interfaces",
		Long:    "Lists all network interfaces of a device.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
	}
}

type networkInterfaceUpdateCommand struct {
	securityGroup string
}

func (n *networkInterfaceUpdateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	device, err := findDevice(ctx, config, args[0])
	if err != nil {
		return err
	}

	interfaces, err := macbaremetal.NewNetworkInterfaceService(config.Client, device.ID).List(ctx)
	if err != nil {
		return fmt.Errorf("fetch network interfaces: %w", err)
	}

	iface, err := filter.FindOne(interfaces, args[1])
	if err != nil {
		return fmt.Errorf("find network interface: %w", err)
	}

	if len(n.securityGroup) > 0 {
		securityGroup, err := findSecurityGroup(ctx, config, n.securityGroup)
		if err != nil {
			return fmt.Errorf("find security group: %w", err)
		}

		update := macbaremetal.NetworkInterfaceSecurityGroupUpdate{
			SecurityGroupID: securityGroup.ID,
		}

		iface, err = macbaremetal.NewNetworkInterfaceService(config.Client, device.ID).UpdateSecurityGroup(ctx, iface.ID, update)
		if err != nil {
			return fmt.Errorf("update network interface security group: %w", err)
		}
	}

	return commands.PrintStdout(iface)
}

func (n *networkInterfaceUpdateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update DEVICE INTERFACE",
		Short:   "Update network interface",
		Long:    "Updates a network interface of a device.",
		Args:    cobra.ExactArgs(2),
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&n.securityGroup, "security-group", "", "security group to be applied to the network interface")

	return cmd
}
