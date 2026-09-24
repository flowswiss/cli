package common

import (
	"fmt"

	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/cli/v2/pkg/filter"
	"github.com/flowswiss/goclient/v2"
	"github.com/flowswiss/goclient/v2/common"
)

var (
	_ filter.Filterable = (*Quota)(nil)
)

type (
	QuotaOverview common.Quota
	Quota         common.QuotaUsage

	QuotaGet = common.QuotaGetReq
)

type GenericQuotaService struct {
	generic.GetService[QuotaGet, common.Quota, QuotaOverview]
}

func QuotaService(client *goclient.Client) GenericQuotaService {
	c := client.Common.Quota
	cast := func(q common.Quota) QuotaOverview {
		return QuotaOverview(q)
	}

	return GenericQuotaService{
		generic.NewGet[QuotaGet, common.Quota, QuotaOverview](c, cast),
	}
}

func (o QuotaOverview) Quotas() []Quota {
	result := make([]Quota, len(o.Usages))
	for i, u := range o.Usages {
		result[i] = Quota(u)
	}
	return result
}

func (q Quota) String() string {
	return fmt.Sprintf("Quota for %s", q.Entity.Name)
}

func (q Quota) Keys() []string {
	return []string{q.Entity.Key, q.Entity.Name, q.EntityType.Key, q.EntityType.Name}
}

func (q Quota) Columns() []string {
	return []string{"entity", "type", "amount", "unit"}
}

func (q Quota) Values() map[string]any {
	return map[string]any{
		"entity": q.Entity.Name,
		"type":   q.EntityType.Name,
		"amount": fmt.Sprintf("%d / %d", q.Current, q.Quota),
		"unit":   q.Unit,
	}
}
