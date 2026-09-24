package common

import (
	"fmt"

	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/cli/v2/pkg/filter"
	"github.com/flowswiss/goclient/v2"
	"github.com/flowswiss/goclient/v2/common"
	"github.com/flowswiss/goclient/v2/core"
)

var (
	_ filter.Filterable = (*Commitment)(nil)
)

type (
	Commitment common.Commitment

	CommitmentList   = core.Cursor
	CommitmentUpdate = common.CommitmentUpdateReq
)

type GenericCommitmentService struct {
	generic.ListService[CommitmentList, common.Commitment, Commitment]
	generic.UpdateService[CommitmentUpdate, common.Commitment, Commitment]
}

func CommitmentService(client *goclient.Client) GenericCommitmentService {
	c := client.Common.Commitment
	cast := func(e common.Commitment) Commitment {
		return Commitment(e)
	}

	return GenericCommitmentService{
		generic.NewList[CommitmentList, common.Commitment, Commitment](c, cast),
		generic.NewUpdate[CommitmentUpdate, common.Commitment, Commitment](c, cast),
	}
}

func (c Commitment) String() string {
	return fmt.Sprintf("Commitment for %s", c.Reference.Name)
}

func (c Commitment) Keys() []string {
	return []string{fmt.Sprint(c.ID), c.Reference.Type, c.Reference.Name, fmt.Sprint(c.Reference.ID)}
}

func (c Commitment) Columns() []string {
	return []string{"id", "type", "reference", "start_date", "end_date", "renew"}
}

func (c Commitment) Values() map[string]any {
	return map[string]any{
		"id":         c.ID,
		"type":       c.Reference.Type,
		"reference":  fmt.Sprintf("%s (ID %d)", c.Reference.Name, c.Reference.ID),
		"start_date": c.StartDate.String(),
		"end_date":   c.EndDate.String(),
		"renew":      c.Renew,
	}
}
