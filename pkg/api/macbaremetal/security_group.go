package macbaremetal

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/macbaremetal"
)

type SecurityGroup macbaremetal.SecurityGroup

func (s SecurityGroup) Keys() []string {
	return []string{fmt.Sprint(s.ID), s.Name}
}

func (s SecurityGroup) Columns() []string {
	return []string{"id", "name", "network"}
}

func (s SecurityGroup) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":      s.ID,
		"name":    s.Name,
		"network": s.Network.Name,
	}
}

type SecurityGroupService struct {
	delegate macbaremetal.SecurityGroupService
}

func NewSecurityGroupService(client goclient.Client) SecurityGroupService {
	return SecurityGroupService{
		delegate: macbaremetal.NewSecurityGroupService(client),
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

type SecurityGroupCreate = macbaremetal.SecurityGroupCreate

func (s SecurityGroupService) Create(ctx context.Context, data SecurityGroupCreate) (SecurityGroup, error) {
	res, err := s.delegate.Create(ctx, data)
	if err != nil {
		return SecurityGroup{}, err
	}

	return SecurityGroup(res), nil
}

type SecurityGroupUpdate = macbaremetal.SecurityGroupUpdate

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
