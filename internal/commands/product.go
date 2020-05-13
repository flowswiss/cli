package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
)

func findProduct(filter string) (*flow.Product, error) {
	products, _, err := client.Product.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	product, err := findOne(products, filter, 2)
	if err != nil {
		return nil, fmt.Errorf("product: %v", err)
	}

	return product.(*flow.Product), nil
}
