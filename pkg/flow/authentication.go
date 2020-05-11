package flow

import (
	"context"
	"net/http"
)

type AuthenticationService interface {
	Login(ctx context.Context, username, password string) (*AuthenticatedUser, *Response, error)
	Verify(ctx context.Context, token, code string) (*AuthenticatedUser, *Response, error)
}

type authenticationService struct {
	client *Client
}

type AuthenticatedUser struct {
	Token     string `json:"token"`
	TwoFactor bool   `json:"two_factor"`
	User
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type verifyRequest struct {
	Token string `json:"token"`
	Code  string `json:"code"`
}

func (s *authenticationService) Login(ctx context.Context, username, password string) (*AuthenticatedUser, *Response, error) {
	path := "/v3/auth"

	body := loginRequest{
		Username: username,
		Password: password,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, body, FlagNoAuthentication)
	if err != nil {
		return nil, nil, err
	}

	user := &AuthenticatedUser{}

	res, err := s.client.Do(req, user)
	if err != nil {
		return nil, nil, err
	}

	return user, res, nil
}

func (s *authenticationService) Verify(ctx context.Context, token, code string) (*AuthenticatedUser, *Response, error) {
	path := "/v3/2fa/verify"

	body := verifyRequest{
		Token: token,
		Code:  code,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, body, FlagNoAuthentication)
	if err != nil {
		return nil, nil, err
	}

	user := &AuthenticatedUser{}

	res, err := s.client.Do(req, user)
	if err != nil {
		return nil, nil, err
	}

	return user, res, nil
}
