package flow

import (
	"context"
	"net/http"
	"path"
	"reflect"
	"testing"
)

var networkAlp1 = &Network{
	Id:   1,
	Name: "Default Network",
	Cidr: "172.31.0.0/20",
	Location: Location{
		Id:   1,
		Name: "ALP1",
	},
	UsedIps:  5,
	TotalIps: 3995,
}

var networkZrh1 = &Network{
	Id:   2,
	Name: "Default Network",
	Cidr: "172.31.16.0/20",
	Location: Location{
		Id:   2,
		Name: "ZRH1",
	},
	UsedIps:  3,
	TotalIps: 3995,
}

func TestNetworkService_List(t *testing.T) {
	setupMockServer(t)

	options := PaginationOptions{NoFilter: 1}

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/networks"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)
		assertPagination(t, req, options)

		response := `[{"id":1,"name":"Default Network","cidr":"172.31.0.0/20","location":{"id":1,"name":"ALP1"},"used_ips":5,"total_ips":3995},{"id":2,"name":"Default Network","cidr":"172.31.16.0/20","location":{"id":2,"name":"ZRH1"},"used_ips":3,"total_ips":3995}]`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	networks, _, err := client.Network.List(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}

	if len(networks) != 2 {
		t.Fatal("expected amount of key pairs to be 2, got", len(networks))
	}

	if !reflect.DeepEqual(networkAlp1, networks[0]) {
		t.Errorf("expected %v, got %v", networkAlp1, networks[0])
	}

	if !reflect.DeepEqual(networkZrh1, networks[1]) {
		t.Errorf("expected %v, got %v", networkZrh1, networks[1])
	}
}

func TestNetworkService_Get(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/networks/1"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)

		response := `{"id":1,"name":"Default Network","description":"Initially created network","cidr":"172.31.0.0\/20","location":{"id":1,"name":"ALP1"},"domain_name_servers":["1.1.1.1","8.8.8.8"],"allocation_pool_start":"172.31.0.100","allocation_pool_end":"172.31.15.254","gateway_ip":"172.31.0.1","used_ips":6,"total_ips":3995}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	network, _, err := client.Network.Get(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if network.Id != 1 || network.Name != "Default Network" || network.Cidr != "172.31.0.0/20" {
		t.Error("error while parsing network")
	}

	if network.AllocationPoolStart != "172.31.0.100" || network.AllocationPoolEnd != "172.31.15.254" {
		t.Error("error while parsing allocation pool")
	}

	if len(network.DomainNameServers) != 2 {
		t.Error("error while parsing domain name server")
	}
}
