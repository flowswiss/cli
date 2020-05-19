package dto

import (
	"bytes"
	"github.com/flowswiss/cli/pkg/flow"
)

type Location struct {
	*flow.Location
	Modules []*flow.Module
}

func (l *Location) Columns() []string {
	return []string{"id", "name", "city", "modules"}
}

func (l *Location) Values() map[string]interface{} {
	buf := &bytes.Buffer{}
	for _, module := range l.Modules {
		if module.AvailableAt(l.Location) {
			buf.WriteString(", ")
			buf.WriteString(module.Name)
		}
	}

	modules := buf.String()
	if modules != "" {
		modules = modules[2:]
	}

	return map[string]interface{}{
		"id":      l.Id,
		"name":    l.Name,
		"city":    l.City,
		"modules": modules,
	}
}
