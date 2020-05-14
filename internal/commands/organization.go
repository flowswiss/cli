package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
)

func findOrganization(filter string) (*flow.Organization, error) {
	organizations, _, err := client.Organization.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	organization, err := findOne(organizations, filter, 2)
	if err != nil {
		return nil, fmt.Errorf("server: %v", err)
	}

	return organization.(*flow.Organization), nil
}
