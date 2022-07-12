package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type Server compute.Server

func (s Server) String() string {
	return s.Name
}

func (s Server) Keys() (identifiers []string) {
	for _, network := range s.Networks {
		identifiers = append(identifiers, network.Name, network.CIDR)

		for _, iface := range network.Interfaces {
			identifiers = append(identifiers, fmt.Sprint(iface.ID), iface.PrivateIP, iface.PublicIP)
		}
	}

	return append(identifiers, fmt.Sprint(s.ID), s.Name)
}

func (s Server) Columns() []string {
	return []string{"id", "name", "status", "product", "operating system", "location", "public ip", "network"}
}

func (s Server) Values() map[string]interface{} {
	networkBuffer := &strings.Builder{}
	publicIPBuffer := &strings.Builder{}

	for i, network := range s.Networks {
		if i != 0 {
			networkBuffer.WriteString(", ")
		}

		networkBuffer.WriteString(fmt.Sprintf("%s (", network.Name))
		for j, iface := range network.Interfaces {
			if j != 0 {
				networkBuffer.WriteString(", ")
			}

			networkBuffer.WriteString(iface.PrivateIP)

			if iface.PublicIP != "" {
				publicIPBuffer.WriteString(fmt.Sprintf("%s, ", iface.PublicIP))
			}
		}
		networkBuffer.WriteRune(')')
	}

	publicIP := publicIPBuffer.String()
	if len(publicIP) > 0 {
		publicIP = publicIP[:len(publicIP)-2]
	}

	return map[string]interface{}{
		"id":               s.ID,
		"name":             s.Name,
		"status":           s.Status.Name,
		"product":          common.Product(s.Product),
		"operating system": Image{Image: s.Image},
		"location":         common.Location(s.Location),
		"public ip":        publicIP,
		"network":          networkBuffer.String(),
	}
}

type ServerService struct {
	client   goclient.Client
	delegate compute.ServerService
}

func NewServerService(client goclient.Client) ServerService {
	return ServerService{
		client:   client,
		delegate: compute.NewServerService(client),
	}
}

func (s ServerService) List(ctx context.Context) ([]Server, error) {
	res, err := s.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Server, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Server(item)
	}

	return items, nil
}

type ServerCreate = compute.ServerCreate

func (s ServerService) Create(ctx context.Context, data ServerCreate) (common.Ordering, error) {
	return s.delegate.Create(ctx, data)
}

type ServerUpdate = compute.ServerUpdate

func (s ServerService) Update(ctx context.Context, id int, data ServerUpdate) (Server, error) {
	res, err := s.delegate.Update(ctx, id, data)
	if err != nil {
		return Server{}, err
	}

	return Server(res), nil
}

type ServerUpgrade = compute.ServerUpgrade

func (s ServerService) Upgrade(ctx context.Context, id int, data ServerUpgrade) (common.Ordering, error) {
	return s.delegate.Upgrade(ctx, id, data)
}

func (s ServerService) Delete(ctx context.Context, id int, deleteElasticIPs bool) error {
	return s.delegate.Delete(ctx, id, deleteElasticIPs)
}
