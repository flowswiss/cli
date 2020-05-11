package commands

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

var authConfig = viper.New()

var (
	authCommand = &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication settings",
	}

	loginCommand = &cobra.Command{
		Use:   "login",
		Short: "Login using username and password",
		Long:  "Checks provided login credentials and stores them in the $HOME/.flow/credentials.json file",
		RunE:  authenticate,
	}
)

func init() {
	authCommand.AddCommand(loginCommand)
}

func readAuthConfig() {
	configureConfig("credentials", authConfig)
	authConfig.SetConfigPermissions(0600)

	authConfig.SetDefault("username", viper.GetString("username"))
	authConfig.SetDefault("password", viper.GetString("password"))
	_ = authConfig.ReadInConfig()
}

func authenticate(cmd *cobra.Command, args []string) error {
	if client.CredentialsProvider.Username() == "" || client.CredentialsProvider.Password() == "" {
		return fmt.Errorf("please provide a username and password")
	}

	_, _, err := client.Authentication.Login(context.Background(), client.CredentialsProvider.Username(), client.CredentialsProvider.Password())
	if err != nil {
		return err
	}

	return authConfig.WriteConfigAs(filepath.Join(configDir, "credentials."+configType))
}
