package common

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func LocationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "Manage datacenter locations",
	}

	commands.Add(cmd, &locationListCommand{})

	return cmd
}

type locationListCommand struct {
	filter string
}

func (l *locationListCommand) Run(cmd *cobra.Command, args []string) (err error) {
	items, err := common.Locations(cmd.Context(), commands.Config.Client)
	if err != nil {
		return err
	}

	if len(l.filter) != 0 {
		items = filter.Find(items, l.filter)
	}

	return commands.PrintStdout(items)
}

func (l *locationListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List datacenter locations",
		Long:  "Lists all datacenter locations including their available modules.",
		RunE:  l.Run,
	}

	cmd.Flags().StringVar(&l.filter, "filter", "", "custom term to filter the results")

	return cmd
}
