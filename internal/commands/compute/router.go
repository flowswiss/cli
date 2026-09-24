package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/pkg/optional"
	"github.com/spf13/cobra"

	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func RouterCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "router",
		Short: "Manage compute routers",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # List all routers
      %[1]s compute router list
      
      # Create a new router
      %[1]s compute router create --name my-router --location ALP1
		`, app.Name)),
	}

	commands.Add(app, cmd,
		&routerListCommand{},
		&routerCreateCommand{},
		&routerUpdateCommand{},
		&routerDeleteCommand{},
	)

	cmd.AddCommand(RouterInterfaceCommand(app), RouterRouteCommand(app))

	return cmd
}

type routerListCommand struct {
	filter string
}

func (r *routerListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.RouterService().List(cmd.Context(), core.CursorAll)
	if err != nil {
		return fmt.Errorf("fetch routers: %w", err)
	}

	if len(r.filter) != 0 {
		items = filter.Find(items, r.filter)
	}

	return commands.PrintStdout(items)
}

func (r *routerListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (r *routerListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List routers",
		Long:              "Lists all routers of the current tenant.",
		ValidArgsFunction: r.CompleteArg,
		RunE:              r.Run,
	}

	cmd.Flags().StringVar(&r.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type routerCreateCommand struct {
	name        string
	description optional.Optional[string]
	location    string
	private     bool
}

func (r *routerCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Client, r.location)
	if err != nil {
		return err
	}

	data := compute.RouterCreate{
		Name:        r.name,
		Description: r.description.Value(),
		LocationID:  location.ID,
		Public:      !r.private,
	}

	item, err := compute.RouterService().Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create router: %w", err)
	}

	return commands.PrintStdout(item)
}

func (r *routerCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (r *routerCreateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Short:             "Create a router",
		Long:              "Creates a new router",
		ValidArgsFunction: r.CompleteArg,
		RunE:              r.Run,
	}

	cmd.Flags().StringVar(&r.name, "name", "", "name of the router")
	cmd.Flags().StringVar(r.description.Configure(cmd, "description", "description of the router"))
	cmd.Flags().StringVar(&r.location, "location", "", "location of the router")
	cmd.Flags().BoolVar(&r.private, "private", false, "create a private router")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

type routerUpdateCommand struct {
	name        optional.Optional[string]
	description optional.Optional[string]
	makePrivate bool
	makePublic  bool
}

func (r *routerUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	router, err := findRouter(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	data := compute.RouterUpdate{ID: uint(router.ID),
		Name:        r.name.Value(),
		Description: r.description.Value(),
	}

	if r.makePrivate {
		data.Public = new(false)
	}

	if r.makePublic {
		data.Public = new(true)
	}

	router, err = compute.RouterService().Update(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("update router: %w", err)
	}

	return commands.PrintStdout(router)
}

func (r *routerUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeRouter(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (r *routerUpdateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update ROUTER",
		Short:             "Update router",
		Long:              "Updates a compute router.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: r.CompleteArg,
		RunE:              r.Run,
	}

	cmd.Flags().StringVar(r.name.Configure(cmd, "name", "name of the router"))
	cmd.Flags().StringVar(r.description.Configure(cmd, "description", "description of the router"))
	cmd.Flags().BoolVar(&r.makePrivate, "private", false, "make router private")
	cmd.Flags().BoolVar(&r.makePublic, "public", false, "make router public")

	cmd.MarkFlagsMutuallyExclusive("private", "public")

	return cmd
}

type routerDeleteCommand struct {
	force bool
}

func (r *routerDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	router, err := findRouter(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !r.force && !commands.ConfirmDeletion("router", router) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.RouterService().Delete(cmd.Context(), compute.RouterDelete{ID: uint(router.ID)})
	if err != nil {
		return fmt.Errorf("delete router: %w", err)
	}

	return nil
}

func (r *routerDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeRouter(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (r *routerDeleteCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete ROUTER",
		Short:             "Delete router",
		Long:              "Deletes a compute router.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: r.CompleteArg,
		RunE:              r.Run,
	}

	cmd.Flags().BoolVar(&r.force, "force", false, "force the deletion of the router without asking for confirmation")

	return cmd
}

func completeRouter(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	routers, err := compute.RouterService().List(ctx, core.CursorAll)
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

func findRouter(ctx context.Context, term string) (compute.Router, error) {
	routers, err := compute.RouterService().List(ctx, core.CursorAll)
	if err != nil {
		return compute.Router{}, fmt.Errorf("fetch routers: %w", err)
	}

	router, err := filter.FindOne(routers, term)
	if err != nil {
		return compute.Router{}, fmt.Errorf("find router: %w", err)
	}

	return router, nil
}
