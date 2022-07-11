package compute

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func SecurityGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security-group",
		Aliases: []string{"security-groups", "securitygroup", "securitygroups"},
		Short:   "Manage mac bare metal security groups",
	}

	commands.Add(cmd, &securityGroupListCommand{}, &securityGroupCreateCommand{}, &securityGroupUpdateCommand{}, &securityGroupDeleteCommand{})
	cmd.AddCommand(SecurityGroupRuleCommand())

	return cmd
}

type securityGroupListCommand struct {
	filter string
}

func (s *securityGroupListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := macbaremetal.NewSecurityGroupService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch security groups: %w", err)
	}

	if len(s.filter) != 0 {
		items = filter.Find(items, s.filter)
	}

	return commands.PrintStdout(items)
}

func (s *securityGroupListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List security groups",
		Long:    "Lists all mac bare metal security groups.",
		RunE:    s.Run,
	}

	cmd.Flags().StringVar(&s.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type securityGroupCreateCommand struct {
	name        string
	description string
	network     string
}

func (s *securityGroupCreateCommand) Run(cmd *cobra.Command, args []string) error {
	networks, err := macbaremetal.NewNetworkService(commands.Config.Client).List(cmd.Context())
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

	item, err := macbaremetal.NewSecurityGroupService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create security group: %w", err)
	}

	return commands.PrintStdout(item)
}

func (s *securityGroupCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "new"},
		Short:   "Create new security group",
		Long:    "Creates a new mac bare metal security group.",
		RunE:    s.Run,
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

func (s *securityGroupUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	service := macbaremetal.NewSecurityGroupService(commands.Config.Client)

	securityGroups, err := service.List(cmd.Context())
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

	securityGroup, err = service.Update(cmd.Context(), securityGroup.ID, update)
	if err != nil {
		return fmt.Errorf("update security group: %w", err)
	}

	return commands.PrintStdout(securityGroup)
}

func (s *securityGroupUpdateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update SECURITY-GROUP",
		Short: "Update security group",
		Long:  "Updates a mac bare metal security group.",
		Args:  cobra.ExactArgs(1),
		RunE:  s.Run,
	}

	cmd.Flags().StringVar(&s.name, "name", "", "name to be applied to the security group")
	cmd.Flags().StringVar(&s.description, "description", "", "description to be applied to the security group")

	return cmd
}

type securityGroupDeleteCommand struct {
	force bool
}

func (s *securityGroupDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	service := macbaremetal.NewSecurityGroupService(commands.Config.Client)

	securityGroups, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch security groups: %w", err)
	}

	securityGroup, err := filter.FindOne(securityGroups, args[0])
	if err != nil {
		return fmt.Errorf("find security group: %w", err)
	}

	if !s.force && !commands.ConfirmDeletion("security group", securityGroup) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = service.Delete(cmd.Context(), securityGroup.ID)
	if err != nil {
		return fmt.Errorf("delete security group: %w", err)
	}

	return nil
}

func (s *securityGroupDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete SECURITY-GROUP",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete security group",
		Long:    "Deletes a mac bare metal security group.",
		Args:    cobra.ExactArgs(1),
		RunE:    s.Run,
	}

	cmd.Flags().BoolVar(&s.force, "force", false, "force the deletion of the security group without asking for confirmation")

	return cmd
}
