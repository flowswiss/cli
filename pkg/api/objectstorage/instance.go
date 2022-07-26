package objectstorage

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/objectstorage"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

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

func (i Instance) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":       i.ID,
		"name":     i.Name,
		"location": common.Location(i.Location),
	}
}

type InstanceService struct {
	delegate objectstorage.InstanceService
}

func NewInstanceService(client goclient.Client) InstanceService {
	return InstanceService{
		delegate: objectstorage.NewInstanceService(client),
	}
}

func (i InstanceService) List(ctx context.Context) ([]Instance, error) {
	res, err := i.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Instance, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Instance(item)
	}

	return items, nil
}

type InstanceCreate = objectstorage.InstanceCreate

func (i InstanceService) Create(ctx context.Context, data InstanceCreate) (Instance, error) {
	res, err := i.delegate.Create(ctx, data)
	if err != nil {
		return Instance{}, err
	}

	return Instance(res), nil
}

func (i InstanceService) Delete(ctx context.Context, id int) error {
	return i.delegate.Delete(ctx, id)
}
