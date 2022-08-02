package common

import (
	"context"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/common"
)

var ErrOrderFailed = common.ErrOrderFailed

type Ordering = common.Ordering

type Order = common.Order

func WaitForOrder(ctx context.Context, client goclient.Client, ordering Ordering) (Order, error) {
	return common.NewOrderService(client).WaitUntilProcessed(ctx, ordering)
}
