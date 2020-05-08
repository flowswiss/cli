package commands

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

var loginCommand = &cobra.Command{
	Use:   "login",
	Short: "Login using username and password",
	Long:  "Checks provided login credentials and stores them in the $HOME/.flow/credentials.json file",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.CredentialsProvider.Username() == "" || client.CredentialsProvider.Password() == "" {
			return fmt.Errorf("please provide a username and password")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, err := client.Authentication.Login(context.Background(), client.CredentialsProvider.Username(), client.CredentialsProvider.Password())
		if err != nil {
			return err
		}

		return authConfig.SafeWriteConfig()
	},
}
