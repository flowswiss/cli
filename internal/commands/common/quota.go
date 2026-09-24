package common

import (
	"strings"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/filter"
	"github.com/spf13/cobra"
)

func Quota(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota",
		Short: "Manage quota",
	}

	commands.Add(app, cmd,
		&quotaGetCommand{},
	)

	return cmd
}

type quotaGetCommand struct {
	scope string

	filter string
}

func (q *quotaGetCommand) Run(cmd *cobra.Command, args []string) error {
	quotaOverview, err := common.QuotaService(commands.Client).Get(cmd.Context(), common.QuotaGet{
		Scope: q.scope,
	})
	if err != nil {
		return err
	}

	quotas := quotaOverview.Quotas()

	if len(q.filter) != 0 {
		quotas = filter.Find(quotas, q.filter)
	}

	if quotaOverview.Scope.Current != "" {
		commands.Stdout.Bold().Print("Current scope: ").Reset()
		commands.Stdout.Println(quotaOverview.Scope.Current)
		commands.Stdout.Println()
	}
	if len(quotaOverview.Scope.Children) != 0 {
		commands.Stdout.Bold().Print("Available scope(s): ").Reset()
		commands.Stdout.Println(quotaOverview.Scope.Children[0])

		for i := 1; i < len(quotaOverview.Scope.Children); i++ {
			commands.Stdout.Println(strings.Repeat(" ", 19), quotaOverview.Scope.Children[i])
		}
		commands.Stdout.Println()
	}

	return commands.PrintStdout(quotas)
}

func (q *quotaGetCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (q *quotaGetCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List quotas",
		Long:              "Lists all quotas across products.",
		ValidArgsFunction: q.CompleteArg,
		RunE:              q.Run,
	}

	cmd.Flags().StringVar(&q.scope, "scope", "", "scope of the quotas to list")
	cmd.Flags().StringVar(&q.filter, "filter", "", "custom term to filter the results")

	return cmd
}
