package flow

import (
	"context"
	"fmt"
	"net/http"
)

type ElasticIpService interface {
	List(ctx context.Context, options PaginationOptions) ([]*ElasticIp, *Response, error)
	Create(ctx context.Context, data *ElasticIpCreate) (*ElasticIp, *Response, error)
	Delete(ctx context.Context, id Id) (*Response, error)
}

type ElasticIpProduct struct {
	Id   Id     `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ElasticIp struct {
	Id               Id               `json:"id"`
	Product          ElasticIpProduct `json:"product"`
	Location         Location         `json:"location"`
	Price            float64          `json:"price"`
	PublicIp         string           `json:"public_ip"`
	PrivateIp        string           `json:"private_ip"`
	AttachedInstance *Server          `json:"attached_instance"`
}

type ElasticIpCreate struct {
	LocationId Id `json:"location_id"`
}

func (e *ElasticIp) String() string {
	return e.PublicIp
}

type elasticIpService struct {
	client *Client
}

func (s *elasticIpService) List(ctx context.Context, options PaginationOptions) ([]*ElasticIp, *Response, error) {
	p := "/v3/organizations/{organization}/compute/elastic-ips"
	p, err := addOptions(p, options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, p, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	var val []*ElasticIp

	res, err := s.client.Do(req, &val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *elasticIpService) Create(ctx context.Context, data *ElasticIpCreate) (*ElasticIp, *Response, error) {
	p := "/v3/organizations/{organization}/compute/elastic-ips"

	req, err := s.client.NewRequest(ctx, http.MethodPost, p, data, 0)
	if err != nil {
		return nil, nil, err
	}

	val := &ElasticIp{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *elasticIpService) Delete(ctx context.Context, id Id) (*Response, error) {
	p := fmt.Sprintf("/v3/organizations/{organization}/compute/elastic-ips/%d", id)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, p, nil, 0)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}
