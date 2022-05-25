package compute

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/internal/commands"
	"github.com/flowswiss/cli/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/pkg/filter"
)

func SecurityGroupRuleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rule",
		Short:   "Manage mac bare metal security group rules",
		Example: "", // TODO
	}

	commands.Add(cmd, &securityGroupRuleListCommand{}, &securityGroupRuleCreateCommand{}, &securityGroupUpdateCommand{}, &securityGroupRuleDeleteCommand{})

	return cmd
}

type securityGroupRuleListCommand struct {
	securityGroup string
	filter        string
}

func (s *securityGroupRuleListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	securityGroup, err := findSecurityGroup(ctx, config, s.securityGroup)
	if err != nil {
		return err
	}

	service := macbaremetal.NewSecurityGroupRuleService(config.Client, securityGroup.ID)

	items, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch security group rules: %w", err)
	}

	if len(s.filter) != 0 {
		items = filter.Find(items, s.filter)
	}

	return commands.PrintStdout(items)
}

func (s *securityGroupRuleListCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List security group rules",
		Long:    "Lists all mac bare metal security group rules.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&s.securityGroup, "security-group", "", "security group to list rules for")
	cmd.Flags().StringVar(&s.filter, "filter", "", "custom term to filter the results")

	_ = cmd.MarkFlagRequired("security-group")

	return cmd
}

type securityGroupRuleCreateCommand struct {
	securityGroup string
	direction     string
	protocol      string
	fromPort      int
	toPort        int
	icmpType      int
	icmpCode      int
	ipRange       net.IPNet
}

func (s *securityGroupRuleCreateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	securityGroup, err := findSecurityGroup(ctx, config, s.securityGroup)
	if err != nil {
		return err
	}

	service := macbaremetal.NewSecurityGroupRuleService(config.Client, securityGroup.ID)

	protocol, found := macbaremetal.ProtocolIDs[strings.ToLower(s.protocol)]
	if !found {
		return fmt.Errorf("invalid protocol: %s", s.protocol)
	}

	data := macbaremetal.SecurityGroupRuleCreate{
		Direction: s.direction,
		Protocol:  protocol,
		FromPort:  s.fromPort,
		ToPort:    s.toPort,
		ICMPType:  s.icmpType,
		ICMPCode:  s.icmpCode,
		IPRange:   s.ipRange.String(),
	}

	item, err := service.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("create security group rule: %w", err)
	}

	return commands.PrintStdout(item)
}

func (s *securityGroupRuleCreateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create new security group",
		Long:    "Creates a new mac bare metal security group.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&s.securityGroup, "security-group", "", "security group to create rule for")
	cmd.Flags().StringVar(&s.direction, "direction", "", "direction of the rule")
	cmd.Flags().StringVar(&s.protocol, "protocol", "", "protocol of the rule")
	cmd.Flags().IntVar(&s.fromPort, "from-port", 0, "from port of the rule (only for TCP and UDP)")
	cmd.Flags().IntVar(&s.toPort, "to-port", 0, "to port of the rule (only for TCP and UDP)")
	cmd.Flags().IntVar(&s.icmpType, "icmp-type", 0, "icmp type of the rule (only for ICMP)")
	cmd.Flags().IntVar(&s.icmpCode, "icmp-code", 0, "icmp code of the rule (only for ICMP)")
	cmd.Flags().IPNetVar(&s.ipRange, "ip-range", macbaremetal.IPRangeAny, "ip range of the rule")

	_ = cmd.MarkFlagRequired("security-group")
	_ = cmd.MarkFlagRequired("direction")
	_ = cmd.MarkFlagRequired("protocol")

	return cmd
}

type securityGroupRuleUpdateCommand struct {
	securityGroup string
	rule          string
	direction     string
	protocol      string
	fromPort      int
	toPort        int
	icmpType      int
	icmpCode      int
	ipRange       net.IPNet
}

func (s *securityGroupRuleUpdateCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	securityGroup, err := findSecurityGroup(ctx, config, s.securityGroup)
	if err != nil {
		return err
	}

	service := macbaremetal.NewSecurityGroupRuleService(config.Client, securityGroup.ID)

	rules, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch security group rules: %w", err)
	}

	rule, err := filter.FindOne(rules, s.rule)
	if err != nil {
		return fmt.Errorf("find security group rule: %w", err)
	}

	protocol, found := macbaremetal.ProtocolIDs[strings.ToLower(s.protocol)]
	if !found {
		return fmt.Errorf("invalid protocol: %s", s.protocol)
	}

	data := macbaremetal.SecurityGroupRuleCreate{
		Direction: s.direction,
		Protocol:  protocol,
		FromPort:  s.fromPort,
		ToPort:    s.toPort,
		ICMPType:  s.icmpType,
		ICMPCode:  s.icmpCode,
		IPRange:   s.ipRange.String(),
	}

	item, err := service.Update(ctx, rule.ID, data)
	if err != nil {
		return fmt.Errorf("create security group rule: %w", err)
	}

	return commands.PrintStdout(item)
}

func (s *securityGroupRuleUpdateCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update security group rule",
		Long:    "Updates a mac bare metal security group rule.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&s.securityGroup, "security-group", "", "security group to update rule for")
	cmd.Flags().StringVar(&s.rule, "rule", "", "rule to update")
	cmd.Flags().StringVar(&s.direction, "direction", "", "direction of the rule")
	cmd.Flags().StringVar(&s.protocol, "protocol", "", "protocol of the rule")
	cmd.Flags().IntVar(&s.fromPort, "from-port", 0, "from port of the rule (only for TCP and UDP)")
	cmd.Flags().IntVar(&s.toPort, "to-port", 0, "to port of the rule (only for TCP and UDP)")
	cmd.Flags().IntVar(&s.icmpType, "icmp-type", 0, "icmp type of the rule (only for ICMP)")
	cmd.Flags().IntVar(&s.icmpCode, "icmp-code", 0, "icmp code of the rule (only for ICMP)")
	cmd.Flags().IPNetVar(&s.ipRange, "ip-range", macbaremetal.IPRangeAny, "ip range of the rule")

	_ = cmd.MarkFlagRequired("security-group")
	_ = cmd.MarkFlagRequired("direction")
	_ = cmd.MarkFlagRequired("protocol")

	return cmd
}

type securityGroupRuleDeleteCommand struct {
	securityGroup     string
	securityGroupRule string
	force             bool
}

func (s *securityGroupRuleDeleteCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	securityGroup, err := findSecurityGroup(ctx, config, s.securityGroup)
	if err != nil {
		return err
	}

	service := macbaremetal.NewSecurityGroupRuleService(config.Client, securityGroup.ID)

	rules, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("fetch security group rules: %w", err)
	}

	rule, err := filter.FindOne(rules, s.securityGroupRule)
	if err != nil {
		return fmt.Errorf("find security group rule: %w", err)
	}

	// TODO ask for confirmation

	err = service.Delete(ctx, rule.ID)
	if err != nil {
		return fmt.Errorf("delete security group rule: %w", err)
	}

	return nil
}

func (s *securityGroupRuleDeleteCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete security group rule",
		Long:    "Deletes a mac bare metal security group rule.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&s.securityGroup, "security-group", "", "security group of the rule to be deleted")
	cmd.Flags().StringVar(&s.securityGroupRule, "rule", "", "security group rule to be deleted")

	_ = cmd.MarkFlagRequired("security-group")
	_ = cmd.MarkFlagRequired("rule")

	return cmd
}

func findSecurityGroup(ctx context.Context, config commands.Config, term string) (macbaremetal.SecurityGroup, error) {
	securityGroups, err := macbaremetal.NewSecurityGroupService(config.Client).List(ctx)
	if err != nil {
		return macbaremetal.SecurityGroup{}, fmt.Errorf("fetch security groups: %w", err)
	}

	securityGroup, err := filter.FindOne(securityGroups, term)
	if err != nil {
		return macbaremetal.SecurityGroup{}, fmt.Errorf("find security group: %w", err)
	}

	return securityGroup, nil
}
