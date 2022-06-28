package compute

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func NetworkInterfaceCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "interface",
		Aliases: []string{"interfaces", "networkinterface", "networkinterfaces", "network-interface", "network-interfaces"},
		Short:   "Manage network interfaces",
		Example: commands.FormatExamples(fmt.Sprintf(`
  			# List network interfaces of a device
	  		%s mac-bare-metal device interface list my-device

			# Update security group of a network interface
			%s mac-bare-metal device interface update my-device 1.1.1.1 --security-group default
		`)),
	}

	commands.Add(cmd, &networkInterfaceListCommand{}, &networkInterfaceUpdateCommand{})

	return cmd
}

type networkInterfaceListCommand struct {
}

func (n *networkInterfaceListCommand) Run(cmd *cobra.Command, args []string) error {
	device, err := findDevice(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	interfaces, err := macbaremetal.NewNetworkInterfaceService(commands.Config.Client, device.ID).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch network interfaces: %w", err)
	}

	return commands.PrintStdout(interfaces)
}

func (n *networkInterfaceListCommand) Build() *cobra.Command {
	return &cobra.Command{
		Use:     "list DEVICE",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List network interfaces",
		Long:    "Lists all network interfaces of a device.",
		Args:    cobra.ExactArgs(1),
		RunE:    n.Run,
	}
}

type networkInterfaceUpdateCommand struct {
	securityGroup string
}

func (n *networkInterfaceUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	device, err := findDevice(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	interfaces, err := macbaremetal.NewNetworkInterfaceService(commands.Config.Client, device.ID).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch network interfaces: %w", err)
	}

	iface, err := filter.FindOne(interfaces, args[1])
	if err != nil {
		return fmt.Errorf("find network interface: %w", err)
	}

	if len(n.securityGroup) > 0 {
		securityGroup, err := findSecurityGroup(cmd.Context(), n.securityGroup)
		if err != nil {
			return fmt.Errorf("find security group: %w", err)
		}

		update := macbaremetal.NetworkInterfaceSecurityGroupUpdate{
			SecurityGroupID: securityGroup.ID,
		}

		iface, err = macbaremetal.NewNetworkInterfaceService(commands.Config.Client, device.ID).UpdateSecurityGroup(cmd.Context(), iface.ID, update)
		if err != nil {
			return fmt.Errorf("update network interface security group: %w", err)
		}
	}

	return commands.PrintStdout(iface)
}

func (n *networkInterfaceUpdateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update DEVICE INTERFACE",
		Short: "Update network interface",
		Long:  "Updates a network interface of a device.",
		Args:  cobra.ExactArgs(2),
		RunE:  n.Run,
	}

	cmd.Flags().StringVar(&n.securityGroup, "security-group", "", "security group to be applied to the network interface")

	return cmd
}
