package compute

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/objectstorage"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func InstanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "instance",
		Aliases: []string{"instances"},
		Short:   "Manage object storage instances",
	}

	commands.Add(cmd, &instanceListCommand{}, &instanceCreateCommand{}, &instanceDeleteCommand{}, &credentialsCommand{})

	return cmd
}

type instanceListCommand struct {
	filter string
}

func (i *instanceListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := objectstorage.NewInstanceService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch object storage instances: %w", err)
	}

	if len(i.filter) != 0 {
		items = filter.Find(items, i.filter)
	}

	return commands.PrintStdout(items)
}

func (i *instanceListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List instances",
		Long:    "Prints a table of all object storage instances belonging to the current organization.",
		RunE:    i.Run,
	}

	cmd.Flags().StringVar(&i.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type instanceCreateCommand struct {
	location string
}

func (i *instanceCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Config.Client, i.location)
	if err != nil {
		return err
	}

	data := objectstorage.InstanceCreate{
		LocationID: location.ID,
	}

	instance, err := objectstorage.NewInstanceService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create object storage instance: %w", err)
	}

	return commands.PrintStdout(instance)
}

func (i *instanceCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "new"},
		Short:   "Create new instance",
		Long:    "Creates a object storage instance.",
		Example: "", // TODO
		RunE:    i.Run,
	}

	cmd.Flags().StringVar(&i.location, "location", "", "location to be used for the instance")
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

type instanceDeleteCommand struct {
	force bool
}

func (i *instanceDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	service := objectstorage.NewInstanceService(commands.Config.Client)

	instances, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch object storage instances: %w", err)
	}

	instance, err := filter.FindOne(instances, args[0])
	if err != nil {
		return fmt.Errorf("find object storage instance: %w", err)
	}

	if !i.force && !commands.ConfirmDeletion("object storage instance", instance) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = service.Delete(cmd.Context(), instance.ID)
	if err != nil {
		return fmt.Errorf("delete instance: %w", err)
	}

	return nil
}

func (i *instanceDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete INSTANCE",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete instance",
		Long:    "Deletes an object storage instance.",
		Args:    cobra.ExactArgs(1),
		RunE:    i.Run,
	}

	cmd.Flags().BoolVar(&i.force, "force", false, "force the deletion of the instance without asking for confirmation")

	return cmd
}

type credentialsCommand struct {
	filter string
}

func (c *credentialsCommand) Run(cmd *cobra.Command, args []string) error {
	credentials, err := objectstorage.NewCredentialService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return err
	}

	if len(c.filter) != 0 {
		credentials = filter.Find(credentials, c.filter)
	}

	return commands.PrintStdout(credentials)
}

func (c *credentialsCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Show object storage credentials",
		Long:  "Prints a table of all object storage credentials belonging to the current organization.",
		RunE:  c.Run,
	}

	cmd.Flags().StringVar(&c.filter, "filter", "", "custom term to filter the results")

	return cmd
}
