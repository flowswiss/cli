package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
)

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