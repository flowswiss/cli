package common

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/internal/commands"
	"github.com/flowswiss/cli/pkg/api/common"
	"github.com/flowswiss/cli/pkg/filter"
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

func (m *moduleListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	items, err := common.Modules(ctx, config.Client)
	if err != nil {
		return err
	}

	if len(m.filter) != 0 {
		items = filter.Find(items, m.filter)
	}

	return commands.PrintStdout(items)
}

func (m *moduleListCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available modules",
		Long:  "Lists all available modules including their location availability.",
	}

	cmd.Flags().StringVar(&m.filter, "filter", "", "custom term to filter the results")

	return cmd
}
