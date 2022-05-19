package common

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/internal/commands"
	"github.com/flowswiss/cli/pkg/api/common"
	"github.com/flowswiss/cli/pkg/filter"
)

func ProductCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product",
		Short: "Manage products",
	}

	commands.Add(cmd, &productListCommand{})

	categoryCmd := &cobra.Command{
		Use:   "category",
		Short: "Manage product categories",
	}

	commands.Add(categoryCmd, &productCategoryListCommand{})
	cmd.AddCommand(categoryCmd)

	return cmd
}

type productListCommand struct {
	productType string
	filter      string
}

func (p *productListCommand) Run(ctx context.Context, config commands.Config, args []string) (err error) {
	var items []common.Product

	if len(p.productType) != 0 {
		items, err = common.ProductsByType(ctx, config.Client, p.productType)
		if err != nil {
			return err
		}
	} else {
		items, err = common.Products(ctx, config.Client)
		if err != nil {
			return err
		}
	}

	if len(p.filter) != 0 {
		items = filter.Find(items, p.filter)
	}

	return commands.PrintStdout(items)
}

func (p *productListCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List products",
		Long:    "Lists all available products.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&p.productType, "category", "", "product category to filter the results")
	cmd.Flags().StringVar(&p.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type productCategoryListCommand struct {
	filter string
}

func (p *productCategoryListCommand) Run(ctx context.Context, config commands.Config, args []string) error {
	items, err := common.ProductTypes(ctx, config.Client)
	if err != nil {
		return err
	}

	if len(p.filter) != 0 {
		items = filter.Find(items, p.filter)
	}

	return commands.PrintStdout(items)
}

func (p *productCategoryListCommand) Desc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List product categories",
		Long:    "Lists all product categories.",
		Example: "", // TODO
	}

	cmd.Flags().StringVar(&p.filter, "filter", "", "custom term to filter the results")

	return cmd
}
