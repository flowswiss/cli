package flow

import (
	"context"
	"net/http"
	"path"
	"testing"
)

func TestServerService_Create(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc(path.Join("/v3/", client.OrganizationPath(), "/compute/instances"), func(res http.ResponseWriter, req *http.Request) {
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
