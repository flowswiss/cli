package dto

import (
	"bytes"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
)

type Product struct {
	*flow.Product
}

func (p *Product) Columns() []string {
	return []string{"id", "name", "configuration", "price", "availability"}
}

func (p *Product) Values() map[string]interface{} {
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
		"id":            p.Id,
		"name":          p.Name,
		"configuration": configurationBuf.String(),
		"price":         getPricePerHour(p.Price),
		"availability":  availabilityBuf.String(),
	}
}

func getPricePerHour(price float64) string {
	return fmt.Sprintf("%.2f CHF/h", price/float64(730))
}
