package dto

import (
	"bytes"
	"github.com/flowswiss/cli/pkg/flow"
)

type Module struct {
	*flow.Module
}

func (m *Module) Columns() []string {
	return []string{"id", "name", "parent", "locations"}
}

func (m *Module) Values() map[string]interface{} {
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
		"id":        m.Id,
		"name":      m.Name,
		"parent":    parent,
		"locations": buf.String(),
	}
}
