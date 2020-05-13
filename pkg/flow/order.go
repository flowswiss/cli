package flow

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type OrderService interface {
	Get(ctx context.Context, id Id) (*Order, *Response, error)
}

type OrderStatus struct {
	Id   Id     `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	Id     Id          `json:"id"`
	Status OrderStatus `json:"status"`
}

type Ordering struct {
	Ref string `json:"ref"`
}

func (o *Ordering) Id() (Id, error) {
	regex, err := regexp.Compile("/organizations/\\d+/orders/(\\d+)$")
	if err != nil {
		return 0, err
	}

	data := regex.FindStringSubmatch(o.Ref)

	id, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return Id(id), nil
}

type orderService struct {
	client *Client
}

func (s *orderService) Get(ctx context.Context, id Id) (*Order, *Response, error) {
	p := fmt.Sprintf("/v3/organizations/{organization}/orders/%d", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, p, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	val := &Order{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}
