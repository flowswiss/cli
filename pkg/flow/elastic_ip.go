package flow

import "context"

type ElasticIpService interface {
	List(ctx context.Context, options PaginationOptions) (*ElasticIp, *Response, error)
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
