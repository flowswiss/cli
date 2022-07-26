package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/kubernetes"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type Node kubernetes.Node

func (n Node) String() string {
	return n.Name
}

func (n Node) Keys() (identifiers []string) {
	return []string{fmt.Sprint(n.ID), n.Name}
}

func (n Node) Columns() []string {
	return []string{"id", "name", "status", "roles", "product", "network"}
}

func (n Node) Values() map[string]interface{} {
	networkBuffer := strings.Builder{}

	networkBuffer.WriteString(fmt.Sprintf("%s (", n.Network.Name))
	for j, iface := range n.Network.Interfaces {
		if j != 0 {
			networkBuffer.WriteString(", ")
		}

		networkBuffer.WriteString(iface.PrivateIP)
	}
	networkBuffer.WriteRune(')')

	roleBuffer := strings.Builder{}
	for idx, role := range n.Roles {
		if idx != 0 {
			roleBuffer.WriteString(", ")
		}

		roleBuffer.WriteString(role.Name)
	}

	return map[string]interface{}{
		"id":      n.ID,
		"name":    n.Name,
		"status":  n.Status.Name,
		"roles":   roleBuffer.String(),
		"product": common.Product(n.Product),
		"network": networkBuffer.String(),
	}
}

type NodeService struct {
	delegate kubernetes.NodeService
}

func NewNodeService(client goclient.Client, clusterID int) NodeService {
	return NodeService{
		delegate: kubernetes.NewNodeService(client, clusterID),
	}
}

func (n NodeService) List(ctx context.Context) ([]Node, error) {
	res, err := n.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Node, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Node(item)
	}

	return items, nil
}

func (n NodeService) Delete(ctx context.Context, id int) error {
	return n.delegate.Delete(ctx, id)
}

type NodePerformAction = kubernetes.NodePerformAction

func (n NodeService) PerformAction(ctx context.Context, id int, data NodePerformAction) (Node, error) {
	cluster, err := n.delegate.PerformAction(ctx, id, data)
	if err != nil {
		return Node{}, err
	}

	return Node(cluster), nil
}
