package flow

import (
	"context"
	"net/http"
	"testing"
)

func TestLocationService_List(t *testing.T) {
	setupMockServer(t)

	options := PaginationOptions{NoFilter: 1}

	serveMux.HandleFunc("/v3/entities/locations", func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)
		assertPagination(t, req, options)

		response := `[{"id":1,"name":"ALP1","key":"key-alp1","city":"Lucerne","available_modules":[{"id":2,"name":"Compute"},{"id":4,"name":"Object Storage"},{"id":5,"name":"Compute Networking"}]},{"id":2,"name":"ZRH1","key":"key-zrh1","city":"Zurich","available_modules":[{"id":2,"name":"Compute"},{"id":3,"name":"Mac Bare Metal"},{"id":4,"name":"Object Storage"},{"id":5,"name":"Compute Networking"},{"id":6,"name":"Mac Bare Metal - Experimental"}]}]`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	locations, _, err := client.Location.List(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}

	if len(locations) != 2 {
		t.Fatal("expected amount of locations to be 2, got", len(locations))
	}

	alp1 := locations[0]
	zrh1 := locations[1]

	if alp1.Id != 1 || alp1.Name != "ALP1" {
		t.Error("error while parsing location alp1")
	}

	if zrh1.Id != 2 || zrh1.Name != "ZRH1" {
		t.Error("error while parsing location zrh1")
	}
}

func TestLocationService_Get(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc("/v3/entities/locations/1", func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)

		response := `{"id":1,"name":"ALP1","key":"key-alp1","city":"Lucerne","available_modules":[{"id":2,"name":"Compute"},{"id":4,"name":"Object Storage"},{"id":5,"name":"Compute Networking"}]}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	location, _, err := client.Location.Get(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}

	if location.Id != 1 || location.Name != "ALP1" {
		t.Error("error while parsing location")
	}
}
