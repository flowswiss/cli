package flow

import "context"

type NetworkInterfaceService interface {
	List(ctx context.Context, options PaginationOptions) ([]*NetworkInterface, *Response, error)
	Create(ctx context.Context, data *NetworkInterfaceCreate) (*NetworkInterface, *Response, error)
	Delete(ctx context.Context, id Id) (*Response, error)
}

type NetworkInterface struct {
	Id                Id               `json:"id"`
	PrivateIp         string           `json:"private_ip"`
	MacAddress        string           `json:"mac_address"`
	Network           *Network         `json:"network"`
	AttachedElasticIp *ElasticIp       `json:"attached_elastic_ip"`
	SecurityGroups    []*SecurityGroup `json:"security_groups"`
	Security          bool             `json:"security"`
}

type AttachedNetworkInterface struct {
	Id        Id     `json:"id"`
	PrivateIp string `json:"private_ip"`
	PublicIp  string `json:"public_ip"`
}

type NetworkInterfaceCreate struct {
	NetworkId Id     `json:"network_id"`
	PrivateIp string `json:"private_ip"`
}
