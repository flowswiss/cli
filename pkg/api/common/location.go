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

var (
	_ filter.Filterable   = (*Location)(nil)
	_ console.Displayable = (*Location)(nil)
)

type Location common.Location

func (l Location) String() string {
	return l.Name
}

func (l Location) Keys() []string {
	return []string{fmt.Sprint(l.ID), l.Name, l.City}
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
		"id":      l.ID,
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

func FindLocation(ctx context.Context, client goclient.Client, term string) (Location, error) {
	locations, err := Locations(ctx, client)
	if err != nil {
		return Location{}, fmt.Errorf("fetch locations: %w", err)
	}

	location, err := filter.FindOne(locations, term)
	if err != nil {
		return Location{}, fmt.Errorf("find location: %w", err)
	}

	return location, nil
}
