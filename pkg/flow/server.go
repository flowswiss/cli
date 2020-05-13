package flow

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

type ServerService interface {
	List(ctx context.Context, options PaginationOptions) ([]*Server, *Response, error)
	Get(ctx context.Context, id Id) (*Server, *Response, error)
	Create(ctx context.Context, data *ServerCreate) (*Ordering, *Response, error)
	Update(ctx context.Context, id Id) (*Server, *Response, error)
	Delete(ctx context.Context, id Id) (*Response, error)

	RunAction(ctx context.Context, id Id, command string) (*Server, *Response, error)
}

type ServerAction struct {
	Id      Id     `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
	Sorting int    `json:"sorting"`
}

type ServerStatus struct {
	Id      Id             `json:"id"`
	Name    string         `json:"name"`
	Key     string         `json:"key"`
	Actions []ServerAction `json:"actions"`
}

type ServerNetworkAttachment struct {
	Network
	Interfaces []AttachedNetworkInterface `json:"network_interfaces"`
}

type Server struct {
	Id       Id                         `json:"id"`
	Name     string                     `json:"name"`
	Status   ServerStatus               `json:"status"`
	Image    Image                      `json:"image"`
	Product  Product                    `json:"product"`
	Location Location                   `json:"location"`
	Networks []*ServerNetworkAttachment `json:"networks"`
	KeyPair  KeyPair                    `json:"key_pair"`
}

type ServerCreate struct {
	Name             string `json:"name"`
	LocationId       Id     `json:"location_id"`
	ImageId          Id     `json:"image_id"`
	ProductId        Id     `json:"product_id"`
	AttachExternalIp bool   `json:"attach_external_ip"`
	NetworkId        Id     `json:"network_id"`
	PrivateIp        string `json:"private_ip,omitempty"`
	KeyPairId        Id     `json:"key_pair_id,omitempty"`
	Password         string `json:"password,omitempty"`
	CloudInit        string `json:"cloud_init,omitempty"`
}

func (s *Server) String() string {
	return s.Name
}

type serverService struct {
	client *Client
}

func (s *serverService) List(ctx context.Context, options PaginationOptions) ([]*Server, *Response, error) {
	p := path.Join("/v3/", s.client.OrganizationPath(), "/compute/instances")
	p, err := addOptions(p, options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, p, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	var val []*Server

	res, err := s.client.Do(req, &val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *serverService) Get(ctx context.Context, id Id) (*Server, *Response, error) {
	p := path.Join("/v3/", s.client.OrganizationPath(), fmt.Sprintf("/compute/instances/%d", id))

	req, err := s.client.NewRequest(ctx, http.MethodGet, p, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	val := &Server{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *serverService) Create(ctx context.Context, data *ServerCreate) (*Ordering, *Response, error) {
	p := path.Join("/v3/", s.client.OrganizationPath(), "/compute/instances")

	req, err := s.client.NewRequest(ctx, http.MethodPost, p, data, 0)
	if err != nil {
		return nil, nil, err
	}

	val := &Ordering{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *serverService) Update(ctx context.Context, id Id) (*Server, *Response, error) {
	return nil, nil, nil
}

func (s *serverService) Delete(ctx context.Context, id Id) (*Response, error) {
	p := path.Join("/v3/", s.client.OrganizationPath(), fmt.Sprintf("/compute/instances/%d", id))

	req, err := s.client.NewRequest(ctx, http.MethodDelete, p, nil, 0)
	if err != nil {
		return nil, err
	}

	val := &Server{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serverService) RunAction(ctx context.Context, id Id, command string) (*Server, *Response, error) {
	p := path.Join("/v3/", s.client.OrganizationPath(), fmt.Sprintf("/compute/instances/%d/action", id))

	data := struct {
		Action string `json:"action"`
	}{command}

	req, err := s.client.NewRequest(ctx, http.MethodPost, p, &data, 0)
	if err != nil {
		return nil, nil, err
	}

	val := &Server{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}
