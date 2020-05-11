package flow

import (
	"context"
	"fmt"
	"net/http"
)

type LocationService interface {
	List(ctx context.Context, options PaginationOptions) ([]*Location, *Response, error)
	Get(ctx context.Context, id Id) (*Location, *Response, error)
}

type Location struct {
	Id   Id     `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	City string `json:"city"`
}

func (l *Location) String() string {
	return l.Name
}

type locationService struct {
	client *Client
}

func (s *locationService) List(ctx context.Context, options PaginationOptions) ([]*Location, *Response, error) {
	path, err := addOptions("/v3/entities/locations", options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil, FlagNoAuthentication)
	if err != nil {
		return nil, nil, err
	}

	var val []*Location

	res, err := s.client.Do(req, &val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *locationService) Get(ctx context.Context, id Id) (*Location, *Response, error) {
	path := fmt.Sprintf("/v3/entities/locations/%d", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil, FlagNoAuthentication)
	if err != nil {
		return nil, nil, err
	}

	val := &Location{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}
