package compute

import (
	"context"
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func DeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device",
		Aliases: []string{"devices"},
		Short:   "Manage mac bare metal devices",
		Example: commands.FormatExamples(fmt.Sprintf(`
			# Create a new device
			%[1]s mac-bare-metal device create --name "my-device" --product "macmini.2018.6-16-256" --network default --password "some-secure-user-password"
		`, commands.Name)),
	}

	commands.Add(cmd, &deviceListCommand{}, &deviceCreateCommand{}, &deviceUpdateCommand{}, &deviceDeleteCommand{}, &deviceVNCCommand{})
	cmd.AddCommand(DeviceActionCommand(), DeviceWorkflowCommand(), NetworkInterfaceCommands())

	commands.Add(cmd,
		deviceActionRunCommandPreset("power-off"),
		deviceActionRunCommandPreset("power-on"),
		deviceActionRunCommandPreset("power-cord-un-plug"),
		deviceActionRunCommandPreset("power-cord-plug-in"),
	)

	return cmd
}

type deviceListCommand struct {
	filter string
}

func (d *deviceListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := macbaremetal.NewDeviceService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch devices: %w", err)
	}

	if len(d.filter) != 0 {
		items = filter.Find(items, d.filter)
	}

	return commands.PrintStdout(items)
}

func (d *deviceListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d *deviceListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List devices",
		Long:    "Prints a table of all mac bare metal devices belonging to the current organization.",
		Example: commands.FormatExamples(fmt.Sprintf(`
			# Print all devices
			%[1]s mac-bare-metal device list

			# Print all devices containing the name "device"
			%[1]s mac-bare-metal device list --filter "device"

			# Print all devices in JSON format
			%[1]s mac-bare-metal device list --format json
		`, commands.Name)), // TODO
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}

	cmd.Flags().StringVar(&d.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type deviceVNCCommand struct {
	open bool
}

func (d *deviceVNCCommand) Run(cmd *cobra.Command, args []string) error {
	device, err := findDevice(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	vnc, err := macbaremetal.NewDeviceService(commands.Config.Client).GetVNC(cmd.Context(), device.ID)
	if err != nil {
		return fmt.Errorf("fetch vnc connection: %w", err)
	}

	if d.open {
		if err = browser.OpenURL(vnc.Ref); err != nil {
			return fmt.Errorf("open vnc connection: %w", err)
		}
	} else {
		commands.Stdout.Println(vnc.Ref)
	}

	return nil
}

func (d *deviceVNCCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeDevice(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d *deviceVNCCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vnc DEVICE",
		Short: "Get VNC of device",
		Long:  "Prints the VNC url of the device.",
		Example: commands.FormatExamples(fmt.Sprintf(`
			# Print the VNC url of the device "my-device"
			%[1]s mac-bare-metal device vnc my-device

			# Open the VNC url of the device "my-device" in the browser
			%[1]s mac-bare-metal device vnc my-device --open
		`, commands.Name)),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}

	cmd.Flags().BoolVar(&d.open, "open", false, "open the VNC url in the browser")

	return cmd
}

type deviceCreateCommand struct {
	name            string
	product         string
	network         string
	attachElasticIP bool
	password        string
}

func (d *deviceCreateCommand) Run(cmd *cobra.Command, args []string) error {
	products, err := common.ProductsByType(cmd.Context(), commands.Config.Client, common.ProductTypeMacBareMetalDevice)
	if err != nil {
		return fmt.Errorf("fetch products: %w", err)
	}

	product, err := filter.FindOne(products, d.product)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}

	networks, err := macbaremetal.NewNetworkService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch networks: %w", err)
	}

	network, err := filter.FindOne(networks, d.network)
	if err != nil {
		return fmt.Errorf("find network: %w", err)
	}

	data := macbaremetal.DeviceCreate{
		Name:            d.name,
		LocationID:      network.Location.ID,
		ProductID:       product.ID,
		NetworkID:       network.ID,
		AttachElasticIP: d.attachElasticIP,
		Password:        d.password,
	}

	service := macbaremetal.NewDeviceService(commands.Config.Client)

	ordering, err := service.Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create device: %w", err)
	}

	order, err := commands.WaitForOrder(cmd.Context(), "Creating device", ordering)
	if err != nil {
		return fmt.Errorf("wait for order: %w", err)
	}

	device, err := service.Get(cmd.Context(), order.Product.ID)
	if err != nil {
		return fmt.Errorf("fetch order: %w", err)
	}

	return commands.PrintStdout(device)
}

