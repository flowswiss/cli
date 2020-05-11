package flow

type NetworkInterface struct {
	Id        Id     `json:"id"`
	PrivateIp string `json:"private_ip"`
	PublicIp  string `json:"public_ip"`
}
