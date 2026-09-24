package kubernetes

import (
	"fmt"
	"strings"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/kubernetes"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	NodeList      = kubernetes.NodeListReq
	NodeDelete    = kubernetes.NodeDeleteReq
	NodeRunAction = kubernetes.NodePerformReq
)

type GenericNodeService struct {
	generic.ListService[NodeList, kubernetes.Node, Node]
	generic.PerformActionService[NodeRunAction, kubernetes.Node, Node]
	generic.DeleteService[NodeDelete]
}

func NodeService() GenericNodeService {
	client := commands.Client.Kubernetes.Node
	cast := func(node kubernetes.Node) Node {
		return Node(node)
	}

	return GenericNodeService{
		generic.NewList[NodeList, kubernetes.Node, Node](client, cast),
		generic.NewPerformAction[NodeRunAction, kubernetes.Node, Node](client, cast),
		generic.NewDelete[NodeDelete](client),
	}
}

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

func (n Node) Values() map[string]any {
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

	return map[string]any{
		"id":      n.ID,
		"name":    n.Name,
		"status":  n.Status.Name,
		"roles":   roleBuffer.String(),
		"product": common.Product(n.Product),
		"network": networkBuffer.String(),
	}
}

type NodeAction kubernetes.NodeAction

func (c NodeAction) Keys() []string {
	return []string{fmt.Sprint(c.ID), c.Name, c.Command}
}

func (c NodeAction) Columns() []string {
	return []string{"id", "name", "command"}
}

func (c NodeAction) Values() map[string]any {
	return map[string]any{
		"id":      c.ID,
		"name":    c.Name,
		"command": c.Command,
	}
}
