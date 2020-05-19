package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/internal/commands/dto"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/spf13/cobra"
	"sort"
)

const flagModule = "module"

var locationsCommand = &cobra.Command{
	Use:   "locations",
	Short: "List locations",
	RunE:  listLocations,
}

func init() {
	locationsCommand.Flags().String(flagModule, "", "filter for available module")
}

func findLocation(filter string) (*flow.Location, error) {
	locations, _, err := client.Location.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	location, err := findOne(locations, filter, 2)
	if err != nil {
		return nil, fmt.Errorf("location: %v", err)
	}

	return location.(*flow.Location), nil
}

func listLocations(cmd *cobra.Command, args []string) error {
	locations, _, err := client.Location.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	moduleFilter, err := cmd.Flags().GetString(flagModule)
	if err != nil {
		return err
	}

	var module *flow.Module
	if moduleFilter != "" {
		module, err = findModule(moduleFilter)
		if err != nil {
			return err
		}
	}

	modules, _, err := client.Module.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	sort.Sort(moduleBySorting{modules})

	var displayable []*dto.Location
	for _, location := range locations {
		if module != nil && !module.AvailableAt(location) {
			continue
		}

		displayable = append(displayable, &dto.Location{Location: location, Modules: modules})
	}

	return display(displayable)
}
