package common

import (
	"bytes"
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/common"

	"github.com/flowswiss/cli/pkg/console"
	"github.com/flowswiss/cli/pkg/filter"
)

var (
	_ filter.Filterable   = (*Location)(nil)
	_ console.Displayable = (*Location)(nil)
)

type Location common.Location

func (l Location) Keys() []string {
	return []string{fmt.Sprint(l.Id), l.Name, l.City}
}

func (l Location) Columns() []string {
	return []string{"id", "name", "city", "modules"}
}

func (l Location) Values() map[string]interface{} {
	buf := &bytes.Buffer{}
	for idx, module := range l.Modules {
		buf.WriteString(module.Name)

		if idx+1 < len(l.Modules) {
			buf.WriteString(", ")
		}
	}

	return map[string]interface{}{
		"id":      l.Id,
		"name":    l.Name,
		"city":    l.City,
		"modules": buf.String(),
	}
}

func Locations(ctx context.Context, client goclient.Client) ([]Location, error) {
	res, err := common.NewLocationService(client).List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Location, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Location(item)
	}

	return items, nil
}
