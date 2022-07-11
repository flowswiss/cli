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
	_ filter.Filterable   = (*Module)(nil)
	_ console.Displayable = (*Module)(nil)
)

type Module common.Module

func (m Module) String() string {
	return m.Name
}

func (m Module) Keys() []string {
	return []string{fmt.Sprint(m.ID), m.Name}
}

func (m Module) Columns() []string {
	return []string{"id", "name", "parent", "locations"}
}

func (m Module) Values() map[string]interface{} {
	parent := ""
	if m.Parent != nil {
		parent = m.Parent.Name
	}

	buf := &bytes.Buffer{}
	for idx, location := range m.Locations {
		buf.WriteString(location.Name)

		if idx+1 < len(m.Locations) {
			buf.WriteString(", ")
		}
	}

	return map[string]interface{}{
		"id":        m.ID,
		"name":      m.Name,
		"parent":    parent,
		"locations": buf.String(),
	}
}

func (m Module) AvailableAt(location Location) bool {
	for _, loc := range m.Locations {
		if loc.ID == location.ID {
			return true
		}
	}
	return false
}

func Modules(ctx context.Context, client goclient.Client) ([]Module, error) {
	res, err := common.NewModuleService(client).List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Module, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Module(item)
	}

	return items, nil
}
