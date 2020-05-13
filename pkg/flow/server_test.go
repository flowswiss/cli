package flow

import (
	"context"
	"net/http"
	"path"
	"testing"
)

func TestServerService_List(t *testing.T) {
	setupMockServer(t)

	options := PaginationOptions{NoFilter: 1}

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/instances"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)
		assertPagination(t, req, options)

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte("[]"))
		if err != nil {
			t.Fatal(err)
		}
	})

	_, _, err := client.Server.List(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerService_Create(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/instances"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodPost)
		assertPayload(t, req, &ServerCreate{
			Name:             "Test Server",
			LocationId:       1,
			ImageId:          3,
			ProductId:        40,
			AttachExternalIp: true,
			NetworkId:        1,
			PrivateIp:        "",
			KeyPairId:        1,
			Password:         "",
			CloudInit:        "",
		})

		response := `{"ref": "https://api.flow.swiss/v3/organizations/1/orders/1"}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	data := &ServerCreate{
		Name:             "Test Server",
		LocationId:       1,
		ImageId:          3,
		ProductId:        40,
		AttachExternalIp: true,
		NetworkId:        1,
		KeyPairId:        1,
	}

	ordering, _, err := client.Server.Create(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}

	id, err := ordering.Id()
	if err != nil {
		t.Fatal(err)
	}

	if id != 1 {
		t.Error("expected order id 1, got", id)
	}
}
