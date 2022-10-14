package common

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func Product(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product",
		Short: "Manage products",
	}

	commands.Add(app, cmd, &productListCommand{})

	categoryCmd := &cobra.Command{
		Use:   "category",
		Short: "Manage product categories",
	}

	commands.Add(app, categoryCmd, &productCategoryListCommand{})
	cmd.AddCommand(categoryCmd)

	return cmd
}

type productListCommand struct {
	filter string
}

func (p *productListCommand) Run(cmd *cobra.Command, args []string) (err error) {
	var items []common.Product

	if len(args) != 0 {
		items, err = common.ProductsByType(cmd.Context(), commands.Config.Client, args[0])
		if err != nil {
			return err
		}
	} else {
		items, err = common.Products(cmd.Context(), commands.Config.Client)
		if err != nil {
			return err
		}
	}

	if len(p.filter) != 0 {
		items = filter.Find(items, p.filter)
	}

	return commands.PrintStdout(items)
}

func (p *productListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeProductCategory(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (p *productListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list [CATEGORY]",
		Short:             "List products",
		Long:              "Lists all available products.",
		Example:           "", // TODO
		ValidArgsFunction: p.CompleteArg,
		RunE:              p.Run,
	}

	cmd.Flags().StringVar(&p.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type productCategoryListCommand struct {
	filter string
}

func (p *productCategoryListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := common.ProductTypes(cmd.Context(), commands.Config.Client)
	if err != nil {
		return err
	}

	if len(p.filter) != 0 {
		items = filter.Find(items, p.filter)
	}

	return commands.PrintStdout(items)
}

func (p *productCategoryListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (p *productCategoryListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Short:             "List product categories",
		Long:              "Lists all product categories.",
		Example:           "", // TODO
		ValidArgsFunction: p.CompleteArg,
		RunE:              p.Run,
	}

	cmd.Flags().StringVar(&p.filter, "filter", "", "custom term to filter the results")

	return cmd
}

func completeProductCategory(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	categories, err := common.ProductTypes(ctx, commands.Config.Client)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(categories, term)

	names := make([]string, len(filtered))
	for i, category := range filtered {
		names[i] = category.Key
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
