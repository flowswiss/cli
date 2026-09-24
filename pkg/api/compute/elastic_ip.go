package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
	"github.com/flowswiss/goclient/v2/core"
)

type (
	ElasticIPCreate = compute.ElasticIPCreateReq
	ElasticIPList   = core.Cursor
	ElasticIPDelete = compute.ElasticIPDeleteReq
	ElasticIPAttach = compute.ElasticIPAttachmentCreateReq
	ElasticIPDetach = compute.ElasticIPAttachmentDeleteReq
)

type GenericElasticIPService struct {
	generic.CFD[
		compute.ElasticIP,
		ElasticIP,
		ElasticIPCreate,
		ElasticIPList,
		ElasticIPDelete,
	]

	client *compute.ElasticIPAttachmentService
}

func (s GenericElasticIPService) Attach(ctx context.Context, req ElasticIPAttach) (ElasticIP, error) {
	elasticIP, err := s.client.Create(ctx, req)
	if err != nil {
		return ElasticIP{}, err
	}

	return ElasticIP(elasticIP), err
}

func (s GenericElasticIPService) Detach(ctx context.Context, req ElasticIPDetach) error {
	return s.client.Delete(ctx, req)
}

func ElasticIPService() GenericElasticIPService {
	return GenericElasticIPService{
		generic.NewCFD[
			compute.ElasticIP,
			ElasticIP,
			ElasticIPCreate,
			ElasticIPList,
			ElasticIPDelete,
		](commands.Client.Compute.ElasticIP, func(ip compute.ElasticIP) ElasticIP {
			return ElasticIP(ip)
		}),
		commands.Client.Compute.ElasticIPAttachment,
	}
}

type ElasticIP compute.ElasticIP

func (e ElasticIP) String() string {
	return e.PublicIP
}

func (e ElasticIP) Keys() []string {
	return []string{fmt.Sprint(e.ID), e.PublicIP, e.PrivateIP}
}

func (e ElasticIP) Columns() []string {
	return []string{"id", "location", "public ip", "attachment"}
}

func (e ElasticIP) Values() map[string]any {
	attachment := ""
	if e.Attachment.ID != 0 {
		attachment = fmt.Sprintf("%s (%s)", e.Attachment.Name, e.PrivateIP)
	}

	return map[string]any{
		"id":         e.ID,
		"location":   e.Location.Name,
		"public ip":  e.PublicIP,
		"attachment": attachment,
	}
}
