package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/console"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ElasticIPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elastic-ip",
		Aliases: []string{"elastic-ips", "elasticip", "elasticips"},
		Short:   "Manage mac bare metal elastic ips",
		Example: commands.FormatExamples(fmt.Sprintf(`
  			# List all mac bare metal elastic ips
	  		%[1]s mac-bare-metal elastic-ip list	

			# Create a new mac bare metal elastic ip
			%[1]s mac-bare-metal elastic-ip create --location=ZRH1

			# Attach a mac bare metal elastic ip to a device
			%[1]s mac-bare-metal elastic-ip attach 1.1.1.1 my-device
		`, commands.Name)),
	}

	commands.Add(cmd, &elasticIPListCommand{}, &elasticIPCreateCommand{}, &elasticIPDeleteCommand{}, &elasticIPAttachCommand{}, &elasticIPDetachCommand{})

	return cmd
}

type elasticIPListCommand struct {
	filter string
}

func (e *elasticIPListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := macbaremetal.NewElasticIPService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch elastic ips: %w", err)
	}

	if len(e.filter) != 0 {
		items = filter.Find(items, e.filter)
	}

	return commands.PrintStdout(items)
}

func (e *elasticIPListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List elastic ips",
		Long:    "Lists all mac bare metal elastic ips.",
		RunE:    e.Run,
	}

	cmd.Flags().StringVar(&e.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type elasticIPCreateCommand struct {
	location string
}

func (e *elasticIPCreateCommand) Run(cmd *cobra.Command, args []string) error {
	locations, err := common.Locations(cmd.Context(), commands.Config.Client)
	if err != nil {
		return fmt.Errorf("fetch locations: %w", err)
	}

	location, err := filter.FindOne(locations, e.location)
	if err != nil {
		return fmt.Errorf("find location: %w", err)
	}

	data := macbaremetal.ElasticIPCreate{
		LocationID: location.ID,
	}

	item, err := macbaremetal.NewElasticIPService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create elastic ip: %w", err)
	}

	return commands.PrintStdout(item)
}

func (e *elasticIPCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "new"},
		Short:   "Create new elastic ip",
		Long:    "Creates a new mac bare metal elastic ip.",
		RunE:    e.Run,
	}

	cmd.Flags().StringVar(&e.location, "location", "", "location where the elastic ip will be created")
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

type elasticIPDeleteCommand struct {
	force bool
}

func (e *elasticIPDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	service := macbaremetal.NewElasticIPService(commands.Config.Client)

	elasticIPs, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch elastic ips: %w", err)
	}

	elasticIP, err := filter.FindOne(elasticIPs, args[0])
	if err != nil {
		return fmt.Errorf("find elastic ip: %w", err)
	}

	if elasticIP.Attachment.ID != 0 {
		commands.Stderr.Errorf("WARNING: The elastic ip is still attached to a device. Connections to the device will be lost.\n")
	}

	if !e.force {
		if !console.Confirm(commands.Stderr, fmt.Sprintf("Are you sure you want to delete the elastic ip %q?", elasticIP.PublicIP)) {
			commands.Stderr.Println("aborted.")
			return nil
		}
	}

	if elasticIP.Attachment.ID != 0 {
		err = service.Detach(cmd.Context(), elasticIP.Attachment.ID, elasticIP.ID)
		if err != nil {
			return fmt.Errorf("detach elastic ip: %w", err)
		}
	}

	err = service.Delete(cmd.Context(), elasticIP.ID)
	if err != nil {
		return fmt.Errorf("delete elastic ip: %w", err)
	}

	return nil
}

func (e *elasticIPDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete ELASTIC-IP",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete elastic ip",
		Long:    "Deletes a mac bare metal elastic ip.",
		Example: commands.FormatExamples(fmt.Sprintf(`
	  		# Delete a mac bare metal elastic ip
			%[1]s mac-bare-metal elastic-ip delete 1.1.1.1

			# Force the deletion a mac bare metal elastic ip without confirmation
			%[1]s mac-bare-metal elastic-ip delete 1.1.1.1 --force
		`, commands.Name)),
		Args: cobra.ExactArgs(1),
		RunE: e.Run,
	}

	cmd.Flags().BoolVar(&e.force, "force", false, "force the deletion of the elastic ip without asking for confirmation")

	return cmd
}

type elasticIPAttachCommand struct {
}

func (e *elasticIPAttachCommand) Run(cmd *cobra.Command, args []string) error {
	elasticIP, err := findElasticIP(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if elasticIP.Attachment.ID != 0 {
		return fmt.Errorf("elastic ip is already attached to a device")
	}

	device, err := findDevice(cmd.Context(), args[1])
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

	_, err = macbaremetal.NewElasticIPService(commands.Config.Client).Attach(cmd.Context(), device.ID, body)
	if err != nil {
		return fmt.Errorf("attach elastic ip: %w", err)
	}

	return nil
}

func (e *elasticIPAttachCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach ELASTIC-IP DEVICE",
		Short: "Attach elastic ip to device",
		Long:  "Attaches a mac bare metal elastic ip to a device.",
		Args:  cobra.ExactArgs(2),
		RunE:  e.Run,
	}

	return cmd
}

type elasticIPDetachCommand struct {
	force bool
}

func (e *elasticIPDetachCommand) Run(cmd *cobra.Command, args []string) error {
	elasticIP, err := findElasticIP(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	device, err := findDevice(cmd.Context(), args[1])
	if err != nil {
		return err
	}

	if elasticIP.Attachment.ID != device.ID {
		return fmt.Errorf("elastic ip not attached to the selected device")
	}

	if !e.force {
		if !console.Confirm(commands.Stderr, fmt.Sprintf("Are you sure you want to detach the elastic ip %q? Any connection to the device will be lost.", elasticIP.PublicIP)) {
			commands.Stderr.Println("aborted.")
			return nil
		}
	}

	err = macbaremetal.NewElasticIPService(commands.Config.Client).Detach(cmd.Context(), device.ID, elasticIP.ID)
	if err != nil {
		return fmt.Errorf("detach elastic ip: %w", err)
	}

	return nil
}

func (e *elasticIPDetachCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach ELASTIC-IP DEVICE",
		Short: "Detach elastic ip from device",
		Long:  "Detaches a mac bare metal elastic ip from a device.",
		Args:  cobra.ExactArgs(2),
		RunE:  e.Run,
	}

	cmd.Flags().BoolVar(&e.force, "force", false, "force the detachment of the elastic ip without asking for confirmation")

	return cmd
}

func findElasticIP(ctx context.Context, term string) (macbaremetal.ElasticIP, error) {
	elasticIPs, err := macbaremetal.NewElasticIPService(commands.Config.Client).List(ctx)
	if err != nil {
		return macbaremetal.ElasticIP{}, fmt.Errorf("fetch elastic ips: %w", err)
	}

	elasticIP, err := filter.FindOne(elasticIPs, term)
	if err != nil {
		return macbaremetal.ElasticIP{}, fmt.Errorf("find elastic ip: %w", err)
	}

	return elasticIP, nil
}
