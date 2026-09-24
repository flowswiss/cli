package generic

import "github.com/flowswiss/cli/v2/pkg/filter"

// Read (Get & Filter)
type Read[RawEntity any, CastedEntity filter.Filterable, Get, List any] struct {
	GetService[Get, RawEntity, CastedEntity]
	FilterService[List, RawEntity, CastedEntity]
}

type readService[Entity, Get, List any] interface {
	getService[Get, Entity]
	listService[List, Entity]
}

func NewRead[RawEntity any, CastedEntity filter.Filterable, Get, List any](
	service readService[RawEntity, Get, List],
	cast Cast[RawEntity, CastedEntity],
) Read[RawEntity, CastedEntity, Get, List] {
	return Read[RawEntity, CastedEntity, Get, List]{
		NewGet[Get, RawEntity, CastedEntity](service, cast),
		NewFilter[List, RawEntity, CastedEntity](service, cast),
	}
}

// CRD = Create - Read (Get & Filter) - Delete
type CRD[RawEntity any, CastedEntity filter.Filterable, Create, Get, List, Delete any] struct {
	CreateService[Create, RawEntity, CastedEntity]
	GetService[Get, RawEntity, CastedEntity]
	FilterService[List, RawEntity, CastedEntity]
	DeleteService[Delete]
}

type crdService[Entity, Create, Get, List, Delete any] interface {
	createService[Create, Entity]
	getService[Get, Entity]
	listService[List, Entity]
	deleteService[Delete]
}

func NewCRD[RawEntity any, CastedEntity filter.Filterable, Create, Get, List, Delete any](
	service crdService[RawEntity, Create, Get, List, Delete],
	cast Cast[RawEntity, CastedEntity],
) CRD[RawEntity, CastedEntity, Create, Get, List, Delete] {
	return CRD[RawEntity, CastedEntity, Create, Get, List, Delete]{
		NewCreate[Create, RawEntity, CastedEntity](service, cast),
		NewGet[Get, RawEntity, CastedEntity](service, cast),
		NewFilter[List, RawEntity, CastedEntity](service, cast),
		NewDelete[Delete](service),
	}
}

// CFD = Create - Filter - Delete
type CFD[RawEntity any, CastedEntity filter.Filterable, Create, List, Delete any] struct {
	CreateService[Create, RawEntity, CastedEntity]
	FilterService[List, RawEntity, CastedEntity]
	DeleteService[Delete]
}

type cfdService[Entity, Create, List, Delete any] interface {
	createService[Create, Entity]
	listService[List, Entity]
	deleteService[Delete]
}

func NewCFD[RawEntity any, CastedEntity filter.Filterable, Create, List, Delete any](
	service cfdService[RawEntity, Create, List, Delete],
	cast Cast[RawEntity, CastedEntity],
) CFD[RawEntity, CastedEntity, Create, List, Delete] {
	return CFD[RawEntity, CastedEntity, Create, List, Delete]{
		NewCreate[Create, RawEntity, CastedEntity](service, cast),
		NewFilter[List, RawEntity, CastedEntity](service, cast),
		NewDelete[Delete](service),
	}
}

// CRUD = Create - Read (Get & Filter) - Update - Delete
type CRUD[RawEntity any, CastedEntity filter.Filterable, Create, Get, List, Update, Delete any] struct {
	CreateService[Create, RawEntity, CastedEntity]
	GetService[Get, RawEntity, CastedEntity]
	FilterService[List, RawEntity, CastedEntity]
	UpdateService[Update, RawEntity, CastedEntity]
	DeleteService[Delete]
}

type crudService[Entity, Create, Get, List, Update, Delete any] interface {
	createService[Create, Entity]
	getService[Get, Entity]
	listService[List, Entity]
	updateService[Update, Entity]
	deleteService[Delete]
}

func NewCRUD[RawEntity any, CastedEntity filter.Filterable, Create, Get, List, Update, Delete any](
	service crudService[RawEntity, Create, Get, List, Update, Delete],
	cast Cast[RawEntity, CastedEntity],
) CRUD[RawEntity, CastedEntity, Create, Get, List, Update, Delete] {
	return CRUD[RawEntity, CastedEntity, Create, Get, List, Update, Delete]{
		NewCreate[Create, RawEntity, CastedEntity](service, cast),
		NewGet[Get, RawEntity, CastedEntity](service, cast),
		NewFilter[List, RawEntity, CastedEntity](service, cast),
		NewUpdate[Update, RawEntity, CastedEntity](service, cast),
		NewDelete[Delete](service),
	}
}
