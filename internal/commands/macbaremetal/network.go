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

func NetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network",
		Short:   "Manage mac bare metal networks",
		Example: "", // TODO
	}

	commands.Add(cmd, &networkListCommand{}, &networkCreateCommand{}, &networkUpdateCommand{}, &networkDeleteCommand{})

	return cmd
}

type networkListCommand struct {
	filter string
}

func (n *networkListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	items, err := macbaremetal.NewNetworkService(config.Client).List(ctx)
	if err != nil {
		return fmt.Errorf("fetch networks: %w", err)
	}

	if len(n.filter) != 0 {
		items = filter.Find(items, n.filter)
	}

	return commands.PrintStdout(items)
}

func (n *networkListCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List networks",
		Long:    "Lists all mac bare metal networks.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&n.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type networkCreateCommand struct {
	name        string
	description string
	location    string
}

func (n *networkCreateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	locations, err := common.Locations(ctx, config.Client)
	if err != nil {
		return fmt.Errorf("fetch locations: %w", err)
	}

	location, err := filter.FindOne(locations, n.location)
	if err != nil {
		return fmt.Errorf("find location: %w", err)
	}

	data := macbaremetal.NetworkCreate{
		Name:        n.name,
		Description: n.description,
		LocationID:  location.Id,
	}

	item, err := macbaremetal.NewNetworkService(config.Client).Create(ctx, data)
	if err != nil {
		return fmt.Errorf("create network: %w", err)
	}

	return commands.PrintStdout(item)
}

func (n *networkCreateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create new network",
		Long:    "Creates a new mac bare metal network.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&n.name, "name", "", "name to be applied to the network")
	cmd.Flags().StringVar(&n.description, "description", "", "description to be applied to the network")
	cmd.Flags().StringVar(&n.location, "location", "", "location where the network will be created")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

type networkUpdateCommand struct {
	name             string
	description      string
	domainName       string
	domainNameServer []string
}

func (n *networkUpdateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	service := macbaremetal.NewNetworkService(config.Client)

	networks, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch networks: %w", err)
	}

	network, err := filter.FindOne(networks, args[0])
	if err != nil {
		return fmt.Errorf("find network: %w", err)
	}

	update := macbaremetal.NetworkUpdate{
		Name:              n.name,
		Description:       n.description,
		DomainName:        n.domainName,
		DomainNameServers: n.domainNameServer,
	}

	network, err = service.Update(ctx, network.ID, update)
	if err != nil {
		return fmt.Errorf("update network: %w", err)
	}

	return commands.PrintStdout(network)
}

func (n *networkUpdateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update NETWORK",
		Short:   "Update network",
		Long:    "Updates a mac bare metal network.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&n.name, "name", "", "name to be applied to the network")
	cmd.Flags().StringVar(&n.description, "description", "", "description to be applied to the network")
	cmd.Flags().StringVar(&n.domainName, "domain-name", "", "domain name to be applied to the network")
	cmd.Flags().StringSliceVar(&n.domainNameServer, "domain-name-server", nil, "domain name server to be applied to the network")

	return cmd
}

type networkDeleteCommand struct {
	force bool
}

func (n *networkDeleteCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	service := macbaremetal.NewNetworkService(config.Client)

	networks, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch networks: %w", err)
	}

	network, err := filter.FindOne(networks, args[0])
	if err != nil {
		return fmt.Errorf("find network: %w", err)
	}

	// TODO ask for confirmation

	err = service.Delete(ctx, network.ID)
	if err != nil {
		return fmt.Errorf("delete network: %w", err)
	}

	return nil
}

func (n *networkDeleteCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete NETWORK",
		Short:   "Delete network",
		Long:    "Deletes a mac bare metal network.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
	}

	return cmd
}
