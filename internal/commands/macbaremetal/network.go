package macbaremetal

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func NetworkCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network",
		Aliases: []string{"networks"},
		Short:   "Manage mac bare metal networks",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # List all networks
      %[1]s mac-bare-metal network list
      
      # Create a new network
      %[1]s mac-bare-metal network create --name my-network --location ZRH1
		`, app.Name)),
	}

	commands.Add(app, cmd,
		&networkListCommand{},
		&networkCreateCommand{},
		&networkUpdateCommand{},
		&networkDeleteCommand{},
	)

	return cmd
}

type networkListCommand struct {
	filter string
}

func (n *networkListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := macbaremetal.NewNetworkService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch networks: %w", err)
	}

	if len(n.filter) != 0 {
		items = filter.Find(items, n.filter)
	}

	return commands.PrintStdout(items)
}

func (n *networkListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List networks",
		Long:              "Lists all mac bare metal networks.",
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().StringVar(&n.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type networkCreateCommand struct {
	name        string
	description string
	location    string
}

func (n *networkCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Config.Client, n.location)
	if err != nil {
		return err
	}

	data := macbaremetal.NetworkCreate{
		Name:        n.name,
		Description: n.description,
		LocationID:  location.ID,
	}

	item, err := macbaremetal.NewNetworkService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create network: %w", err)
	}

	return commands.PrintStdout(item)
}

func (n *networkCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkCreateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Aliases:           []string{"add", "new"},
		Short:             "Create new network",
		Long:              "Creates a new mac bare metal network.",
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
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

func (n *networkUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	service := macbaremetal.NewNetworkService(commands.Config.Client)

	networks, err := service.List(cmd.Context())
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

	network, err = service.Update(cmd.Context(), network.ID, update)
	if err != nil {
		return fmt.Errorf("update network: %w", err)
	}

	return commands.PrintStdout(network)
}

func (n *networkUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeNetwork(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkUpdateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update NETWORK",
		Short:             "Update network",
		Long:              "Updates a mac bare metal network.",
		Example:           "", // TODO
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
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

func (n *networkDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	service := macbaremetal.NewNetworkService(commands.Config.Client)

	networks, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch networks: %w", err)
	}

	network, err := filter.FindOne(networks, args[0])
	if err != nil {
		return fmt.Errorf("find network: %w", err)
	}

	if !n.force && !commands.ConfirmDeletion("network", network) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = service.Delete(cmd.Context(), network.ID)
	if err != nil {
		return fmt.Errorf("delete network: %w", err)
	}

	return nil
}

func (n *networkDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeNetwork(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (n *networkDeleteCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete NETWORK",
		Aliases:           []string{"del", "remove", "rm"},
		Short:             "Delete network",
		Long:              "Deletes a mac bare metal network.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: n.CompleteArg,
		RunE:              n.Run,
	}

	cmd.Flags().BoolVar(&n.force, "force", false, "force the deletion of the network without asking for confirmation")

	return cmd
}

func completeNetwork(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	networks, err := macbaremetal.NewNetworkService(commands.Config.Client).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(networks, term)

	names := make([]string, len(filtered))
	for i, network := range filtered {
		names[i] = network.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
