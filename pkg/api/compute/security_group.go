package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

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

func (s SecurityGroup) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":       s.ID,
		"name":     s.Name,
		"location": common.Location(s.Location),
	}
}

type SecurityGroupService struct {
	delegate compute.SecurityGroupService
}

func NewSecurityGroupService(client goclient.Client) SecurityGroupService {
	return SecurityGroupService{
		delegate: compute.NewSecurityGroupService(client),
	}
}

func (s SecurityGroupService) List(ctx context.Context) ([]SecurityGroup, error) {
	res, err := s.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]SecurityGroup, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = SecurityGroup(item)
	}

	return items, nil
}

type SecurityGroupCreate = compute.SecurityGroupCreate

func (s SecurityGroupService) Create(ctx context.Context, data SecurityGroupCreate) (SecurityGroup, error) {
	res, err := s.delegate.Create(ctx, data)
	if err != nil {
		return SecurityGroup{}, err
	}

	return SecurityGroup(res), nil
}

type SecurityGroupUpdate = compute.SecurityGroupUpdate

func (s SecurityGroupService) Update(ctx context.Context, id int, data SecurityGroupUpdate) (SecurityGroup, error) {
	res, err := s.delegate.Update(ctx, id, data)
	if err != nil {
		return SecurityGroup{}, err
	}

	return SecurityGroup(res), nil
}

func (s SecurityGroupService) Delete(ctx context.Context, id int) error {
	return s.delegate.Delete(ctx, id)
}
