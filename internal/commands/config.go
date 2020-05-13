package commands

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const configType = "json"

var (
	config = &CliConfig{}

	configFile string
	configDir  string
)

type CliConfig struct {
	Endpoint string `mapstructure:"endpoint_url"`
	Format string `mapstructure:"format"`

	TwoFactorCode string
	Verbosity     int
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