package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func DeviceActionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "action",
		Short:   "Manage mac bare metal device actions",
		Example: "", // TODO
	}

	commands.Add(cmd, &deviceActionListCommand{}, &deviceActionRunCommand{})

	return cmd
}

type deviceActionListCommand struct {
}

func (d *deviceActionListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	device, err := findDevice(ctx, config, args[0])
	if err != nil {
		return err
	}

	availableActions := make([]macbaremetal.DeviceAction, len(device.Status.Actions))
	for i, action := range device.Status.Actions {
		availableActions[i] = macbaremetal.DeviceAction(action)
	}

	return commands.PrintStdout(availableActions)
}

func (d *deviceActionListCommand) Desc() *cobra.Command {
	return &cobra.Command{
		Use:     "list DEVICE",
		Short:   "List actions of device",
		Long:    "Lists all available actions of the specified device.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
	}
}

type deviceActionRunCommand struct {
}

func (d *deviceActionRunCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	return runAction(ctx, config, args[0], args[1])
}

func (d *deviceActionRunCommand) Desc() *cobra.Command {
	return &cobra.Command{
		Use:     "run DEVICE ACTION",
		Short:   "Run action on device",
		Long:    "Runs the specified action on the specified device.",
		Args:    cobra.ExactArgs(2),
		Example: "", // TODO
	}
}

type deviceActionRunCommandPreset string

func (d deviceActionRunCommandPreset) Run(ctx context.Context, config commands.Config, args []string) error {
	return runAction(ctx, config, args[0], string(d))
}

func (d deviceActionRunCommandPreset) Desc() *cobra.Command {
	return &cobra.Command{
		Use:   string(d) + " DEVICE",
		Short: "Run " + string(d) + " action on device",
		Long:  "Runs the " + string(d) + " action on the specified device.",
		Args:  cobra.ExactArgs(1),
	}
}

func DeviceWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workflow",
		Short:   "Manage mac bare metal device workflows",
		Example: "", // TODO
	}

	commands.Add(cmd, &deviceWorkflowListCommand{}, &deviceWorkflowRunCommand{})

	return cmd
}

type deviceWorkflowListCommand struct {
}

func (d *deviceWorkflowListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	device, err := findDevice(ctx, config, args[0])
	if err != nil {
		return err
	}

	service := macbaremetal.NewDeviceWorkflowService(config.Client, device.ID)

	workflows, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch workflows: %w", err)
	}

	return commands.PrintStdout(workflows)
}

func (d *deviceWorkflowListCommand) Desc() *cobra.Command {
	return &cobra.Command{
		Use:     "list DEVICE",
		Short:   "List device workflows",
		Long:    "Lists the available workflows on the specified device.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
	}
}

type deviceWorkflowRunCommand struct {
}

func (d *deviceWorkflowRunCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	device, err := findDevice(ctx, config, args[0])
	if err != nil {
		return err
	}

	service := macbaremetal.NewDeviceWorkflowService(config.Client, device.ID)

	workflows, err := service.List(ctx)
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

	device, err = service.Run(ctx, body)
	if err != nil {
		return fmt.Errorf("run workflow: %w", err)
	}

	return commands.PrintStdout(device)
}

func (d *deviceWorkflowRunCommand) Desc() *cobra.Command {
	return &cobra.Command{
		Use:     "run DEVICE WORKFLOW",
		Short:   "Run workflow on device",
		Long:    "Runs the specified workflow on the specified device.",
		Args:    cobra.ExactArgs(2),
		Example: "", // TODO
	}
}

func runAction(ctx context.Context, config commands.Config, deviceTerm, actionTerm string) error {
	device, err := findDevice(ctx, config, deviceTerm)
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

	device, err = macbaremetal.NewDeviceActionService(config.Client, device.ID).Run(ctx, body)
	if err != nil {
		return fmt.Errorf("run action: %w", err)
	}

	return commands.PrintStdout(device)
}
