package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/internal/commands/dto"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/spf13/cobra"
)

var (
	productsCommand = &cobra.Command{
		Use:   "products",
		Short: "Lists products of type",
	}

	productsListComputeCommand = &cobra.Command{
		Use:   "compute",
		Short: "Lists all compute products",
		RunE:  listComputeProducts,
	}

	productsListMacBareMetalCommand = &cobra.Command{
		Use:   "mac-bare-metal",
		Short: "Lists all mac bare metal products",
		RunE:  listMacBareMetalProducts,
	}
)

func init() {
	productsCommand.AddCommand(productsListComputeCommand)
	productsCommand.AddCommand(productsListMacBareMetalCommand)

	productsCommand.PersistentFlags().String(flagLocation, "", "filter for availability at location")
}

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

func listProducts(cmd *cobra.Command, t string) error {
	products, _, err := client.Product.ListByType(context.Background(), flow.PaginationOptions{NoFilter: 1}, t)
	if err != nil {
		return err
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

	var displayable []*dto.Product
	for _, product := range products {
		if len(product.Availability) == 0 {
			continue
		}

		if location != nil && !product.AvailableAt(location) {
			continue
		}

		displayable = append(displayable, &dto.Product{Product: product})
	}

	return display(displayable)
}

func listComputeProducts(cmd *cobra.Command, args []string) error {
	return listProducts(cmd, "compute-engine-vm")
}

func listMacBareMetalProducts(cmd *cobra.Command, args []string) error {
	return listProducts(cmd, "bare-metal-device")
}