func (d *deviceCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d *deviceCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Aliases:           []string{"add", "new"},
		Short:             "Create new device",
		Long:              "Creates a new mac bare metal device.",
		Example:           "", // TODO
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}

	cmd.Flags().StringVar(&d.name, "name", "", "name to be applied to the device")
	cmd.Flags().StringVar(&d.product, "product", "", "product for the device")
	cmd.Flags().StringVar(&d.network, "network", "", "network to be attached to the device")
	cmd.Flags().BoolVar(&d.attachElasticIP, "attach-elastic-ip", false, "whether to attach an elastic ip to the device")
	cmd.Flags().StringVar(&d.password, "password", "", "password to be applied to the device") // TODO this is insecure and should be removed

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("product")
	_ = cmd.MarkFlagRequired("network")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}

type deviceUpdateCommand struct {
	name string
}

func (d *deviceUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	service := macbaremetal.NewDeviceService(commands.Config.Client)

	devices, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch devices: %w", err)
	}

	device, err := filter.FindOne(devices, args[0])
	if err != nil {
		return fmt.Errorf("find device: %w", err)
	}

	update := macbaremetal.DeviceUpdate{
		Name: d.name,
	}

	device, err = service.Update(cmd.Context(), device.ID, update)
	if err != nil {
		return fmt.Errorf("update device: %w", err)
	}

	return commands.PrintStdout(device)
}

func (d *deviceUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeDevice(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d *deviceUpdateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update DEVICE",
		Short:             "Update device",
		Long:              "Updates a mac bare metal device.",
		Example:           "", // TODO
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}

	cmd.Flags().StringVar(&d.name, "name", "", "name to be applied to the device")

	_ = cmd.MarkFlagRequired("device")

	return cmd
}

type deviceDeleteCommand struct {
	force bool
}

func (d *deviceDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	service := macbaremetal.NewDeviceService(commands.Config.Client)

	devices, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch devices: %w", err)
	}

	device, err := filter.FindOne(devices, args[0])
	if err != nil {
		return fmt.Errorf("find device: %w", err)
	}

	if !d.force && !commands.ConfirmDeletion("device", device) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = service.Delete(cmd.Context(), device.ID)
	if err != nil {
		return fmt.Errorf("delete device: %w", err)
	}

	return nil
}

func (d *deviceDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeDevice(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d *deviceDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete DEVICE",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete device",
		Long:    "Deletes a mac bare metal device.",
		Example: commands.FormatExamples(fmt.Sprintf(`
			# Delete a device
			%[1]s mac-bare-metal device delete my-device

			# Force the deletion of a device (without confirmation)
			%[1]s mac-bare-metal device delete my-device --force
		`, commands.Name)),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}

	cmd.Flags().BoolVar(&d.force, "force", false, "force the deletion of the device without asking for confirmation")

	return cmd
}

func completeDevice(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	devices, err := macbaremetal.NewDeviceService(commands.Config.Client).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(devices, term)

	names := make([]string, len(filtered))
	for i, d := range filtered {
		names[i] = d.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findDevice(ctx context.Context, term string) (macbaremetal.Device, error) {
	elasticIPs, err := macbaremetal.NewDeviceService(commands.Config.Client).List(ctx)
	if err != nil {
		return macbaremetal.Device{}, fmt.Errorf("fetch devices: %w", err)
	}

	elasticIP, err := filter.FindOne(elasticIPs, term)
	if err != nil {
		return macbaremetal.Device{}, fmt.Errorf("find device: %w", err)
	}

	return elasticIP, nil
}
