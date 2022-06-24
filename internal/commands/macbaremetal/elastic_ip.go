package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ElasticIPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elastic-ip",
		Short:   "Manage mac bare metal elastic ips",
		Example: "", // TODO
	}

	commands.Add(cmd, &elasticIPListCommand{}, &elasticIPCreateCommand{}, &elasticIPDeleteCommand{}, &elasticIPAttachCommand{}, &elasticIPDetachCommand{})

	return cmd
}

type elasticIPListCommand struct {
	filter string
}

func (e *elasticIPListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	items, err := macbaremetal.NewElasticIPService(config.Client).List(ctx)
	if err != nil {
		return fmt.Errorf("fetch elastic ips: %w", err)
	}

	if len(e.filter) != 0 {
		items = filter.Find(items, e.filter)
	}

	return commands.PrintStdout(items)
}

func (e *elasticIPListCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List elastic ips",
		Long:    "Lists all mac bare metal elastic ips.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&e.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type elasticIPCreateCommand struct {
	location string
}

func (e *elasticIPCreateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	locations, err := common.Locations(ctx, config.Client)
	if err != nil {
		return fmt.Errorf("fetch locations: %w", err)
	}

	location, err := filter.FindOne(locations, e.location)
	if err != nil {
		return fmt.Errorf("find location: %w", err)
	}

	data := macbaremetal.ElasticIPCreate{
		LocationID: location.Id,
	}

	item, err := macbaremetal.NewElasticIPService(config.Client).Create(ctx, data)
	if err != nil {
		return fmt.Errorf("create elastic ip: %w", err)
	}

	return commands.PrintStdout(item)
}

func (e *elasticIPCreateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create new elastic ip",
		Long:    "Creates a new mac bare metal elastic ip.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&e.location, "location", "", "location where the elastic ip will be created")
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

type elasticIPDeleteCommand struct {
	force bool
}

func (e *elasticIPDeleteCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	service := macbaremetal.NewElasticIPService(config.Client)

	elasticIPs, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch elastic ips: %w", err)
	}

	elasticIP, err := filter.FindOne(elasticIPs, args[0])
	if err != nil {
		return fmt.Errorf("find elastic ip: %w", err)
	}

	// TODO ask for confirmation

	err = service.Delete(ctx, elasticIP.ID)
	if err != nil {
		return fmt.Errorf("delete elastic ip: %w", err)
	}

	return nil
}

func (e *elasticIPDeleteCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete ELASTIC-IP",
		Short:   "Delete elastic ip",
		Long:    "Deletes a mac bare metal elastic ip.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
	}

	return cmd
}

type elasticIPAttachCommand struct {
}

func (e *elasticIPAttachCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	elasticIP, err := findElasticIP(ctx, config, args[0])
	if err != nil {
		return err
	}

	if elasticIP.Attachment.ID != 0 {
		return fmt.Errorf("elastic ip is already attached to a device")
	}

	device, err := findDevice(ctx, config, args[1])
	if err != nil {
		return err
	}

	networkInterfaceID := 0
	for _, networkInterface := range device.NetworkInterfaces {
		if networkInterface.PublicIP == "" {
			networkInterfaceID = networkInterface.ID
			break
		}
	}

	if networkInterfaceID == 0 {
		return fmt.Errorf("device has no free network interface to attach the elastic ip to")
	}

	body := macbaremetal.ElasticIPAttach{
		ElasticIPID:        elasticIP.ID,
		NetworkInterfaceID: networkInterfaceID,
	}

	_, err = macbaremetal.NewElasticIPService(config.Client).Attach(ctx, device.ID, body)
	if err != nil {
		return fmt.Errorf("attach elastic ip: %w", err)
	}

	return nil
}

func (e *elasticIPAttachCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "attach ELASTIC-IP DEVICE",
		Short:   "Attach elastic ip to device",
		Long:    "Attaches a mac bare metal elastic ip to a device.",
		Args:    cobra.ExactArgs(2),
		Example: "", // TODO
	}

	return cmd
}

type elasticIPDetachCommand struct {
}

func (e *elasticIPDetachCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	elasticIP, err := findElasticIP(ctx, config, args[0])
	if err != nil {
		return err
	}

	device, err := findDevice(ctx, config, args[1])
	if err != nil {
		return err
	}

	if elasticIP.Attachment.ID != device.ID {
		return fmt.Errorf("elastic ip not attached to the selected device")
	}

	err = macbaremetal.NewElasticIPService(config.Client).Detach(ctx, device.ID, elasticIP.ID)
	if err != nil {
		return fmt.Errorf("detach elastic ip: %w", err)
	}

	return nil
}

func (e *elasticIPDetachCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "detach ELASTIC-IP DEVICE",
		Short:   "Detach elastic ip from device",
		Long:    "Detaches a mac bare metal elastic ip from a device.",
		Args:    cobra.ExactArgs(2),
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

func findElasticIP(ctx context.Context, config commands.Config, term string) (macbaremetal.ElasticIP, error) {
	elasticIPs, err := macbaremetal.NewElasticIPService(config.Client).List(ctx)
	if err != nil {
		return macbaremetal.ElasticIP{}, fmt.Errorf("fetch elastic ips: %w", err)
	}

	elasticIP, err := filter.FindOne(elasticIPs, term)
	if err != nil {
		return macbaremetal.ElasticIP{}, fmt.Errorf("find elastic ip: %w", err)
	}

	return elasticIP, nil
}
