package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
)

func findNetwork(filter string) (*flow.Network, error) {
	networks, _, err := client.Network.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	network, err := findOne(networks, filter, 2)
	if err != nil {
		return nil, fmt.Errorf("network: %v", err)
	}

	return network.(*flow.Network), nil
}
