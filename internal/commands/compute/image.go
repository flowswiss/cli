package compute

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "image",
		Aliases: []string{"images"},
		Short:   "Manage images",
		Example: commands.FormatExamples(fmt.Sprintf(`
	  		# List all available images
			%[1]s compute image list
		`, commands.Name)),
	}

	commands.Add(cmd, &imageListCommand{})

	return cmd
}

type imageListCommand struct {
	filter string
}

func (i *imageListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.Images(cmd.Context(), commands.Config.Client)
	if err != nil {
		return fmt.Errorf("fetch images: %w", err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Sorting < items[j].Sorting
	})

	if len(i.filter) != 0 {
		items = filter.Find(items, i.filter)
	}

	return commands.PrintStdout(items)
}

func (i *imageListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (i *imageListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List server images",
		Long:              "Lists all server images.",
		ValidArgsFunction: i.CompleteArg,
		RunE:              i.Run,
	}

	cmd.Flags().StringVar(&i.filter, "filter", "", "custom term to filter the results")

	return cmd
}
