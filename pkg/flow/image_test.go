package flow

import (
	"context"
	"net/http"
	"testing"
)

func TestImageService_List(t *testing.T) {
	setupMockServer(t)

	options := PaginationOptions{Page: 1, PerPage: 3}

	serveMux.HandleFunc("/v3/entities/images", func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)
		assertPagination(t, req, options)

		response := `[{"id":1,"os":"Ubuntu","version":"18.04 LTS","key":"linux-ubuntu-18.04-lts","category":"Linux","type":"distribution","min_root_disk_size":10,"sorting":2,"required_licenses":[],"available_locations":[1,2]},{"id":2,"os":"Ubuntu","version":"16.04 LTS","key":"linux-ubuntu-16.04-lts","category":"Linux","type":"distribution","min_root_disk_size":10,"sorting":3,"required_licenses":[],"available_locations":[1,2]},{"id":9,"os":"Windows Server","version":"2019 Standard","key":"microsoft-windows-server-2019","category":"Windows","type":"distribution","min_root_disk_size":40,"sorting":9,"required_licenses":[{"id":10,"product_name":"Windows Server","type":{"id":2,"name":"License","key":"license"},"visibility":"public","usage_cycle":{"id":3,"name":"Monthly","duration":730},"items":[{"id":10,"name":"Windows Server 2016 Standard","description":"Windows Server 2016 Standard License","amount":1}],"price":10,"availability":[{"location":{"id":1,"name":"ALP1"},"available":-1},{"location":{"id":2,"name":"ZRH1"},"available":-1}],"category":null,"deployment_fees":[]}],"available_locations":[1,2]}]`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	images, _, err := client.Image.List(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}

	if len(images) != 3 {
		t.Fatal("expected amount of images to be 3, got", len(images))
	}

	ubuntu := images[0]
	windows := images[2]

	if ubuntu.OperatingSystem != "Ubuntu" || ubuntu.Version != "18.04 LTS" {
		t.Error("error while parsing ubuntu image")
	}

	if windows.OperatingSystem != "Windows Server" || windows.Version != "2019 Standard" || windows.MinRootDiskSize != 40 {
		t.Error("error while parsing windows image")
	}

	if len(windows.RequiredLicenses) != 1 {
		t.Fatal("windows should require window server license")
	}

	license := windows.RequiredLicenses[0]

	if license.UsageCycle.Name != "Monthly" || license.Price != 10 {
		t.Error("error while paring windows license")
	}
}

func TestImageService_Get(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc("/v3/entities/images/1", func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodGet)

		response := `{"id":1,"os":"Ubuntu","version":"18.04 LTS","key":"linux-ubuntu-18.04-lts","category":"Linux","type":"distribution","min_root_disk_size":10,"sorting":2,"required_licenses":[],"available_locations":[1,2]}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	})

	image, _, err := client.Image.Get(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}

	if image.Id != 1 || image.OperatingSystem != "Ubuntu" || image.Version != "18.04 LTS" {
		t.Error("error while parsing image")
	}

	if len(image.AvailableLocations) != 2 {
		t.Error("error while parsing location availability")
	}
}
