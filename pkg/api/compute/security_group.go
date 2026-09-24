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
	SecurityGroupCreate = compute.SecurityGroupCreateReq
	SecurityGroupGet    = compute.SecurityGroupGetReq
	SecurityGroupList   = core.Cursor
	SecurityGroupUpdate = compute.SecurityGroupUpdateReq
	SecurityGroupDelete = compute.SecurityGroupDeleteReq
)

type GenericSecurityGroupService = generic.CRUD[
	compute.SecurityGroup,
	SecurityGroup,
	SecurityGroupCreate,
	SecurityGroupGet,
	SecurityGroupList,
	SecurityGroupUpdate,
	SecurityGroupDelete,
]

func SecurityGroupService() GenericSecurityGroupService {
	return generic.NewCRUD[
		compute.SecurityGroup,
		SecurityGroup,
		SecurityGroupCreate,
		SecurityGroupGet,
		SecurityGroupList,
		SecurityGroupUpdate,
		SecurityGroupDelete,
	](commands.Client.Compute.SecurityGroup, func(group compute.SecurityGroup) SecurityGroup {
		return SecurityGroup(group)
	})
}

type SecurityGroup compute.SecurityGroup

func (s SecurityGroup) String() string {
	return s.Name
}

func (s SecurityGroup) Keys() []string {
	return []string{fmt.Sprint(s.ID), s.Name}
}

func (s SecurityGroup) Columns() []string {
	return []string{"id", "name", "location"}
}

func (s SecurityGroup) Values() map[string]any {
	return map[string]any{
		"id":       s.ID,
		"name":     s.Name,
		"location": common.Location(s.Location),
	}
}
