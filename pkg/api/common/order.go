package common

import (
	"context"

	"github.com/flowswiss/goclient/v2"
	"github.com/flowswiss/goclient/v2/common"
)

var ErrOrderFailed = common.ErrOrderFailed

type Ordering = common.Ordering

type Order = common.Order

func WaitForOrder(ctx context.Context, client *goclient.Client, ordering Ordering) (Order, error) {
	return client.Common.Order.WaitUntilProcessed(ctx, ordering)
}
