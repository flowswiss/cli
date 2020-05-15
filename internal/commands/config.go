package commands

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	configType = "json"

	flagConfig        = "config"
	flagVerbosity     = "verbosity"
	flagEndpointUrl   = "endpoint-url"
	flagUsername      = "username"
	flagPassword      = "password"
	flagOrganization  = "organization"
	flagTwoFactorCode = "two-factor-code"
)

var (
	config = &CliConfig{}

	configFile string
	configDir  string
)

type CliConfig struct {
	Endpoint string `mapstructure:"endpoint_url"`
	Format   string `mapstructure:"format"`

	TwoFactorCode string
	Verbosity     int
}

func init() {
	root.PersistentFlags().StringVar(&configFile, flagConfig, "", "config file (default is $HOME/.flow/config.json")

	root.PersistentFlags().CountVarP(&config.Verbosity, flagVerbosity, "v", "enable a more verbose output (repeat up to 3 times to see entire output)")
	root.PersistentFlags().String(flagEndpointUrl, "https://api.flow.swiss/", "base endpoint to use for all api requests")
	root.PersistentFlags().String(flagOrganization, "", "the organization context to use for every request")
	root.PersistentFlags().String(flagUsername, "", "name of the user to authenticate with")
	root.PersistentFlags().String(flagPassword, "", "password of the user to authenticate with")
	root.PersistentFlags().StringVar(&config.TwoFactorCode, flagTwoFactorCode, "", "two factor code")

	handleError(viper.BindPFlag("endpoint_url", root.PersistentFlags().Lookup(flagEndpointUrl)))
	handleError(viper.BindPFlag("organization", root.PersistentFlags().Lookup(flagOrganization)))
	handleError(authConfig.BindPFlag("username", root.PersistentFlags().Lookup(flagUsername)))
	handleError(authConfig.BindPFlag("password", root.PersistentFlags().Lookup(flagPassword)))
}

func configureConfig(name string, conf *viper.Viper) {
	conf.AddConfigPath(configDir)
	conf.SetConfigName(name)
	conf.SetConfigType(configType)
	conf.SetEnvPrefix("flow")
	conf.AutomaticEnv()
}

func initConfig() {
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			handleError(err)
		}

		configDir = filepath.Join(home, ".flow")
	}

	_ = os.Mkdir(configDir, 0755)

	configureConfig("config", viper.GetViper())

	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	_ = viper.ReadInConfig()
	handleError(viper.Unmarshal(config))

	readAuthConfig()
}
