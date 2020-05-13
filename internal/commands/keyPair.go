package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
)

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
