package commands

import (
	"net/http"
	"net/http/httputil"

	"github.com/flowswiss/cli/v2/pkg/console"
)

var _ http.RoundTripper = (*dryRunTransport)(nil)

type dryRunTransport struct {
	delegate http.RoundTripper
}

func (d dryRunTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodGet {
		return d.base().RoundTrip(req)
	}

	Stderr.Color(console.Bright + console.Black)
	_ = req.Write(Stderr)
	Stderr.Reset()

	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
	}, nil
}

func (d dryRunTransport) base() http.RoundTripper {
	if d.delegate == nil {
		return http.DefaultTransport
	}

	return d.delegate
}

var _ http.RoundTripper = (*logRequestTransport)(nil)

type logRequestTransport struct {
	delegate http.RoundTripper
}

func (l logRequestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := l.base().RoundTrip(req)

	if err == nil {
		Stderr.Color(console.Bright+console.Black).Printf("request to `%s %s` resulted in `%s`\n", req.Method, req.URL, res.Status).Reset()
	} else {
		Stderr.Errorf("request to `%s %s` resulted in error `%v`\n", req.Method, req.URL, err)
	}

	return res, err
}

func (l logRequestTransport) base() http.RoundTripper {
	if l.delegate == nil {
		return http.DefaultTransport
	}

	return l.delegate
}

var _ http.RoundTripper = (*dumpRequestTransport)(nil)

type dumpRequestTransport struct {
	delegate http.RoundTripper
}

func (d dumpRequestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	Stderr.Color(console.Bright + console.Black)
	defer Stderr.Reset()

	// dump request
	data, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	Stderr.Println(string(data))

	// make request
	res, err := d.base().RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// dump response
	data, err = httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}
	Stderr.Println(string(data))

	return res, nil
}

func (d dumpRequestTransport) base() http.RoundTripper {
	if d.delegate == nil {
		return http.DefaultTransport
	}

	return d.delegate
}
