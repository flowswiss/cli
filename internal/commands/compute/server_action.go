package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ServerActionCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "action",
		Aliases: []string{"actions"},
		Short:   "Manage compute server actions",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # List all available actions for a specific server
      %[1]s compute server action list my-server
      
      # Run an action on a server
      %[1]s compute server action run my-server stop
		`, app.Name)),
	}

	commands.Add(app, cmd, &serverActionListCommand{}, &serverActionRunCommand{})

	return cmd
}

type serverActionListCommand struct {
}

func (s *serverActionListCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	availableActions := make([]compute.ServerAction, len(server.Status.Actions))
	for i, action := range server.Status.Actions {
		availableActions[i] = compute.ServerAction(action)
	}

	return commands.PrintStdout(availableActions)
}

func (s *serverActionListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverActionListCommand) Build(app commands.Application) *cobra.Command {
	return &cobra.Command{
		Use:               "list SERVER",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List actions of server",
		Long:              "Lists all available actions of the specified server.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}
}

type serverActionRunCommand struct {
}

func (s *serverActionRunCommand) Run(cmd *cobra.Command, args []string) error {
	return runAction(cmd.Context(), args[0], args[1])
}

func (s *serverActionRunCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		server, err := findServer(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeServerAction(cmd.Context(), server, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverActionRunCommand) Build(app commands.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "run SERVER ACTION",
		Short: "Run action on server",
		Long: commands.FormatHelp(fmt.Sprintf(`
			Runs the specified action on the specified server.

			To get a list of all available actions for a specific server, run "%[1]s compute server action list SERVER".
		`, app.Name)),
		Example: commands.FormatExamples(fmt.Sprintf(`
      # Shutdown a server
      %[1]s server action run my-server stop
      
      # Use the predefined action aliases
      %[1]s server stop my-server
      %[1]s server start my-server
		`, app.Name)), // TODO
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}
}

type serverActionRunCommandPreset string

func (s serverActionRunCommandPreset) Run(cmd *cobra.Command, args []string) error {
	return runAction(cmd.Context(), args[0], string(s))
}

func (s *serverActionRunCommandPreset) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s serverActionRunCommandPreset) Build(app commands.Application) *cobra.Command {
	return &cobra.Command{
		Use:   string(s) + " SERVER",
		Short: "Run " + string(s) + " action on the server",
		Long: commands.FormatHelp(fmt.Sprintf(`
			Runs the %[2]s action on the specified server.

			This is a shortcut for "%[1]s compute server action run SERVER %[2]s".
		`, app.Name, string(s))),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}
}

func completeServerAction(ctx context.Context, server compute.Server, term string) ([]string, cobra.ShellCompDirective) {
	actions := make([]compute.ServerAction, len(server.Status.Actions))
	for i, action := range server.Status.Actions {
		actions[i] = compute.ServerAction(action)
	}

	filtered := filter.Find(actions, term)

	names := make([]string, len(filtered))
	for i, action := range filtered {
		names[i] = action.Command
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func runAction(ctx context.Context, serverTerm, actionTerm string) error {
	server, err := findServer(ctx, serverTerm)
	if err != nil {
		return err
	}

	availableActions := make([]compute.ServerAction, len(server.Status.Actions))
	for i, action := range server.Status.Actions {
		availableActions[i] = compute.ServerAction(action)
	}

	action, err := filter.FindOne(availableActions, actionTerm)
	if err != nil {
		return fmt.Errorf("the selected action does not exist or is currently not possible")
	}

	body := compute.ServerRunAction{
		Action: action.Command,
	}

	server, err = compute.NewServerActionService(commands.Config.Client).Run(ctx, server.ID, body)
	if err != nil {
		return fmt.Errorf("run action: %w", err)
	}

	return commands.PrintStdout(server)
}
