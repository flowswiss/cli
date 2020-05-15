package flow

import (
	"context"
	"fmt"
	"net/http"
)

type ProductService interface {
	List(ctx context.Context, options PaginationOptions) ([]*Product, *Response, error)
	ListByType(ctx context.Context, options PaginationOptions, t string) ([]*Product, *Response, error)
	Get(ctx context.Context, id Id) (*Product, *Response, error)
}

type ProductType struct {
	Id   Id     `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type ProductUsageCycle struct {
	Id       Id     `json:"id"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

type ProductItem struct {
	Id          Id     `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
}

type ProductAvailability struct {
	Location  Location `json:"location"`
	Available int      `json:"available"`
}

type DeploymentFee struct {
	Location        Location `json:"location"`
	Price           float64  `json:"price"`
	FreeDeployments int      `json:"free_deployments"`
}

type Product struct {
	Id             Id                    `json:"id"`
	Name           string                `json:"product_name"`
	Type           ProductType           `json:"type"`
	Visibility     string                `json:"visibility"`
	UsageCycle     ProductUsageCycle     `json:"usage_cycle"`
	Items          []*ProductItem        `json:"items"`
	Price          float64               `json:"price"`
	Availability   []ProductAvailability `json:"availability"`
	Category       string                `json:"category"`
	DeploymentFees []*DeploymentFee      `json:"deployment_fees"`
}

func (p *Product) String() string {
	return p.Name
}

func (p *Product) AvailableAt(location *Location) bool {
	for _, availability := range p.Availability {
		if availability.Location.Id == location.Id {
			return availability.Available == -1 || availability.Available > 0
		}
	}
	return false
}

func (p *Product) FindItem(id Id) *ProductItem {
	for _, item := range p.Items {
		if item.Id == id {
			return item
		}
	}
	return nil
}

type productService struct {
	client *Client
}

func (s *productService) List(ctx context.Context, options PaginationOptions) ([]*Product, *Response, error) {
	path, err := addOptions("/v3/products", options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil, FlagNoAuthentication)
	if err != nil {
		return nil, nil, err
	}

	var val []*Product

	res, err := s.client.Do(req, &val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *productService) ListByType(ctx context.Context, options PaginationOptions, t string) ([]*Product, *Response, error) {
	path := fmt.Sprintf("/v3/products/%s", t)
	path, err := addOptions(path, options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil, FlagNoAuthentication)
	if err != nil {
		return nil, nil, err
	}

	var val []*Product

	res, err := s.client.Do(req, &val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *productService) Get(ctx context.Context, id Id) (*Product, *Response, error) {
	path := fmt.Sprintf("/v3/products/%d", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil, FlagNoAuthentication)
	if err != nil {
		return nil, nil, err
	}

	val := &Product{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}
