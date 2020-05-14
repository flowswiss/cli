package flow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

var (
	server   *httptest.Server
	serveMux *http.ServeMux

	client           *Client
	organizationPath string
)

func setupMockServer(t *testing.T) {
	var err error

	serveMux = http.NewServeMux()

	server = httptest.NewServer(serveMux)
	base, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	client = NewClient(base)
	client.Flags |= FlagNoAuthentication
	client.SelectedOrganization = 1

	organizationPath = fmt.Sprintf("organizations/%d", client.SelectedOrganization)
}

func assertPagination(t *testing.T, req *http.Request, expectation PaginationOptions) {
	values := req.URL.Query()

	page, _ := strconv.ParseInt(values.Get("page"), 10, 32)
	if int(page) != expectation.Page {
		t.Error(fmt.Sprintf("expected page to be %d, got %d", expectation.Page, page))
	}

	perPage, _ := strconv.ParseInt(values.Get("per_page"), 10, 32)
	if int(perPage) != expectation.PerPage {
		t.Error(fmt.Sprintf("expected per_page to be %d, got %d", expectation.PerPage, perPage))
	}

	noFilter, _ := strconv.ParseInt(values.Get("no_filter"), 10, 32)
	if int(noFilter) != expectation.NoFilter {
		t.Error(fmt.Sprintf("expected no_filter to be %d, got %d", expectation.NoFilter, noFilter))
	}
}

func assertMethod(t *testing.T, req *http.Request, method string) {
	if req.Method != method {
		t.Errorf("expected method %s, got %s", method, req.Method)
	}
}

func assertPayload(t *testing.T, req *http.Request, expected interface{}) {
	parsed := reflect.New(reflect.TypeOf(expected).Elem()).Interface()

	err := json.NewDecoder(req.Body).Decode(parsed)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, parsed) {
		t.Errorf("expected %v, got %v", expected, parsed)
	}

	_ = req.Body.Close()
}
