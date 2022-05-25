package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/internal/commands"
	"github.com/flowswiss/cli/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/pkg/filter"
)

func SecurityGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security-group",
		Short:   "Manage mac bare metal security groups",
		Example: "", // TODO
	}

	commands.Add(cmd, &securityGroupListCommand{}, &securityGroupCreateCommand{}, &securityGroupUpdateCommand{}, &securityGroupDeleteCommand{})
	cmd.AddCommand(SecurityGroupRuleCommand())

	return cmd
}

type securityGroupListCommand struct {
	filter string
}

func (s *securityGroupListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	items, err := macbaremetal.NewSecurityGroupService(config.Client).List(ctx)
	if err != nil {
		return fmt.Errorf("fetch security groups: %w", err)
	}

	if len(s.filter) != 0 {
		items = filter.Find(items, s.filter)
	}

	return commands.PrintStdout(items)
}

func (s *securityGroupListCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List security groups",
		Long:    "Lists all mac bare metal security groups.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&s.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type securityGroupCreateCommand struct {
	name        string
	description string
	network     string
}

func (s *securityGroupCreateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	networks, err := macbaremetal.NewNetworkService(config.Client).List(ctx)
	if err != nil {
		return fmt.Errorf("fetch networks: %w", err)
	}

	network, err := filter.FindOne(networks, s.network)
	if err != nil {
		return fmt.Errorf("find network: %w", err)
	}

	data := macbaremetal.SecurityGroupCreate{
		Name:        s.name,
		Description: s.description,
		NetworkID:   network.ID,
	}

	item, err := macbaremetal.NewSecurityGroupService(config.Client).Create(ctx, data)
	if err != nil {
		return fmt.Errorf("create security group: %w", err)
	}

	return commands.PrintStdout(item)
}

func (s *securityGroupCreateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create new security group",
		Long:    "Creates a new mac bare metal security group.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&s.name, "name", "", "name to be applied to the security group")
	cmd.Flags().StringVar(&s.description, "description", "", "description to be applied to the security group")
	cmd.Flags().StringVar(&s.network, "network", "", "network in which the security group will be created")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("network")

	return cmd
}

type securityGroupUpdateCommand struct {
	name        string
	description string
}

func (s *securityGroupUpdateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	service := macbaremetal.NewSecurityGroupService(config.Client)

	securityGroups, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch security groups: %w", err)
	}

	securityGroup, err := filter.FindOne(securityGroups, args[0])
	if err != nil {
		return fmt.Errorf("find security group: %w", err)
	}

	update := macbaremetal.SecurityGroupUpdate{
		Name:        s.name,
		Description: s.description,
	}

	securityGroup, err = service.Update(ctx, securityGroup.ID, update)
	if err != nil {
		return fmt.Errorf("update security group: %w", err)
	}

	return commands.PrintStdout(securityGroup)
}

func (s *securityGroupUpdateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update SECURITY-GROUP",
		Short:   "Update security group",
		Long:    "Updates a mac bare metal security group.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&s.name, "name", "", "name to be applied to the security group")
	cmd.Flags().StringVar(&s.description, "description", "", "description to be applied to the security group")

	return cmd
}

type securityGroupDeleteCommand struct {
	force bool
}

func (s *securityGroupDeleteCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	service := macbaremetal.NewSecurityGroupService(config.Client)

	securityGroups, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch security groups: %w", err)
	}

	securityGroup, err := filter.FindOne(securityGroups, args[0])
	if err != nil {
		return fmt.Errorf("find security group: %w", err)
	}

	// TODO ask for confirmation

	err = service.Delete(ctx, securityGroup.ID)
	if err != nil {
		return fmt.Errorf("delete security group: %w", err)
	}

	return nil
}

func (s *securityGroupDeleteCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete SECURITY-GROUP",
		Short:   "Delete security group",
		Long:    "Deletes a mac bare metal security group.",
		Args:    cobra.ExactArgs(1),
		Example: "", // TODO
	}

	return cmd
}
