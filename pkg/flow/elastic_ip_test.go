package flow

import (
	"context"
	"net/http"
	"path"
	"reflect"
	"testing"
)

var testElasticIp = &ElasticIp{
	Id: 1,
	Product: ElasticIpProduct{
		Id:   8,
		Name: "Elastic IP",
		Type: "public",
	},
	Location: Location{
		Id:   1,
		Name: "ALP1",
	},
	Price:            5.037,
	PublicIp:         "1.1.1.1",
	PrivateIp:        "",
	AttachedInstance: nil,
}

func TestElasticIpService_List(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/elastic-ips"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)
		assertPagination(t, req, PaginationOptions{NoFilter: 1})

		body := `[{"id":1,"product":{"id":8,"name":"Elastic IP","type":"public"},"location":{"id":1,"name":"ALP1"},"price":5.037,"public_ip":"1.1.1.1","private_ip":null,"attached_instance":null}]`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(body))
		if err != nil {
			t.Fatal(err)
		}
	})

	elasticIps, _, err := client.ElasticIp.List(context.Background(), PaginationOptions{NoFilter: 1})
	if err != nil {
		t.Fatal(err)
	}

	if len(elasticIps) != 1 {
		t.Errorf("expected 1 elastic ip, got %d", len(elasticIps))
	}

	if !reflect.DeepEqual(testElasticIp, elasticIps[0]) {
		t.Errorf("expected %v, got %v", testElasticIp, elasticIps[0])
	}
}

func TestElasticIpService_Create(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/elastic-ips"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodPost)
		assertPayload(t, req, &ElasticIpCreate{
			LocationId: 1,
		})

		body := `{"id":1,"product":{"id":8,"name":"Elastic IP","type":"public"},"location":{"id":1,"name":"ALP1"},"price":5.037,"public_ip":"1.1.1.1","private_ip":null,"attached_instance":null}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(body))
		if err != nil {
			t.Fatal(err)
		}
	})

	elasticIp, _, err := client.ElasticIp.Create(context.Background(), &ElasticIpCreate{LocationId: 1})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testElasticIp, elasticIp) {
		t.Errorf("expected %v, got %v", testElasticIp, elasticIp)
	}
}

func TestElasticIpService_Delete(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/elastic-ips/1"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodDelete)

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(204)
	})

	_, err := client.ElasticIp.Delete(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
}
