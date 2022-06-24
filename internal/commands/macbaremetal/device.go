package compute

import (
	"context"
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/console"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func DeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device",
		Short:   "Manage mac bare metal devices",
		Example: "", // TODO
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

func (d *deviceListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	items, err := macbaremetal.NewDeviceService(config.Client).List(ctx)
	if err != nil {
		return fmt.Errorf("fetch devices: %w", err)
	}

	if len(d.filter) != 0 {
		items = filter.Find(items, d.filter)
	}

	return commands.PrintStdout(items)
}

func (d *deviceListCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List devices",
		Long:    "Lists all mac bare metal devices.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&d.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type deviceVNCCommand struct {
	open bool
}

func (d *deviceVNCCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	device, err := findDevice(ctx, config, args[0])
	if err != nil {
		return err
	}

	vnc, err := macbaremetal.NewDeviceService(config.Client).GetVNC(ctx, device.ID)
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

func (d *deviceVNCCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vnc DEVICE",
		Short:   "Get VNC of device",
		Long:    "Returns the VNC url of the device.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
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

func (d *deviceCreateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	products, err := common.ProductsByType(ctx, config.Client, common.ProductTypeMacBareMetal)
	if err != nil {
		return fmt.Errorf("fetch products: %w", err)
	}

	product, err := filter.FindOne(products, d.product)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}

	networks, err := macbaremetal.NewNetworkService(config.Client).List(ctx)
	if err != nil {
		return fmt.Errorf("fetch networks: %w", err)
	}

	network, err := filter.FindOne(networks, d.network)
	if err != nil {
		return fmt.Errorf("find network: %w", err)
	}

	data := macbaremetal.DeviceCreate{
		Name:            d.name,
		LocationID:      network.Location.Id,
		ProductID:       product.Id,
		NetworkID:       network.ID,
		AttachElasticIP: d.attachElasticIP,
		Password:        d.password,
	}

	ordering, err := macbaremetal.NewDeviceService(config.Client).Create(ctx, data)
	if err != nil {
		return fmt.Errorf("create device: %w", err)
	}

	progress := console.NewProgress("Creating device")
	go progress.Display(commands.Stderr)

	err = common.WaitForOrder(ctx, config.Client, ordering)
	if err != nil {
		progress.Complete("Order filed")

		return fmt.Errorf("wait for order: %w", err)
	}

	progress.Complete("Order completed")

	// TODO find device created through order and print it

	return nil
}

func (d *deviceCreateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create new device",
		Long:    "Creates a new mac bare metal device.",
		Example: "", // TODO
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

func (d *deviceUpdateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	service := macbaremetal.NewDeviceService(config.Client)

	devices, err := service.List(ctx)
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

	device, err = service.Update(ctx, device.ID, update)
	if err != nil {
		return fmt.Errorf("update device: %w", err)
	}

	return commands.PrintStdout(device)
}

func (d *deviceUpdateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update DEVICE",
		Short:   "Update device",
		Long:    "Updates a mac bare metal device.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&d.name, "name", "", "name to be applied to the device")

	_ = cmd.MarkFlagRequired("device")

	return cmd
}

type deviceDeleteCommand struct {
	force bool
}

func (d *deviceDeleteCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	service := macbaremetal.NewDeviceService(config.Client)

	devices, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch devices: %w", err)
	}

	device, err := filter.FindOne(devices, args[0])
	if err != nil {
		return fmt.Errorf("find device: %w", err)
	}

	// TODO ask for confirmation

	err = service.Delete(ctx, device.ID)
	if err != nil {
		return fmt.Errorf("delete device: %w", err)
	}

	return nil
}

func (d *deviceDeleteCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete DEVICE",
		Short:   "Delete device",
		Long:    "Deletes a mac bare metal device.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
	}

	return cmd
}

func findDevice(ctx context.Context, config commands.Config, term string) (macbaremetal.Device, error) {
	elasticIPs, err := macbaremetal.NewDeviceService(config.Client).List(ctx)
	if err != nil {
		return macbaremetal.Device{}, fmt.Errorf("fetch devices: %w", err)
	}

	elasticIP, err := filter.FindOne(elasticIPs, term)
	if err != nil {
		return macbaremetal.Device{}, fmt.Errorf("find device: %w", err)
	}

	return elasticIP, nil
}
