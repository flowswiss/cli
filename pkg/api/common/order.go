package common

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/common"
)

var ErrOrderFailed = errors.New("order failed")

type Ordering = common.Ordering

type Order = common.Order

func WaitForOrder(ctx context.Context, client goclient.Client, ordering Ordering) error {
	id, err := ordering.ExtractIdentifier()
	if err != nil {
		return fmt.Errorf("extract order id from %s: %w", ordering.Ref, err)
	}

	service := common.NewOrderService(client)
	timer := time.NewTimer(time.Second)

	for {
		order, err := service.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("fetch order: %w", err)
		}

		switch order.Status.Id {
		case common.OrderStatusSucceeded:
			return nil

		case common.OrderStatusFailed:
			return ErrOrderFailed
		}

		select {
		case <-timer.C:
			timer.Reset(time.Second)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
