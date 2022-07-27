package common

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ModuleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module",
		Short: "Manage modules",
	}

	commands.Add(cmd, &moduleListCommand{})

	return cmd
}

type moduleListCommand struct {
	filter string
}

func (m *moduleListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := common.Modules(cmd.Context(), commands.Config.Client)
	if err != nil {
		return err
	}

	if len(m.filter) != 0 {
		items = filter.Find(items, m.filter)
	}

	return commands.PrintStdout(items)
}

func (m *moduleListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (m *moduleListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Short:             "List available modules",
		Long:              "Lists all available modules including their location availability.",
		ValidArgsFunction: m.CompleteArg,
		RunE:              m.Run,
	}

	cmd.Flags().StringVar(&m.filter, "filter", "", "custom term to filter the results")

	return cmd
}
