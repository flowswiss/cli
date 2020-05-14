package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/internal/commands/dto"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/spf13/cobra"
)

var (
	keyPairCommand = &cobra.Command{
		Use:   "key-pair",
		Short: "Manage your ssh key pairs",
	}

	keyPairListCommand = &cobra.Command{
		Use:   "list",
		Short: "List all key pairs",
		RunE:  listKeyPair,
	}
)

func init() {
	keyPairCommand.AddCommand(keyPairListCommand)
}

func findKeyPair(filter string) (*flow.KeyPair, error) {
	keyPairs, _, err := client.KeyPair.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	keyPair, err := findOne(keyPairs, filter, 2)
	if err != nil {
		return nil, fmt.Errorf("key-pair: %v", err)
	}

	return keyPair.(*flow.KeyPair), nil
}

func listKeyPair(cmd *cobra.Command, args []string) error {
	keyPairs, _, err := client.KeyPair.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	var displayable []*dto.KeyPair
	for _, keyPair := range keyPairs {
		displayable = append(displayable, &dto.KeyPair{KeyPair: keyPair})
	}

	return display(displayable)
}
