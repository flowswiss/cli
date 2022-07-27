package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func RouterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "router",
		Aliases: []string{"routers"},
		Short:   "Manage mac bare metal routers",
	}

	commands.Add(cmd, &routerListCommand{}, &routerUpdateCommand{})

	return cmd
}

type routerListCommand struct {
	filter string
}

func (r *routerListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := macbaremetal.NewRouterService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch routers: %w", err)
	}

	if len(r.filter) != 0 {
		items = filter.Find(items, r.filter)
	}

	return commands.PrintStdout(items)
}

func (r *routerListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (r *routerListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List routers",
		Long:              "Lists all mac bare metal routers.",
		ValidArgsFunction: r.CompleteArg,
		RunE:              r.Run,
	}

	cmd.Flags().StringVar(&r.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type routerUpdateCommand struct {
	name        string
	description string
}

func (r *routerUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	service := macbaremetal.NewRouterService(commands.Config.Client)

	routers, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch routers: %w", err)
	}

	router, err := filter.FindOne(routers, args[0])
	if err != nil {
		return fmt.Errorf("find router: %w", err)
	}

	update := macbaremetal.RouterUpdate{
		Name:        r.name,
		Description: r.description,
	}

	router, err = service.Update(cmd.Context(), router.ID, update)
	if err != nil {
		return fmt.Errorf("update router: %w", err)
	}

	return commands.PrintStdout(router)
}

func (r *routerUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeRouter(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (r *routerUpdateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update ROUTER",
		Short:             "Update router",
		Long:              "Updates a mac bare metal router.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: r.CompleteArg,
		RunE:              r.Run,
	}

	cmd.Flags().StringVar(&r.name, "name", "", "name to be applied to the router")
	cmd.Flags().StringVar(&r.description, "description", "", "description to be applied to the router")

	return cmd
}

func completeRouter(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	routers, err := macbaremetal.NewRouterService(commands.Config.Client).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(routers, term)

	names := make([]string, len(filtered))
	for i, router := range filtered {
		names[i] = router.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
