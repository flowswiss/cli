package flow

import (
	"context"
	"net/http"
	"path"
	"testing"
)

func TestKeyPairService_List(t *testing.T) {
	setupMockServer(t)

	options := PaginationOptions{NoFilter: 1}

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/key-pairs"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)
		assertPagination(t, req, options)

		response := `[{"id":1,"name":"Sample Key Pair","fingerprint":"3a:8c:ff:f8:db:c2:ab:7e:a4:1a:bc:fb:31:ec:21:5b"}]`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	keyPairs, _, err := client.KeyPair.List(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}

	if len(keyPairs) != 1 {
		t.Fatal("expected amount of key pairs to be 1, got", len(keyPairs))
	}

	key := keyPairs[0]

	if key.Id != 1 || key.Name != "Sample Key Pair" || key.Fingerprint != "3a:8c:ff:f8:db:c2:ab:7e:a4:1a:bc:fb:31:ec:21:5b" {
		t.Error("error while parsing key pair")
	}
}

func TestKeyPairService_Create(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/key-pairs"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodPost)
		assertPayload(t, req, &KeyPairCreate{
			Name:      "Sample Key Pair",
			PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCrDNXDLrYMpfrP6c5F8jAtoQaCC26qpgg4zM9dlFsWXZs6SoQ46efS1SaVlpHKklrp+semLjXzCcIy9a6mLzRNh7+0mxJUl456ZnB/gGtoi0tv5EAX9OhcPwHkMXT9JbaqtSIhmJ3obh0UN9lL0XobWIH9l0RkEocb4HbsF4hZouM6XPQjBjK1Uv4h6KvAdVQSXoh7CC3Ud2IuMWe9Q80bIWOknBQ7nbclsrf7Fn1efrRQt88LPXoZNbGWem2pzulVm200n5fdpXb79Ro7m2Ghcg8qqsy/FOnJeIrsRk+cXmgPwwc3o4gKkR1/p/gDkq8j2gCxLao7po3L1BtJOx/uL58KejpTVh6X3Z2ITn9VTE1UINtflU1BV2eD8oBzD+ILOYkTUSD0ufMrRTzKNlL7qRYP4/agBFuGIA4s06yNQlnIl3Nh5AYTneq9Vw2QDCmdW4HBU3clfjp5nYtjiKLvmAdZLxyicjWHrfub6+JbM3k2xdM68E7d1DudC7kytEc= someone@flow.swiss",
		})

		response := `{"id":1,"name":"Sample Key Pair","fingerprint":"3a:8c:ff:f8:db:c2:ab:7e:a4:1a:bc:fb:31:ec:21:5b"}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	data := &KeyPairCreate{
		Name:      "Sample Key Pair",
		PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCrDNXDLrYMpfrP6c5F8jAtoQaCC26qpgg4zM9dlFsWXZs6SoQ46efS1SaVlpHKklrp+semLjXzCcIy9a6mLzRNh7+0mxJUl456ZnB/gGtoi0tv5EAX9OhcPwHkMXT9JbaqtSIhmJ3obh0UN9lL0XobWIH9l0RkEocb4HbsF4hZouM6XPQjBjK1Uv4h6KvAdVQSXoh7CC3Ud2IuMWe9Q80bIWOknBQ7nbclsrf7Fn1efrRQt88LPXoZNbGWem2pzulVm200n5fdpXb79Ro7m2Ghcg8qqsy/FOnJeIrsRk+cXmgPwwc3o4gKkR1/p/gDkq8j2gCxLao7po3L1BtJOx/uL58KejpTVh6X3Z2ITn9VTE1UINtflU1BV2eD8oBzD+ILOYkTUSD0ufMrRTzKNlL7qRYP4/agBFuGIA4s06yNQlnIl3Nh5AYTneq9Vw2QDCmdW4HBU3clfjp5nYtjiKLvmAdZLxyicjWHrfub6+JbM3k2xdM68E7d1DudC7kytEc= someone@flow.swiss",
	}

	key, _, err := client.KeyPair.Create(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}

	if key.Name != data.Name || key.Fingerprint != "3a:8c:ff:f8:db:c2:ab:7e:a4:1a:bc:fb:31:ec:21:5b" {
		t.Error("error while parsing key pair")
	}
}

func TestKeyPairService_Delete(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc(path.Join("/v3/", organizationPath, "/compute/key-pairs/1"), func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodDelete)

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(204)
		_, err := res.Write([]byte{})
		if err != nil {
			t.Fatal(err)
		}
	})

	_, err := client.KeyPair.Delete(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
}
