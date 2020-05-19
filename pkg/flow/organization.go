package flow

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type OrganizationService interface {
	List(ctx context.Context, options PaginationOptions) ([]*Organization, *Response, error)
	Get(ctx context.Context, id Id) (*Organization, *Response, error)
	GetCurrent(ctx context.Context) (*Organization, *Response, error)
}

type Country struct {
	Id          Id     `json:"id"`
	Name        string `json:"name"`
	IsoAlpha2   string `json:"iso_alpha_2"`
	IsoAlpha3   string `json:"iso_alpha_3"`
	CallingCode string `json:"calling_code"`
}

type Organization struct {
	Id                    Id        `json:"id"`
	Name                  string    `json:"name"`
	Address               string    `json:"address"`
	Zip                   string    `json:"zip"`
	City                  string    `json:"city"`
	PhoneNumber           string    `json:"phone_number"`
	InvoiceDeploymentFees bool      `json:"invoice_deployment_fees"`
	CreatedAt             time.Time `json:"created_at"`

	Status struct {
		Id            Id         `json:"id"`
		Name          string     `json:"name"`
		RetentionTime *time.Time `json:"retention_time"`
	} `json:"status"`

	RegisteredModules []*Module `json:"registered_modules"`

	Contacts struct {
		Primary   *User  `json:"primary"`
		Billing   *User  `json:"billing"`
		Technical []User `json:"technical"`
	} `json:"contacts"`
}

type organizationService struct {
	client *Client
}

func (s *organizationService) List(ctx context.Context, options PaginationOptions) ([]*Organization, *Response, error) {
	p, err := addOptions("/v3/organizations", options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, p, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	var val []*Organization

	res, err := s.client.Do(req, &val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, err
}

func (s *organizationService) Get(ctx context.Context, id Id) (*Organization, *Response, error) {
	p := fmt.Sprintf("/v3/organizations/%d", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, p, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	val := &Organization{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, err
}

func (s *organizationService) GetCurrent(ctx context.Context) (*Organization, *Response, error) {
	p := "/v3/organizations/{organization}"

	req, err := s.client.NewRequest(ctx, http.MethodGet, p, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	val := &Organization{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, err
}
