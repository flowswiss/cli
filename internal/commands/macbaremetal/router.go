package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/internal/commands"
	"github.com/flowswiss/cli/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/pkg/filter"
)

func RouterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "router",
		Short:   "Manage mac bare metal routers",
		Example: "", // TODO
	}

	commands.Add(cmd, &routerListCommand{}, &routerUpdateCommand{})

	return cmd
}

type routerListCommand struct {
	filter string
}

func (r *routerListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	items, err := macbaremetal.NewRouterService(config.Client).List(ctx)
	if err != nil {
		return fmt.Errorf("fetch locations: %w", err)
	}

	if len(r.filter) != 0 {
		items = filter.Find(items, r.filter)
	}

	return commands.PrintStdout(items)
}

func (r *routerListCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List routers",
		Long:    "Lists all mac bare metal routers.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&r.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type routerUpdateCommand struct {
	router      string
	name        string
	description string
}

func (r *routerUpdateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	service := macbaremetal.NewRouterService(config.Client)

	routers, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch routers: %w", err)
	}

	router, err := filter.FindOne(routers, r.router)
	if err != nil {
		return fmt.Errorf("find router: %w", err)
	}

	update := macbaremetal.RouterUpdate{
		Name:        r.name,
		Description: r.description,
	}

	router, err = service.Update(ctx, router.ID, update)
	if err != nil {
		return fmt.Errorf("update router: %w", err)
	}

	return commands.PrintStdout(router)
}

func (r *routerUpdateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update router",
		Long:    "Updates a mac bare metal router.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&r.router, "router", "", "router to be updated")
	cmd.Flags().StringVar(&r.name, "name", "", "name to be applied to the router")
	cmd.Flags().StringVar(&r.description, "description", "", "description to be applied to the router")

	_ = cmd.MarkFlagRequired("router")

	return cmd
}
