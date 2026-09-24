package generic

import (
	"context"

	"github.com/flowswiss/cli/v2/pkg/filter"
	"github.com/flowswiss/goclient/v2/common"
)

// Create

type createService[CreateInput, CreateOutput any] interface {
	Create(context.Context, CreateInput) (CreateOutput, error)
}

type CreateService[Input, RawOutput, CastedOutput any] struct {
	service createService[Input, RawOutput]
	cast    Cast[RawOutput, CastedOutput]
}

func (c CreateService[Input, RawOutput, CastedOutput]) Create(ctx context.Context, in Input) (CastedOutput, error) {
	result, err := c.service.Create(ctx, in)
	return c.cast(result), err
}

func NewCreate[Input, RawOutput, CastedOutput any](
	service createService[Input, RawOutput],
	cast Cast[RawOutput, CastedOutput],
) CreateService[Input, RawOutput, CastedOutput] {
	return CreateService[Input, RawOutput, CastedOutput]{
		service: service,
		cast:    cast,
	}
}

type orderedCreateService[CreateInput any] interface {
	Create(context.Context, CreateInput) (common.Ordering, error)
}

type OrderedCreateService[Input any] struct {
	service orderedCreateService[Input]
}

func (c OrderedCreateService[Input]) Create(ctx context.Context, in Input) (common.Ordering, error) {
	return c.service.Create(ctx, in)
}

func NewOrderedCreate[Input any](
	service orderedCreateService[Input],
) OrderedCreateService[Input] {
	return OrderedCreateService[Input]{
		service: service,
	}
}

// Read

type getService[GetInput, GetOutput any] interface {
	Get(context.Context, GetInput) (GetOutput, error)
}

type GetService[Input, RawOutput, CastedOutput any] struct {
	service getService[Input, RawOutput]
	cast    Cast[RawOutput, CastedOutput]
}

func (g GetService[Input, RawOutput, CastedOutput]) Get(ctx context.Context, in Input) (c CastedOutput, r error) {
	item, err := g.service.Get(ctx, in)
	if err != nil {
		return c, err
	}

	return g.cast(item), nil
}

type listService[ListInput, ListOutput any] interface {
	List(context.Context, ListInput) (common.List[ListOutput], error)
}

type ListService[Input, RawOutput, CastedOutput any] struct {
	service listService[Input, RawOutput]
	cast    Cast[RawOutput, CastedOutput]
}

func (l ListService[Input, RawOutput, CastedOutput]) List(ctx context.Context, in Input) ([]CastedOutput, error) {
	listOutput, err := l.service.List(ctx, in)
	if err != nil {
		return nil, err
	}

	items := make([]CastedOutput, len(listOutput.Items))
	for idx, item := range listOutput.Items {
		items[idx] = l.cast(item)
	}

	return items, nil
}

type FilterService[Input, RawOutput any, CastedOutput filter.Filterable] struct {
	ListService[Input, RawOutput, CastedOutput]
}

func (f FilterService[Input, RawOutput, CastedOutput]) Filter(
	ctx context.Context,
	in Input,
	appliedFilter string,
) ([]CastedOutput, error) {
	items, err := f.ListService.List(ctx, in)
	if err != nil {
		return nil, err
	}

	if len(appliedFilter) != 0 {
		items = filter.Find(items, appliedFilter)
	}

	return items, nil
}

func (f FilterService[Input, RawOutput, CastedOutput]) FindOne(
	ctx context.Context,
	in Input,
	appliedFilter string,
) (c CastedOutput, r error) {
	items, err := f.ListService.List(ctx, in)
	if err != nil {
		return c, err
	}

	item, err := filter.FindOne(items, appliedFilter)
	if err != nil {
		return c, err
	}

	return item, nil
}

func NewGet[Input, RawOutput, CastedOutput any](
	service getService[Input, RawOutput],
	cast Cast[RawOutput, CastedOutput],
) GetService[Input, RawOutput, CastedOutput] {
	return GetService[Input, RawOutput, CastedOutput]{
		service: service,
		cast:    cast,
	}
}

func NewList[Input, RawOutput, CastedOutput any](
	service listService[Input, RawOutput],
	cast Cast[RawOutput, CastedOutput],
) ListService[Input, RawOutput, CastedOutput] {
	return ListService[Input, RawOutput, CastedOutput]{
		service: service,
		cast:    cast,
	}
}

func NewFilter[Input, RawOutput any, CastedOutput filter.Filterable](
	service listService[Input, RawOutput],
	cast Cast[RawOutput, CastedOutput],
) FilterService[Input, RawOutput, CastedOutput] {
	return FilterService[Input, RawOutput, CastedOutput]{
		ListService: ListService[Input, RawOutput, CastedOutput]{
			service: service,
			cast:    cast,
		},
	}
}

// Update

type updateService[UpdateInput, UpdateOutput any] interface {
	Update(context.Context, UpdateInput) (UpdateOutput, error)
}

type UpdateService[Input, RawOutput, CastedOutput any] struct {
	service updateService[Input, RawOutput]
	cast    Cast[RawOutput, CastedOutput]
}

func (u UpdateService[Input, RawOutput, CastedOutput]) Update(ctx context.Context, in Input) (c CastedOutput, r error) {
	item, err := u.service.Update(ctx, in)
	if err != nil {
		return c, err
	}

	return u.cast(item), nil
}

func NewUpdate[Input, RawOutput, CastedOutput any](
	service updateService[Input, RawOutput],
	cast Cast[RawOutput, CastedOutput],
) UpdateService[Input, RawOutput, CastedOutput] {
	return UpdateService[Input, RawOutput, CastedOutput]{
		service: service,
		cast:    cast,
	}
}

// Perform Action

type performActionService[ActionInput, ActionOutput any] interface {
	Perform(context.Context, ActionInput) (ActionOutput, error)
}

type PerformActionService[ActionInput, RawOutput, CastedOutput any] struct {
	service performActionService[ActionInput, RawOutput]
	cast    Cast[RawOutput, CastedOutput]
}

func (p PerformActionService[ActionInput, RawOutput, CastedOutput]) Perform(
	ctx context.Context,
	in ActionInput,
) (c CastedOutput, r error) {
	item, err := p.service.Perform(ctx, in)
	if err != nil {
		return c, err
	}

	return p.cast(item), nil
}

func NewPerformAction[Input, RawOutput, CastedOutput any](
	service performActionService[Input, RawOutput],
	cast Cast[RawOutput, CastedOutput],
) PerformActionService[Input, RawOutput, CastedOutput] {
	return PerformActionService[Input, RawOutput, CastedOutput]{
		service: service,
		cast:    cast,
	}
}

// Delete

type deleteService[DeleteInput any] interface {
	Delete(context.Context, DeleteInput) error
}

type DeleteService[Input any] struct {
	service deleteService[Input]
}

func (d DeleteService[Input]) Delete(ctx context.Context, in Input) error {
	return d.service.Delete(ctx, in)
}

func NewDelete[Input any](
	service deleteService[Input],
) DeleteService[Input] {
	return DeleteService[Input]{
		service: service,
	}
}
