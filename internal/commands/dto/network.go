package dto

import (
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
)

type Network struct {
	*flow.Network
}

func (n *Network) Columns() []string {
	return []string{"id", "name", "location", "cidr", "usage"}
}

func (n *Network) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":       n.Id,
		"name":     n.Name,
		"location": n.Location.Name,
		"cidr":     n.Cidr,
		"usage":    fmt.Sprintf("%d/%d", n.UsedIps, n.TotalIps),
	}
}
