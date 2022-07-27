package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func SecurityGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security-group",
		Aliases: []string{"security-groups", "securitygroup", "securitygroups"},
		Short:   "Manage compute security groups",
	}

	commands.Add(cmd, &securityGroupListCommand{}, &securityGroupCreateCommand{}, &securityGroupUpdateCommand{}, &securityGroupDeleteCommand{})
	cmd.AddCommand(SecurityGroupRuleCommand())

	return cmd
}

type securityGroupListCommand struct {
	filter string
}

func (s *securityGroupListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.NewSecurityGroupService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch security groups: %w", err)
	}

	if len(s.filter) != 0 {
		items = filter.Find(items, s.filter)
	}

	return commands.PrintStdout(items)
}

func (s *securityGroupListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *securityGroupListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List security groups",
		Long:              "Lists all compute security groups.",
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type securityGroupCreateCommand struct {
	name        string
	description string
	location    string
}

func (s *securityGroupCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Config.Client, s.location)
	if err != nil {
		return err
	}

	data := compute.SecurityGroupCreate{
		Name:        s.name,
		Description: s.description,
		LocationID:  location.ID,
	}

	item, err := compute.NewSecurityGroupService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create security group: %w", err)
	}

	return commands.PrintStdout(item)
}

func (s *securityGroupCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *securityGroupCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Aliases:           []string{"add", "new"},
		Short:             "Create new security group",
		Long:              "Creates a new compute security group.",
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.name, "name", "", "name to be applied to the security group")
	cmd.Flags().StringVar(&s.description, "description", "", "description to be applied to the security group")
	cmd.Flags().StringVar(&s.location, "location", "", "location where the security group will be created")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

type securityGroupUpdateCommand struct {
	name        string
	description string
}

func (s *securityGroupUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	service := compute.NewSecurityGroupService(commands.Config.Client)

	securityGroups, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch security groups: %w", err)
	}

	securityGroup, err := filter.FindOne(securityGroups, args[0])
	if err != nil {
		return fmt.Errorf("find security group: %w", err)
	}

	update := compute.SecurityGroupUpdate{
		Name:        s.name,
		Description: s.description,
	}

	securityGroup, err = service.Update(cmd.Context(), securityGroup.ID, update)
	if err != nil {
		return fmt.Errorf("update security group: %w", err)
	}

	return commands.PrintStdout(securityGroup)
}

func (s *securityGroupUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeSecurityGroup(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *securityGroupUpdateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update SECURITY-GROUP",
		Short:             "Update security group",
		Long:              "Updates a compute security group.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.name, "name", "", "name to be applied to the security group")
	cmd.Flags().StringVar(&s.description, "description", "", "description to be applied to the security group")

	return cmd
}

type securityGroupDeleteCommand struct {
	force bool
}

func (s *securityGroupDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	service := compute.NewSecurityGroupService(commands.Config.Client)

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

func (s *securityGroupDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeSecurityGroup(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *securityGroupDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete SECURITY-GROUP",
		Aliases:           []string{"del", "remove", "rm"},
		Short:             "Delete security group",
		Long:              "Deletes a compute security group.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().BoolVar(&s.force, "force", false, "force the deletion of the security group without asking for confirmation")

	return cmd
}

func completeSecurityGroup(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	securityGroups, err := compute.NewSecurityGroupService(commands.Config.Client).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(securityGroups, term)

	names := make([]string, len(filtered))
	for i, securityGroup := range filtered {
		names[i] = securityGroup.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findSecurityGroup(ctx context.Context, term string) (compute.SecurityGroup, error) {
	securityGroups, err := compute.NewSecurityGroupService(commands.Config.Client).List(ctx)
	if err != nil {
		return compute.SecurityGroup{}, fmt.Errorf("fetch security groups: %w", err)
	}

	securityGroup, err := filter.FindOne(securityGroups, term)
	if err != nil {
		return compute.SecurityGroup{}, fmt.Errorf("find security group: %w", err)
	}

	return securityGroup, nil
}
