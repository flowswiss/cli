package flow

import (
	"context"
	"net/http"
	"testing"
)

func TestProductService_List(t *testing.T) {
	setupMockServer(t)

	options := PaginationOptions{Page: 1, PerPage: 3}

	serveMux.HandleFunc("/v3/products", func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)
		assertPagination(t, req, options)

		response := `[{"id":1,"product_name":"Elastic IP Free","type":{"id":5,"name":"Elastic IP","key":"compute-engine-elastic-ip"},"visibility":"public","usage_cycle":{"id":2,"name":"Hour","duration":1},"items":[{"id":4,"name":"IPv4 Adresse","description":"IPv4 Adresse","amount":1}],"price":0,"availability":[{"location":{"id":1,"name":"ALP1"},"available":-1},{"location":{"id":2,"name":"ZRH1"},"available":-1}],"category":null,"deployment_fees":[]},{"id":8,"product_name":"Elastic IP","type":{"id":5,"name":"Elastic IP","key":"compute-engine-elastic-ip"},"visibility":"public","usage_cycle":{"id":2,"name":"Hour","duration":1},"items":[{"id":4,"name":"IPv4 Adresse","description":"IPv4 Adresse","amount":1}],"price":5.037,"availability":[{"location":{"id":1,"name":"ALP1"},"available":-1},{"location":{"id":2,"name":"ZRH1"},"available":-1}],"category":null,"deployment_fees":[]},{"id":10,"product_name":"Windows Server","type":{"id":2,"name":"License","key":"license"},"visibility":"public","usage_cycle":{"id":3,"name":"Monthly","duration":730},"items":[{"id":10,"name":"Windows Server 2016 Standard","description":"Windows Server 2016 Standard License","amount":1}],"price":10,"availability":[{"location":{"id":1,"name":"ALP1"},"available":-1},{"location":{"id":2,"name":"ZRH1"},"available":-1}],"category":null,"deployment_fees":[]}]`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	products, _, err := client.Product.List(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}

	if len(products) != 3 {
		t.Fatal("expected amount of images to be 3, got", len(products))
	}
}

func TestProductService_Get(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc("/v3/products/10", func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)

		response := `{"id":10,"product_name":"Windows Server","type":{"id":2,"name":"License","key":"license"},"visibility":"public","usage_cycle":{"id":3,"name":"Monthly","duration":730},"items":[{"id":10,"name":"Windows Server 2016 Standard","description":"Windows Server 2016 Standard License","amount":1}],"price":10,"availability":[{"location":{"id":1,"name":"ALP1"},"available":-1},{"location":{"id":2,"name":"ZRH1"},"available":-1}],"category":null,"deployment_fees":[]}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	product, _, err := client.Product.Get(context.Background(), 10)
	if err != nil {
		t.Error(err)
	}

	if product.Id != 10 || product.Name != "Windows Server" || product.Price != 10 {
		t.Error("error while parsing product")
	}

	if len(product.Availability) != 2 {
		t.Error("error while parsing product availability")
	}

	if product.UsageCycle.Name != "Monthly" {
		t.Error("error while parsing product usage period")
	}
}
