package macbaremetal

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/flowswiss/goclient/v2/macbaremetal"
)

type (
	RouterGet    = macbaremetal.RouterGetReq
	RouterList   = core.Cursor
	RouterUpdate = macbaremetal.RouterUpdateReq
)

type GenericRouterService struct {
	generic.Read[macbaremetal.Router, Router, RouterGet, RouterList]
	generic.UpdateService[RouterUpdate, macbaremetal.Router, Router]
}

func RouterService() GenericRouterService {
	client := commands.Client.MacBareMetal.Router
	cast := func(router macbaremetal.Router) Router {
		return Router(router)
	}

	return GenericRouterService{
		generic.NewRead[macbaremetal.Router, Router, RouterGet, RouterList](client, cast),
		generic.NewUpdate[RouterUpdate, macbaremetal.Router, Router](client, cast),
	}
}

type Router macbaremetal.Router

func (r Router) Keys() []string {
	return []string{fmt.Sprint(r.ID), r.Name, r.PublicIP}
}

func (r Router) Columns() []string {
	return []string{"id", "name", "location", "public ip"}
}

func (r Router) Values() map[string]any {
	return map[string]any{
		"id":        r.ID,
		"name":      r.Name,
		"location":  r.Location.Name,
		"public ip": r.PublicIP,
	}
}
