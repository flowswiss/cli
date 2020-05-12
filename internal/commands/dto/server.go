package dto

import (
	"bytes"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
)

type Server struct {
	*flow.Server
}

func (s *Server) Columns() []string {
	return []string{"id", "name", "status", "product", "operating system", "location", "public ip", "network"}
}

func (s *Server) Values() map[string]interface{} {
	networkBuffer := &bytes.Buffer{}
	publicIpBuffer := &bytes.Buffer{}

	for i, network := range s.Networks {
		networkBuffer.WriteString(fmt.Sprintf("%s (", network.Name))
		for j, iface := range network.Interfaces {
			networkBuffer.WriteString(iface.PrivateIp)

			if iface.PublicIp != "" {
				publicIpBuffer.WriteString(fmt.Sprintf("%s, ", iface.PublicIp))
			}

			if j+1 < len(network.Interfaces) {
				networkBuffer.WriteString(", ")
			}
		}
		networkBuffer.WriteRune(')')

		if i+1 < len(s.Networks) {
			networkBuffer.WriteString(", ")
		}
	}

	publicIp := publicIpBuffer.String()
	if len(publicIp) > 0 {
		publicIp = publicIp[:len(publicIp)-2]
	}

	pricePerHour := s.Product.Price / float64(730)

	return map[string]interface{}{
		"id":               s.Id,
		"name":             s.Name,
		"status":           s.Status.Name,
		"location":         s.Location.Name,
		"product":          fmt.Sprintf("%s (%.2f CHF/h)", s.Product.Name, pricePerHour),
		"operating system": fmt.Sprintf("%s %s", s.Image.OperatingSystem, s.Image.Version),
		"public ip":        publicIp,
		"network":          networkBuffer.String(),
	}
}
