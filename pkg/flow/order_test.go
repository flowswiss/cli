package flow

import (
	"context"
	"net/http"
	"path"
	"testing"
)

func TestOrderService_Get(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/orders/1"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)

		response := `{"id":1,"status":{"id":1,"name":"created"},"product_instance":null}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	order, _, err := client.Order.Get(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if order.Id != 1 {
		t.Error("expected id to be 1, got", order.Id)
	}

	if order.Status.Id != 1 {
		t.Error("expected status to be 1, got", order.Status.Id)
	}
}
