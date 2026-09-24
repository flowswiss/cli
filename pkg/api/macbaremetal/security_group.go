package macbaremetal

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/flowswiss/goclient/v2/macbaremetal"
)

type (
	SecurityGroupCreate = macbaremetal.SecurityGroupCreateReq
	SecurityGroupGet    = macbaremetal.SecurityGroupGetReq
	SecurityGroupList   = core.Cursor
	SecurityGroupUpdate = macbaremetal.SecurityGroupUpdateReq
	SecurityGroupDelete = macbaremetal.SecurityGroupDeleteReq
)

type GenericSecurityGroupService = generic.CRUD[
	macbaremetal.SecurityGroup,
	SecurityGroup,
	SecurityGroupCreate,
	SecurityGroupGet,
	SecurityGroupList,
	SecurityGroupUpdate,
	SecurityGroupDelete,
]

func SecurityGroupService() GenericSecurityGroupService {
	return generic.NewCRUD[
		macbaremetal.SecurityGroup,
		SecurityGroup,
		SecurityGroupCreate,
		SecurityGroupGet,
		SecurityGroupList,
		SecurityGroupUpdate,
		SecurityGroupDelete,
	](commands.Client.MacBareMetal.SecurityGroup, func(group macbaremetal.SecurityGroup) SecurityGroup {
		return SecurityGroup(group)
	})
}

type SecurityGroup macbaremetal.SecurityGroup

func (s SecurityGroup) String() string {
	return s.Name
}

func (s SecurityGroup) Keys() []string {
	return []string{fmt.Sprint(s.ID), s.Name}
}

func (s SecurityGroup) Columns() []string {
	return []string{"id", "name", "network"}
}

func (s SecurityGroup) Values() map[string]any {
	return map[string]any{
		"id":      s.ID,
		"name":    s.Name,
		"network": s.Network.Name,
	}
}
