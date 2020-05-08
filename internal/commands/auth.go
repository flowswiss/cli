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

func readAuthConfig(dir string) {
	authConfig.AddConfigPath(dir)
	authConfig.SetConfigName("credentials")
	authConfig.SetConfigType(configType)
	authConfig.SetConfigPermissions(0600)

	_ = authConfig.ReadInConfig()
}
