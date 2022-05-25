package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/internal/commands"
	"github.com/flowswiss/cli/pkg/api/common"
	"github.com/flowswiss/cli/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/pkg/filter"
)

func ElasticIPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elastic-ip",
		Short:   "Manage mac bare metal elastic ips",
		Example: "", // TODO
	}

	commands.Add(cmd, &elasticIPListCommand{}, &elasticIPCreateCommand{}, &elasticIPDeleteCommand{})

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
	elasticIP string
	force     bool
}

func (e *elasticIPDeleteCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	service := macbaremetal.NewElasticIPService(config.Client)

	elasticIPs, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch elastic ips: %w", err)
	}

	elasticIP, err := filter.FindOne(elasticIPs, e.elasticIP)
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
		Use:     "delete",
		Short:   "Delete elastic ip",
		Long:    "Deletes a mac bare metal elastic ip.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&e.elasticIP, "elastic-ip", "", "elastic ip to be deleted")
	_ = cmd.MarkFlagRequired("elastic-ip")

	return cmd
}
