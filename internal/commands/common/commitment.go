package common

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/filter"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/spf13/cobra"
)

func Commitment(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commitment",
		Short: "Manage product commitments",
	}

	commands.Add(app, cmd,
		&commitmentListCommand{},
		&commitmentUpdateCommand{},
	)

	return cmd
}

type commitmentListCommand struct {
	filter string
}

func (c *commitmentListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := common.CommitmentService(commands.Client).List(cmd.Context(), core.CursorAll)
	if err != nil {
		return err
	}

	if len(c.filter) != 0 {
		items = filter.Find(items, c.filter)
	}

	return commands.PrintStdout(items)
}

func (c *commitmentListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *commitmentListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List commitments",
		Long:              "Lists all commitments across products.",
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().StringVar(&c.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type commitmentUpdateCommand struct {
	renew bool
}

func (c *commitmentUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := common.CommitmentService(commands.Client).List(cmd.Context(), core.CursorAll)
	if err != nil {
		return err
	}

	commitment, err := filter.FindOne(items, args[0])
	if err != nil {
		return fmt.Errorf("find commitment: %w", err)
	}

	data := common.CommitmentUpdate{
		ID:    uint(commitment.ID),
		Renew: &c.renew,
	}

	commitment, err = common.CommitmentService(commands.Client).Update(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("update commitment: %w", err)
	}

	return commands.PrintStdout(commitment)
}

func (c *commitmentUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *commitmentUpdateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update COMMITMENT",
		Short:             "Update commitment",
		Long:              "Updates a commitment.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().BoolVar(&c.renew, "renew", true, "if the commitment should be renewed on expiration")

	return cmd
}
