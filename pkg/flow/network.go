package flow

type Network struct {
	Id       Id       `json:"id"`
	Name     string   `json:"name"`
	Cidr     string   `json:"cidr"`
	Location Location `json:"location"`
	UsedIps  int      `json:"used_ips"`
	TotalIps int      `json:"total_ips"`
}
