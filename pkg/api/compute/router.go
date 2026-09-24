package compute

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	RouterCreate = compute.RouterCreateReq
	RouterGet    = compute.RouterGetReq
	RouterList   = core.Cursor
	RouterUpdate = compute.RouterUpdateReq
	RouterDelete = compute.RouterDeleteReq
)

type GenericRouterService = generic.CRUD[
	compute.Router,
	Router,
	RouterCreate,
	RouterGet,
	RouterList,
	RouterUpdate,
	RouterDelete,
]

func RouterService() GenericRouterService {
	return generic.NewCRUD[
		compute.Router,
		Router,
		RouterCreate,
		RouterGet,
		RouterList,
		RouterUpdate,
		RouterDelete,
	](commands.Client.Compute.Router, func(router compute.Router) Router {
		return Router(router)
	})
}

type Router compute.Router

func (r Router) String() string {
	return r.Name
}

func (r Router) Keys() []string {
	return []string{fmt.Sprint(r.ID), r.Name, r.PublicIP}
}

func (r Router) Columns() []string {
	return []string{"id", "name", "location", "public ip", "snat"}
}

func (r Router) Values() map[string]any {
	snat := "disabled"
	if r.SourceNAT {
		snat = "enabled"
	}

	return map[string]any{
		"id":        r.ID,
		"name":      r.Name,
		"location":  common.Location(r.Location),
		"public ip": r.PublicIP,
		"snat":      snat,
	}
}
