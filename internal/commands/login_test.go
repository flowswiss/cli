package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/flowswiss/cli/pkg/flow"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type mockAuthenticationService struct {
}

func (m *mockAuthenticationService) Login(ctx context.Context, username, password string) (*flow.AuthenticatedUser, *flow.Response, error) {
	return &flow.AuthenticatedUser{
		Token:     "testtoken",
		TwoFactor: false,
		User:      flow.User{},
	}, nil, nil
}

func (m *mockAuthenticationService) Verify(ctx context.Context, token, code string) (*flow.AuthenticatedUser, *flow.Response, error) {
	return m.Login(ctx, "", "")
}

func Test_Authenticate(t *testing.T) {
	setupTests()

	configFile := filepath.Join(configDir, "credentials."+configType)
	_ = os.Remove(configFile)

	readAuthConfig()
	authConfig.Set("username", "example@flow.swiss")
	authConfig.Set("password", "SuperSecretPassword")

	client.Authentication = &mockAuthenticationService{}

	err := authenticate(loginCommand, []string{})
	if err != nil {
		t.Fatal(err)
	}

	stat, err := os.Stat(configFile)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Mode().Perm() != 0600 {
		t.Error("wrong permissions of credentials file")
	}

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	if err = json.Compact(buf, content); err != nil {
		t.Fatal(err)
	}

	expectedContent := "{\"password\":\"SuperSecretPassword\",\"username\":\"example@flow.swiss\"}"

	if buf.String() != expectedContent {
		t.Error("invalid content in file", buf.String())
	}
}
