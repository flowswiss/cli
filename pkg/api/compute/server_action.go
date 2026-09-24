package compute

import (
	"fmt"

	"github.com/flowswiss/goclient/v2/compute"
)

type ServerAction compute.ServerAction

func (d ServerAction) Keys() []string {
	return []string{fmt.Sprint(d.ID), d.Name, d.Command}
}

func (d ServerAction) Columns() []string {
	return []string{"id", "name", "command"}
}

func (d ServerAction) Values() map[string]any {
	return map[string]any{
		"id":      d.ID,
		"name":    d.Name,
		"command": d.Command,
	}
}
