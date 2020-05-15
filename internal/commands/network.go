package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/internal/commands/dto"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/spf13/cobra"
)

var (
	networkCommand = &cobra.Command{
		Use:   "network",
		Short: "Manage compute networks",
	}

	networkListCommand = &cobra.Command{
		Use:   "list",
		Short: "List all networks",
		RunE:  listNetwork,
	}
)

func init() {
	networkCommand.AddCommand(networkListCommand)
}

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

func listNetwork(cmd *cobra.Command, args []string) error {
	networks, _, err := client.Network.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	var displayable []*dto.Network
	for _, network := range networks {
		displayable = append(displayable, &dto.Network{Network: network})
	}

	return display(displayable)
}
