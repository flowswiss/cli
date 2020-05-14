package flow

import (
	"context"
	"fmt"
	"net/http"
)

type KeyPairService interface {
	List(ctx context.Context, options PaginationOptions) ([]*KeyPair, *Response, error)
	Create(ctx context.Context, data *KeyPairCreate) (*KeyPair, *Response, error)
	Delete(ctx context.Context, id Id) (*Response, error)
}

type KeyPair struct {
	Id          Id     `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

func (k *KeyPair) String() string {
	return k.Name
}

type KeyPairCreate struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

type keyPairService struct {
	client *Client
}

func (s *keyPairService) List(ctx context.Context, options PaginationOptions) ([]*KeyPair, *Response, error) {
	p := "/v3/organizations/{organization}/compute/key-pairs"
	p, err := addOptions(p, options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, p, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	var val []*KeyPair

	res, err := s.client.Do(req, &val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *keyPairService) Create(ctx context.Context, data *KeyPairCreate) (*KeyPair, *Response, error) {
	p := "/v3/organizations/{organization}/compute/key-pairs"

	req, err := s.client.NewRequest(ctx, http.MethodPost, p, data, 0)
	if err != nil {
		return nil, nil, err
	}

	val := &KeyPair{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *keyPairService) Delete(ctx context.Context, id Id) (*Response, error) {
	p := fmt.Sprintf("/v3/organizations/{organization}/compute/key-pairs/%d", id)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, p, nil, 0)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
