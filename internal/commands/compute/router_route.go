package compute

import (
	"context"
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func RouterRouteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "Manage compute router routes",
	}

	commands.Add(cmd, &routeListCommand{}, &routeCreateCommand{}, &routeDeleteCommand{})

	return cmd
}

type routeListCommand struct {
	filter string
}

func (r *routeListCommand) Run(cmd *cobra.Command, args []string) error {
	router, err := findRouter(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	items, err := compute.NewRouteService(commands.Config.Client, router.ID).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch routes: %w", err)
	}

	if len(r.filter) != 0 {
		items = filter.Find(items, r.filter)
	}

	return commands.PrintStdout(items)
}

func (r *routeListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeRouter(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (r *routeListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list ROUTER",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List routes",
		Long:              "Lists all routes of the selected router.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: r.CompleteArg,
		RunE:              r.Run,
	}

	cmd.Flags().StringVar(&r.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type routeCreateCommand struct {
	destination net.IPNet
	nextHop     net.IP
}

func (r *routeCreateCommand) Run(cmd *cobra.Command, args []string) error {
	router, err := findRouter(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	data := compute.RouteCreate{
		Destination: r.destination.String(),
		NextHop:     r.nextHop.String(),
	}

	item, err := compute.NewRouteService(commands.Config.Client, router.ID).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create route: %w", err)
	}

	return commands.PrintStdout(item)
}

func (r *routeCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeRouter(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (r *routeCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create ROUTER",
		Short:             "Create a route",
		Long:              "Creates a new route.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: r.CompleteArg,
		RunE:              r.Run,
	}

	cmd.Flags().IPNetVar(&r.destination, "destination", net.IPNet{}, "destination of the route")
	cmd.Flags().IPVar(&r.nextHop, "next-hop", net.IP{}, "next hop of the route")

	_ = cmd.MarkFlagRequired("destination")
	_ = cmd.MarkFlagRequired("next-hop")

	return cmd
}

type routeDeleteCommand struct {
	force bool
}

func (r *routeDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	router, err := findRouter(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	routes, err := compute.NewRouteService(commands.Config.Client, router.ID).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch routes: %w", err)
	}

	route, err := filter.FindOne(routes, args[1])
	if err != nil {
		return fmt.Errorf("find route: %w", err)
	}

	if !r.force && !commands.ConfirmDeletion("route", route) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NewRouteService(commands.Config.Client, router.ID).Delete(cmd.Context(), route.ID)
	if err != nil {
		return fmt.Errorf("delete route: %w", err)
	}

	return nil
}

func (r *routeDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeRouter(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		router, err := findRouter(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeRouterRoute(cmd.Context(), router, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (r *routeDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete ROUTER ROUTE",
		Short:             "Delete route",
		Long:              "Deletes a compute route.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: r.CompleteArg,
		RunE:              r.Run,
	}

	cmd.Flags().BoolVar(&r.force, "force", false, "force the deletion of the route without asking for confirmation")

	return cmd
}

func completeRouterRoute(ctx context.Context, router compute.Router, term string) ([]string, cobra.ShellCompDirective) {
	routes, err := compute.NewRouteService(commands.Config.Client, router.ID).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(routes, term)

	names := make([]string, len(filtered))
	for i, route := range filtered {
		names[i] = route.Destination
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
