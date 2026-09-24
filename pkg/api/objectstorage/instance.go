package objectstorage

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/flowswiss/goclient/v2/objectstorage"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	InstanceCreate = objectstorage.InstanceCreateReq
	InstanceList   = core.Cursor
	InstanceDelete = objectstorage.InstanceDeleteReq
)

type GenericInstanceService = generic.CFD[objectstorage.Instance, Instance, InstanceCreate, InstanceList, InstanceDelete]

func InstanceService() GenericInstanceService {
	return generic.NewCFD[
		objectstorage.Instance,
		Instance,
		InstanceCreate,
		InstanceList,
		InstanceDelete,
	](commands.Client.ObjectStorage.Instance, func(instance objectstorage.Instance) Instance {
		return Instance(instance)
	})
}

type Instance objectstorage.Instance

func (i Instance) String() string {
	return i.Name
}

func (i Instance) Keys() []string {
	keys := []string{fmt.Sprint(i.ID), i.Name}
	keys = append(keys, common.Location(i.Location).Keys()...)
	return keys
}

func (i Instance) Columns() []string {
	return []string{"id", "name", "location"}
}

func (i Instance) Values() map[string]any {
	return map[string]any{
		"id":       i.ID,
		"name":     i.Name,
		"location": common.Location(i.Location),
	}
}
