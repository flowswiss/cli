package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/internal/commands/dto"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/spf13/cobra"
	"sort"
)

const flagAll = "all"

var modulesCommand = &cobra.Command{
	Use:   "modules",
	Short: "List all modules",
	RunE:  listModules,
}

func init() {
	modulesCommand.Flags().Bool(flagAll, false, "list all modules instead of only those you have access to")
	modulesCommand.Flags().String(flagLocation, "", "filter for availability at location")
}

func findModule(filter string) (*flow.Module, error) {
	modules, _, err := client.Module.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	module, err := findOne(modules, filter, 2)
	if err != nil {
		return nil, fmt.Errorf("module: %v", err)
	}

	return module.(*flow.Module), nil
}

func listModules(cmd *cobra.Command, args []string) error {
	var modules []*flow.Module

	all, err := cmd.Flags().GetBool(flagAll)
	if err != nil {
		return err
	}

	if all {
		modules, _, err = client.Module.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
		if err != nil {
			return err
		}
	} else {
		organization, _, err := client.Organization.GetCurrent(context.Background())
		if err != nil {
			return err
		}

		modules = organization.RegisteredModules
	}

	locationFilter, err := cmd.Flags().GetString(flagLocation)
	if err != nil {
		return err
	}

	var location *flow.Location
	if locationFilter != "" {
		location, err = findLocation(locationFilter)
		if err != nil {
			return err
		}
	}

	sort.Sort(moduleBySorting{modules})

	var displayable []*dto.Module
	for _, module := range modules {
		if location != nil && !module.AvailableAt(location) {
			continue
		}

		displayable = append(displayable, &dto.Module{Module: module})
	}

	return display(displayable)
}

type moduleBySorting struct {
	Modules []*flow.Module
}

func (s moduleBySorting) Len() int {
	return len(s.Modules)
}

func (s moduleBySorting) Swap(i, j int) {
	s.Modules[i], s.Modules[j] = s.Modules[j], s.Modules[i]
}

func (s moduleBySorting) Less(i, j int) bool {
	return s.Modules[i].Sorting < s.Modules[j].Sorting
}
