package macbaremetal

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/macbaremetal"
)

type ElasticIP macbaremetal.ElasticIP

func (e ElasticIP) Keys() []string {
	return []string{fmt.Sprint(e.ID), e.PublicIP}
}

func (e ElasticIP) Columns() []string {
	return []string{"id", "location", "public ip", "attachment"}
}

func (e ElasticIP) Values() map[string]interface{} {
	attachment := ""
	if e.Attachment.ID != 0 {
		attachment = fmt.Sprintf("%s (%s)", e.Attachment.Name, e.PrivateIP)
	}

	return map[string]interface{}{
		"id":         e.ID,
		"location":   e.Location.Name,
		"public ip":  e.PublicIP,
		"attachment": attachment,
	}
}

type ElasticIPService struct {
	delegate macbaremetal.ElasticIPService
}

func NewElasticIPService(client goclient.Client) ElasticIPService {
	return ElasticIPService{
		delegate: macbaremetal.NewElasticIPService(client),
	}
}

func (e ElasticIPService) List(ctx context.Context) ([]ElasticIP, error) {
	res, err := e.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]ElasticIP, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = ElasticIP(item)
	}

	return items, nil
}

type ElasticIPCreate = macbaremetal.ElasticIPCreate

func (e ElasticIPService) Create(ctx context.Context, data ElasticIPCreate) (ElasticIP, error) {
	res, err := e.delegate.Create(ctx, data)
	if err != nil {
		return ElasticIP{}, err
	}

	return ElasticIP(res), nil
}

func (e ElasticIPService) Delete(ctx context.Context, id int) error {
	return e.delegate.Delete(ctx, id)
}
