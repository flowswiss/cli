package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

type ElasticIP compute.ElasticIP

func (e ElasticIP) Keys() []string {
	return []string{fmt.Sprint(e.ID), e.PublicIP, e.PrivateIP}
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
	delegate compute.ElasticIPService
}

func NewElasticIPService(client goclient.Client) ElasticIPService {
	return ElasticIPService{
		delegate: compute.NewElasticIPService(client),
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

func (e ElasticIPService) Delete(ctx context.Context, id int) error {
	return e.delegate.Delete(ctx, id)
}
