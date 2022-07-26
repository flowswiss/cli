package common

import (
	"bytes"
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/common"

	"github.com/flowswiss/cli/v2/pkg/console"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

const HoursPerMonth = 730

const (
	ProductTypeComputeServer      = "compute-engine-vm"
	ProductTypeMacBareMetalDevice = "bare-metal-device"
	ProductTypeKubernetesNode     = "compute-kubernetes-node"
)

var (
	_ filter.Filterable   = (*Product)(nil)
	_ console.Displayable = (*Product)(nil)
)

type Product common.Product

func (p Product) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.PricePerHour())
}

func (p Product) PricePerHour() string {
	return fmt.Sprintf("%.2f CHF/h", p.Price/float64(HoursPerMonth))
}

func (p Product) Keys() []string {
	return []string{fmt.Sprint(p.ID), p.Name, p.Category, p.Type.Name}
}

func (p Product) Columns() []string {
	return []string{"id", "name", "type", "configuration", "price", "availability"}
}

func (p Product) Values() map[string]interface{} {
	configurationBuf := &bytes.Buffer{}
	for idx, item := range p.Items {
		if item.Description == "" {
			configurationBuf.WriteString(fmt.Sprintf("%d %s", item.Amount, item.Name))
		} else {
			configurationBuf.WriteString(fmt.Sprintf("%d %s %s", item.Amount, item.Description, item.Name))
		}

		if idx+1 < len(p.Items) {
			configurationBuf.WriteString(", ")
		}
	}

	availabilityBuf := &bytes.Buffer{}
	for idx, availability := range p.Availability {
		availabilityBuf.WriteString(availability.Location.Name)

		if idx+1 < len(p.Availability) {
			availabilityBuf.WriteString(", ")
		}
	}

	return map[string]interface{}{
		"id":            p.ID,
		"name":          p.Name,
		"type":          p.Type.Name,
		"configuration": configurationBuf.String(),
		"price":         p.PricePerHour(),
		"availability":  availabilityBuf.String(),
	}
}

func Products(ctx context.Context, client goclient.Client) ([]Product, error) {
	res, err := common.NewProductService(client).List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Product, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Product(item)
	}

	return items, nil
}

func ProductsByType(ctx context.Context, client goclient.Client, productTypeFilter string) ([]Product, error) {
	productTypes, err := ProductTypes(ctx, client)
	if err != nil {
		return nil, err
	}

	productType, err := filter.FindOne(productTypes, productTypeFilter)
	if err != nil {
		return nil, err
	}

	res, err := common.NewProductService(client).ListByType(ctx, productType.Key, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Product, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Product(item)
	}

	return items, nil
}

var (
	_ filter.Filterable   = (*ProductType)(nil)
	_ console.Displayable = (*ProductType)(nil)
)

type ProductType common.ProductType

func (p ProductType) Keys() []string {
	return []string{fmt.Sprint(p.ID), p.Name, p.Key}
}

func (p ProductType) Columns() []string {
	return []string{"id", "name", "key"}
}

func (p ProductType) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":   p.ID,
		"name": p.Name,
		"key":  p.Key,
	}
}

func ProductTypes(ctx context.Context, client goclient.Client) ([]ProductType, error) {
	res, err := common.NewProductService(client).ListTypes(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]ProductType, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = ProductType(item)
	}

	return items, nil
}
