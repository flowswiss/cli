package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authConfig = viper.New()

var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication settings",
}

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
