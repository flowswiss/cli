package macbaremetal

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func DeviceActionCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "action",
		Short: "Manage mac bare metal device actions",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # List all available actions for a specific device
      %[1]s mac-bare-metal device action list my-device
      
      # Run an action on a device
      %[1]s mac-bare-metal device action run my-device power-off
		`, app.Name)),
	}

	commands.Add(app, cmd,
		&deviceActionListCommand{},
		&deviceActionRunCommand{},
	)

	return cmd
}

type deviceActionListCommand struct {
}

func (d *deviceActionListCommand) Run(cmd *cobra.Command, args []string) error {
	device, err := findDevice(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	availableActions := make([]macbaremetal.DeviceAction, len(device.Status.Actions))
	for i, action := range device.Status.Actions {
		availableActions[i] = macbaremetal.DeviceAction(action)
	}

	return commands.PrintStdout(availableActions)
}

func (d *deviceActionListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeDevice(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d *deviceActionListCommand) Build(app commands.Application) *cobra.Command {
	return &cobra.Command{
		Use:               "list DEVICE",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List actions of device",
		Long:              "Lists all available actions of the specified device.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}
}

type deviceActionRunCommand struct {
}

func (d *deviceActionRunCommand) Run(cmd *cobra.Command, args []string) error {
	return runAction(cmd.Context(), args[0], args[1])
}

func (d *deviceActionRunCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeDevice(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		device, err := findDevice(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeAction(cmd.Context(), device, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d *deviceActionRunCommand) Build(app commands.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "run DEVICE ACTION",
		Short: "Run action on device",
		Long: commands.FormatHelp(fmt.Sprintf(`
			Runs the specified action on the specified device.

			To get a list of all available actions for a specific device, run "%[1]s mac-bare-metal device action list DEVICE".
		`, app.Name)),
		Example: commands.FormatExamples(fmt.Sprintf(`
      # Shutdown a device
      %[1]s device action run my-device power-off
      
      # Use the predefined action aliases
      %[1]s device power-off my-device
      %[1]s device power-on my-device
		`, app.Name)),
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}
}

type deviceActionRunCommandPreset string

func (d deviceActionRunCommandPreset) Run(cmd *cobra.Command, args []string) error {
	return runAction(cmd.Context(), args[0], string(d))
}

func (d *deviceActionRunCommandPreset) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeDevice(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d deviceActionRunCommandPreset) Build(app commands.Application) *cobra.Command {
	return &cobra.Command{
		Use:   string(d) + " DEVICE",
		Short: "Run " + string(d) + " action on device",
		Long: commands.FormatHelp(fmt.Sprintf(`
			Runs the %[2]s action on the specified device.

			This is a shortcut for "%[1]s mac-bare-metal device action run DEVICE %[2]s".
		`, app.Name, string(d))),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}
}

func DeviceWorkflowCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage mac bare metal device workflows",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # List all available workflows for a specific device
      %[1]s mac-bare-metal device workflow list my-device
      
      # Run the create snaphot workflow on a device
      %[1]s mac-bare-metal device workflow run my-device create_snapshot
		`, app.Name)),
	}

	commands.Add(app, cmd,
		&deviceWorkflowListCommand{},
		&deviceWorkflowRunCommand{},
	)

	return cmd
}

type deviceWorkflowListCommand struct {
}

func (d *deviceWorkflowListCommand) Run(cmd *cobra.Command, args []string) error {
	device, err := findDevice(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	service := macbaremetal.NewDeviceWorkflowService(commands.Config.Client, device.ID)

	workflows, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch workflows: %w", err)
	}

	return commands.PrintStdout(workflows)
}

func (d *deviceWorkflowListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeDevice(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d *deviceWorkflowListCommand) Build(app commands.Application) *cobra.Command {
	return &cobra.Command{
		Use:               "list DEVICE",
		Short:             "List device workflows",
		Long:              "Lists the available workflows on the specified device.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}
}

type deviceWorkflowRunCommand struct {
}

func (d *deviceWorkflowRunCommand) Run(cmd *cobra.Command, args []string) error {
	device, err := findDevice(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	service := macbaremetal.NewDeviceWorkflowService(commands.Config.Client, device.ID)

	workflows, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch workflows: %w", err)
	}

	workflow, err := filter.FindOne(workflows, args[1])
	if err != nil {
		return fmt.Errorf("find workflow: %w", err)
	}

	body := macbaremetal.DeviceRunWorkflow{
		Workflow: workflow.Command,
	}

	device, err = service.Run(cmd.Context(), body)
	if err != nil {
		return fmt.Errorf("run workflow: %w", err)
	}

	return commands.PrintStdout(device)
}

func (d *deviceWorkflowRunCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeDevice(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		device, err := findDevice(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeWorkflow(cmd.Context(), device, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (d *deviceWorkflowRunCommand) Build(app commands.Application) *cobra.Command {
	return &cobra.Command{
		Use:               "run DEVICE WORKFLOW",
		Short:             "Run workflow on device",
		Long:              "Runs the specified workflow on the specified device.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: d.CompleteArg,
		RunE:              d.Run,
	}
}

func completeAction(ctx context.Context, device macbaremetal.Device, term string) ([]string, cobra.ShellCompDirective) {
	actions := make([]string, len(device.Status.Actions))
	for i, action := range device.Status.Actions {
		actions[i] = action.Command
	}

	return actions, cobra.ShellCompDirectiveNoFileComp
}

func completeWorkflow(ctx context.Context, device macbaremetal.Device, term string) ([]string, cobra.ShellCompDirective) {
	workflows, err := macbaremetal.NewDeviceWorkflowService(commands.Config.Client, device.ID).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(workflows, term)

	names := make([]string, len(filtered))
	for i, w := range filtered {
		names[i] = w.Command
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func runAction(ctx context.Context, deviceTerm, actionTerm string) error {
	device, err := findDevice(ctx, deviceTerm)
	if err != nil {
		return err
	}

	availableActions := make([]macbaremetal.DeviceAction, len(device.Status.Actions))
	for i, action := range device.Status.Actions {
		availableActions[i] = macbaremetal.DeviceAction(action)
	}

	action, err := filter.FindOne(availableActions, actionTerm)
	if err != nil {
		return fmt.Errorf("the selected action does not exist or is currently not possible")
	}

	body := macbaremetal.DeviceRunAction{
		Action: action.Command,
	}

	device, err = macbaremetal.NewDeviceActionService(commands.Config.Client, device.ID).Run(ctx, body)
	if err != nil {
		return fmt.Errorf("run action: %w", err)
	}

	return commands.PrintStdout(device)
}
