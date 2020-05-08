package flow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
)

const (
	VersionMajor = 1
	VersionMinor = 0
	VersionPatch = 0

	FlagNoAuthentication = 1

	ErrorMissingCredentials = ClientError("missing credentials provider for authenticated request")
)

type ClientFlag uint
type ClientError string
type ClientRequestCallback func(req *http.Request)
type ClientResponseCallback func(res *http.Response)

type Client struct {
	Base      *url.URL
	UserAgent string

	Client *http.Client

	CredentialsProvider CredentialsProvider
	TokenStorage        TokenStorage

	OnRequest  ClientRequestCallback
	OnResponse ClientResponseCallback

	Authentication AuthenticationService
}

type Response struct {
	*http.Response
	Pagination
}

type ErrorResponse struct {
	Response  *http.Response
	Message   string
	RequestID string
}

func (e ClientError) Error() string {
	return string(e)
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s %s (request %s) resulted in %s: %s", e.Response.Request.Method, e.Response.Request.URL, e.RequestID, e.Response.Status, e.Message)
}

func NewClient(base *url.URL) *Client {
	client := &Client{
		Base:      base,
		UserAgent: fmt.Sprintf("flow/%d.%d.%d", VersionMajor, VersionMinor, VersionPatch),
		Client:    &http.Client{},
	}

	client.Authentication = NewAuthenticationService(client)

	return client
}

func (c *Client) AddUserAgent(userAgent string) {
	c.UserAgent = fmt.Sprintf("%s %s", userAgent, c.UserAgent)
}

func (c *Client) AuthenticationToken(ctx context.Context) (string, error) {
	if c.TokenStorage != nil && c.TokenStorage.IsValid() {
		return c.TokenStorage.Token(), nil
	}

	if c.CredentialsProvider == nil {
		return "", ErrorMissingCredentials
	}

	user, _, err := c.Authentication.Login(ctx, c.CredentialsProvider.Username(), c.CredentialsProvider.Password())
	if err != nil {
		return "", err
	}

	if user.TwoFactor {
		user, _, err = c.Authentication.Verify(ctx, user.Token, c.CredentialsProvider.TwoFactorCode())
		if err != nil {
			return "", err
		}
	}

	if c.TokenStorage != nil {
		c.TokenStorage.SetToken(user.Token)
	}

	return user.Token, nil
}

func (c *Client) NewRequest(ctx context.Context, method string, path string, body interface{}, flags ClientFlag) (*http.Request, error) {
	u, err := c.Base.Parse(path)
	if err != nil {
		return nil, err
	}

	reader, ok := body.(io.Reader)
	if !ok {
		buf := &bytes.Buffer{}
		if body != nil {
			err := json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}
		reader = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.UserAgent)

	if flags&FlagNoAuthentication == 0 {
		token, err := c.AuthenticationToken(ctx)
		if err != nil {
			return nil, err
		}
		req.Header.Add("X-Auth-Token", token)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, val interface{}) (*Response, error) {
	if c.OnRequest != nil {
		c.OnRequest(req)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if c.OnResponse != nil {
		c.OnResponse(res)
	}

	if res.StatusCode >= 400 {
		apiError := &struct {
			Error struct {
				Message struct {
					En string `json:"en"`
				} `json:"message"`
			} `json:"error"`
		}{}

		err := json.NewDecoder(res.Body).Decode(apiError)
		if err != nil {
			return nil, err
		}

		return nil, &ErrorResponse{
			Response:  res,
			Message:   apiError.Error.Message.En,
			RequestID: res.Header.Get("X-Request-Id"),
		}
	}

	if writer, ok := val.(io.Writer); ok {
		_, err := io.Copy(writer, res.Body)
		if err != nil {
			return nil, err
		}
	} else if val != nil {
		err := json.NewDecoder(res.Body).Decode(val)
		if err != nil {
			return nil, err
		}
	}

	return &Response{
		Response:   res,
		Pagination: parsePagination(res),
	}, nil
}

func addOptions(path string, options interface{}) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	newQuery, err := query.Values(options)
	if err != nil {
		return "", err
	}

	prevQuery := u.Query()
	for key, arr := range newQuery {
		for _, val := range arr {
			prevQuery.Add(key, val)
		}
	}
	u.RawQuery = prevQuery.Encode()

	return u.String(), nil
}
