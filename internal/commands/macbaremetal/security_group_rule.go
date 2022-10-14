package macbaremetal

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/macbaremetal"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func SecurityGroupRuleCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rule",
		Aliases: []string{"rules"},
		Short:   "Manage mac bare metal security group rules",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # List all security group rules
      %[1]s mac-bare-metal security-group rule list default
      
      # Create new security group rule to allow tcp traffic on port 80 (HTTP) from any source IP
      %[1]s mac-bare-metal security-group rule create default --direction ingress --protocol tcp --from-port 80 --to-port 80
		`, app.Name)),
	}

	commands.Add(app, cmd,
		&securityGroupRuleListCommand{},
		&securityGroupRuleCreateCommand{},
		&securityGroupRuleUpdateCommand{},
		&securityGroupRuleDeleteCommand{},
	)

	return cmd
}

type securityGroupRuleListCommand struct {
	filter string
}

func (s *securityGroupRuleListCommand) Run(cmd *cobra.Command, args []string) error {
	securityGroup, err := findSecurityGroup(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	service := macbaremetal.NewSecurityGroupRuleService(commands.Config.Client, securityGroup.ID)

	items, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch security group rules: %w", err)
	}

	if len(s.filter) != 0 {
		items = filter.Find(items, s.filter)
	}

	return commands.PrintStdout(items)
}

func (s *securityGroupRuleListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeSecurityGroup(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *securityGroupRuleListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list SECURITY-GROUP",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List security group rules",
		Long:              "Lists all mac bare metal security group rules.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type securityGroupRuleCreateCommand struct {
	direction string
	protocol  string
	fromPort  int
	toPort    int
	icmpType  int
	icmpCode  int
	ipRange   net.IPNet
}

func (s *securityGroupRuleCreateCommand) Run(cmd *cobra.Command, args []string) error {
	securityGroup, err := findSecurityGroup(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	service := macbaremetal.NewSecurityGroupRuleService(commands.Config.Client, securityGroup.ID)

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

	item, err := service.Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create security group rule: %w", err)
	}

	return commands.PrintStdout(item)
}

func (s *securityGroupRuleCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeSecurityGroup(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *securityGroupRuleCreateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create SECURITY-GROUP",
		Aliases: []string{"add", "new"},
		Short:   "Create new security group",
		Long:    "Creates a new mac bare metal security group.",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # Create rule to allow tcp traffic on port 80 (HTTP) from any source IP
      %[1]s mac-bare-metal security-group rule create default --direction ingress --protocol tcp --from-port 80 --to-port 80
      
      # Create rule to allow tcp traffic on port 22 (SSH) only from subnet 1.1.1.0/24
      %[1]s mac-bare-metal security-group rule create default --direction ingress --protocol tcp --from-port 22 --to-port 22 --ip-range 1.1.1.0/24
		`, app.Name)),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.direction, "direction", "", "direction of the rule")
	cmd.Flags().StringVar(&s.protocol, "protocol", "", "protocol of the rule")
	cmd.Flags().IntVar(&s.fromPort, "from-port", 0, "from port of the rule (only for TCP and UDP)")
	cmd.Flags().IntVar(&s.toPort, "to-port", 0, "to port of the rule (only for TCP and UDP)")
	cmd.Flags().IntVar(&s.icmpType, "icmp-type", 0, "icmp type of the rule (only for ICMP)")
	cmd.Flags().IntVar(&s.icmpCode, "icmp-code", 0, "icmp code of the rule (only for ICMP)")
	cmd.Flags().IPNetVar(&s.ipRange, "ip-range", macbaremetal.IPRangeAny, "ip range of the rule")

	_ = cmd.MarkFlagRequired("direction")
	_ = cmd.MarkFlagRequired("protocol")

	return cmd
}

type securityGroupRuleUpdateCommand struct {
	direction string
	protocol  string
	fromPort  int
	toPort    int
	icmpType  int
	icmpCode  int
	ipRange   net.IPNet
}

func (s *securityGroupRuleUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	securityGroup, err := findSecurityGroup(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	service := macbaremetal.NewSecurityGroupRuleService(commands.Config.Client, securityGroup.ID)

	rules, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch security group rules: %w", err)
	}

	rule, err := filter.FindOne(rules, args[1])
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

	item, err := service.Update(cmd.Context(), rule.ID, data)
	if err != nil {
		return fmt.Errorf("create security group rule: %w", err)
	}

	return commands.PrintStdout(item)
}

func (s *securityGroupRuleUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeSecurityGroup(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		securityGroup, err := findSecurityGroup(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeSecurityGroupRule(cmd.Context(), securityGroup, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *securityGroupRuleUpdateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update SECURITY-GROUP RULE",
		Short: "Update security group rule",
		Long:  "Updates a mac bare metal security group rule.",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # Update SSH rule to allow tcp traffic from broader subnet 1.1.0.0/16
      %[1]s mac-bare-metal security-group rule update default 1234 --direction ingress --protocol tcp --from-port 22 --to-port 22 --ip-range 1.1.0.0/16
		`, app.Name)), // TODO
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.direction, "direction", "", "direction of the rule")
	cmd.Flags().StringVar(&s.protocol, "protocol", "", "protocol of the rule")
	cmd.Flags().IntVar(&s.fromPort, "from-port", 0, "from port of the rule (only for TCP and UDP)")
	cmd.Flags().IntVar(&s.toPort, "to-port", 0, "to port of the rule (only for TCP and UDP)")
	cmd.Flags().IntVar(&s.icmpType, "icmp-type", 0, "icmp type of the rule (only for ICMP)")
	cmd.Flags().IntVar(&s.icmpCode, "icmp-code", 0, "icmp code of the rule (only for ICMP)")
	cmd.Flags().IPNetVar(&s.ipRange, "ip-range", macbaremetal.IPRangeAny, "ip range of the rule")

	_ = cmd.MarkFlagRequired("direction")
	_ = cmd.MarkFlagRequired("protocol")

	return cmd
}

type securityGroupRuleDeleteCommand struct {
	force bool
}

func (s *securityGroupRuleDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	securityGroup, err := findSecurityGroup(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	service := macbaremetal.NewSecurityGroupRuleService(commands.Config.Client, securityGroup.ID)

	rules, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch security group rules: %w", err)
	}

	rule, err := filter.FindOne(rules, args[1])
	if err != nil {
		return fmt.Errorf("find security group rule: %w", err)
	}

	if !s.force && !commands.ConfirmDeletion("security group", securityGroup) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = service.Delete(cmd.Context(), rule.ID)
	if err != nil {
		return fmt.Errorf("delete security group rule: %w", err)
	}

	return nil
}

func (s *securityGroupRuleDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeSecurityGroup(cmd.Context(), toComplete)
	}

	if len(args) == 1 {
		securityGroup, err := findSecurityGroup(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return completeSecurityGroupRule(cmd.Context(), securityGroup, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *securityGroupRuleDeleteCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete SECURITY-GROUP RULE",
		Aliases:           []string{"del", "remove", "rm"},
		Short:             "Delete security group rule",
		Long:              "Deletes a mac bare metal security group rule.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	return cmd
}

func completeSecurityGroupRule(ctx context.Context, securityGroup macbaremetal.SecurityGroup, term string) ([]string, cobra.ShellCompDirective) {
	rules, err := macbaremetal.NewSecurityGroupRuleService(commands.Config.Client, securityGroup.ID).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(rules, term)

	names := make([]string, len(filtered))
	for i, rule := range filtered {
		names[i] = fmt.Sprint(rule.ID)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
