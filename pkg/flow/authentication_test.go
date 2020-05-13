package flow

import (
	"context"
	"net/http"
	"testing"
)

func TestAuthenticationService_Login(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc("/v3/auth", func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodPost)
		assertPayload(t, req, &loginRequest{
			Username: "example@flow.swiss",
			Password: "SuperSecretPassword",
		})

		body := `{"admin":false,"token":"1279f86771c51bba551dde6f7c293853","id":1,"username":"example@flow.swiss","firstname":"Example","lastname":"User","phone_number":"","email_alternative":"","assigned_organizations":[{"created_at":"2020-05-07T17:00:57+00:00","id":1,"name":"Flow Swiss","address":"Somestreet 1","zip":"0000","city":"Somewhere","country":{"id":1,"name":"Switzerland","iso_alpha2":"CH","iso_alpha3":"CHE","calling_code":"+41"},"phone_number":"","status":{"id":3,"name":"Active","retention_time":null,"detailed":{"id":1,"name":"Not notified"}},"registered_modules":[{"id":2,"name":"Compute","parent":null,"sorting":1,"locations":[{"id":1,"name":"ALP1"},{"id":2,"name":"ZRH1"}]}],"contacts":{"primary":null,"billing":null,"technical":[]},"invoice_deployment_fees":true}],"default_organization":{"created_at":"2020-05-07T17:00:57+00:00","id":1,"name":"Flow Swiss","address":"Somestreet 1","zip":"0000","city":"Somewhere","country":{"id":1,"name":"Switzerland","iso_alpha2":"CH","iso_alpha3":"CHE","calling_code":"+41"},"phone_number":"","status":{"id":3,"name":"Active","retention_time":null,"detailed":{"id":1,"name":"Not notified"}},"registered_modules":[{"id":2,"name":"Compute","parent":null,"sorting":1,"locations":[{"id":1,"name":"ALP1"},{"id":2,"name":"ZRH1"}]}],"contacts":{"primary":null,"billing":null,"technical":[]},"invoice_deployment_fees":true}}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(body))
		if err != nil {
			t.Fatal(err)
		}
	})

	user, _, err := client.Authentication.Login(context.Background(), "example@flow.swiss", "SuperSecretPassword")
	if err != nil {
		t.Fatal(err)
	}

	if user.Token == "" {
		t.Error("service is unable to parse token")
	}

	if user.Username != "example@flow.swiss" {
		t.Error("expected username example@flow.swiss got", user.Username)
	}

	if user.DefaultOrganization.Name != "Flow Swiss" {
		t.Error("expected default organization Flow Swiss got", user.DefaultOrganization.Name)
	}
}

func TestAuthenticationService_Verify(t *testing.T) {
	setupMockServer(t)

	serveMux.HandleFunc("/v3/auth", func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodPost)
		assertPayload(t, req, &loginRequest{
			Username: "example@flow.swiss",
			Password: "SuperSecretPassword",
		})

		body := `{"token":"1279f86771c51bba551dde6f7c293853","two_factor":true}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(body))
		if err != nil {
			t.Fatal(err)
		}
	})

	serveMux.HandleFunc("/v3/2fa/verify", func(res http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, http.MethodPost)
		assertPayload(t, req, &verifyRequest{
			Token: "1279f86771c51bba551dde6f7c293853",
			Code:  "123456",
		})

		body := `{"admin":false,"token":"1279f86771c51bba551dde6f7c293853","id":1,"username":"example@flow.swiss","firstname":"Example","lastname":"User","phone_number":"","email_alternative":"","assigned_organizations":[{"created_at":"2020-05-07T17:00:57+00:00","id":1,"name":"Flow Swiss","address":"Somestreet 1","zip":"0000","city":"Somewhere","country":{"id":1,"name":"Switzerland","iso_alpha2":"CH","iso_alpha3":"CHE","calling_code":"+41"},"phone_number":"","status":{"id":3,"name":"Active","retention_time":null,"detailed":{"id":1,"name":"Not notified"}},"registered_modules":[{"id":2,"name":"Compute","parent":null,"sorting":1,"locations":[{"id":1,"name":"ALP1"},{"id":2,"name":"ZRH1"}]}],"contacts":{"primary":null,"billing":null,"technical":[]},"invoice_deployment_fees":true}],"default_organization":{"created_at":"2020-05-07T17:00:57+00:00","id":1,"name":"Flow Swiss","address":"Somestreet 1","zip":"0000","city":"Somewhere","country":{"id":1,"name":"Switzerland","iso_alpha2":"CH","iso_alpha3":"CHE","calling_code":"+41"},"phone_number":"","status":{"id":3,"name":"Active","retention_time":null,"detailed":{"id":1,"name":"Not notified"}},"registered_modules":[{"id":2,"name":"Compute","parent":null,"sorting":1,"locations":[{"id":1,"name":"ALP1"},{"id":2,"name":"ZRH1"}]}],"contacts":{"primary":null,"billing":null,"technical":[]},"invoice_deployment_fees":true}}`

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		_, err := res.Write([]byte(body))
		if err != nil {
			t.Fatal(err)
		}
	})

	user, _, err := client.Authentication.Login(context.Background(), "example@flow.swiss", "SuperSecretPassword")
	if err != nil {
		t.Fatal(err)
	}

	if !user.TwoFactor {
		t.Fatal("expected two factor to be enabled")
	}

	user, _, err = client.Authentication.Verify(context.Background(), user.Token, "123456")
	if err != nil {
		t.Fatal(err)
	}

	if user.Token == "" {
		t.Error("service is unable to parse token")
	}

	if user.Username != "example@flow.swiss" {
		t.Error("expected username example@flow.swiss got", user.Username)
	}

	if user.DefaultOrganization.Name != "Flow Swiss" {
		t.Error("expected default organization Flow Swiss got", user.DefaultOrganization.Name)
	}
}
