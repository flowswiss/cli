package flow

import (
	"context"
	"fmt"
	"net/http"
)

type ServerAttachmentService interface {
	ListAttachedElasticIps(ctx context.Context, server Id, options PaginationOptions) ([]*ElasticIp, *Response, error)
	AttachElasticIp(ctx context.Context, server Id, data *AttachElasticIp) (*ElasticIp, *Response, error)
	DetachElasticIp(ctx context.Context, server Id, elasticIp Id) (*Response, error)
}

type AttachElasticIp struct {
	ElasticIpId        Id `json:"elastic_ip_id"`
	NetworkInterfaceId Id `json:"network_interface_id"`
}

type serverAttachmentService struct {
	client *Client
}

func (s *serverAttachmentService) ListAttachedElasticIps(ctx context.Context, server Id, options PaginationOptions) ([]*ElasticIp, *Response, error) {
	p := fmt.Sprintf("/v3/organizations/{organization}/compute/instances/%d/elastic-ips", server)
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

func (s *serverAttachmentService) AttachElasticIp(ctx context.Context, server Id, data *AttachElasticIp) (*ElasticIp, *Response, error) {
	p := fmt.Sprintf("/v3/organizations/{organization}/compute/instances/%d/elastic-ips", server)

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

func (s *serverAttachmentService) DetachElasticIp(ctx context.Context, server Id, elasticIp Id) (*Response, error) {
	p := fmt.Sprintf("/v3/organizations/{organization}/compute/instances/%d/elastic-ips/%d", server, elasticIp)

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
