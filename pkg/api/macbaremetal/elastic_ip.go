package macbaremetal

import (
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/flowswiss/goclient/v2/macbaremetal"
)

type (
	ElasticIPCreate = macbaremetal.ElasticIPCreateReq
	ElasticIPList   = core.Cursor
	ElasticIPDelete = macbaremetal.ElasticIPDeleteReq
	ElasticIPAttach = macbaremetal.ElasticIPAttachmentCreateReq
	ElasticIPDetach = macbaremetal.ElasticIPAttachmentDeleteReq
)

type GenericElasticIPService struct {
	generic.CFD[macbaremetal.ElasticIP, ElasticIP, ElasticIPCreate, ElasticIPList, ElasticIPDelete]

	client *macbaremetal.ElasticIPAttachmentService
}

func (e GenericElasticIPService) Attach(ctx context.Context, req ElasticIPAttach) (ElasticIP, error) {
	elasticIP, err := e.client.Create(ctx, req)
	if err != nil {
		return ElasticIP{}, err
	}

	return ElasticIP(elasticIP), nil
}

func (e GenericElasticIPService) Detach(ctx context.Context, req ElasticIPDetach) error {
	return e.client.Delete(ctx, req)
}

func ElasticIPService() GenericElasticIPService {
	client := commands.Client.MacBareMetal.ElasticIP
	cast := func(elasticIP macbaremetal.ElasticIP) ElasticIP {
		return ElasticIP(elasticIP)
	}

	return GenericElasticIPService{
		generic.NewCFD[
			macbaremetal.ElasticIP,
			ElasticIP,
			ElasticIPCreate,
			ElasticIPList,
			ElasticIPDelete,
		](client, cast),
		commands.Client.MacBareMetal.ElasticIPAttachment,
	}
}

type ElasticIP macbaremetal.ElasticIP

func (e ElasticIP) String() string {
	return e.PublicIP
}

func (e ElasticIP) Keys() []string {
	return []string{fmt.Sprint(e.ID), e.PublicIP}
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
