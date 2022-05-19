package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flowswiss/goclient"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/flowswiss/cli/pkg/console"
)

const (
	FlagEndpoint = "endpoint"
	FlagToken    = "token"
	FlagFormat   = "format"
)

const (
	FormatJSON  = "json"
	FormatTable = "table"
	FormatCSV   = "csv"
)

var (
	configFile string
	configDir  string

	baseFlagSet *pflag.FlagSet
)

type Config struct {
	Client goclient.Client
}

func Format(val interface{}) error {
	format := viper.GetString(FlagFormat)
	if format == FormatJSON {
		return json.NewEncoder(Stdout.Writer).Encode(val)
	}

	separator := "   "
	pretty := true

	if format == FormatCSV {
		separator = ","
		pretty = false
	}

	table := console.Table{}

	err := table.Insert(val)
	if err != nil {
		return err
	}

	table.Format(Stdout, separator, pretty)

	Stderr.Printf("Found a total of %d items\n", len(table.Rows))
	return nil
}

func loadConfig() (Config, error) {
	if err := initViper(); err != nil {
		return Config{}, err
	}

	endpoint := viper.GetString(FlagEndpoint)
	token := viper.GetString(FlagToken)

	if len(token) == 0 {
		return Config{}, fmt.Errorf("missing authentication token")
	}

	return Config{
		Client: goclient.NewClient(
			goclient.WithBase(endpoint),
			goclient.WithToken(token),
			goclient.WithUserAgent(fmt.Sprintf("%s/%s", Name, Version)),
		),
	}, nil
}

func initViper() error {
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		configDir = filepath.Join(home, ".flow")
	}

	if err := os.Mkdir(configDir, 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	viper.SetConfigPermissions(0600)

	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	if len(configFile) != 0 {
		viper.SetConfigFile(configFile)
	}

	viper.SetEnvPrefix("flow")
	viper.AutomaticEnv()

	if err := viper.BindPFlags(baseFlagSet); err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		// ignore config not found error if not manually specified
		if len(configFile) != 0 || !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return err
		}
	}

	return nil
}

func init() {
	baseFlagSet = pflag.NewFlagSet("base", pflag.ContinueOnError)
	baseFlagSet.String(FlagEndpoint, DefaultEndpoint, "base endpoint to use for all api requests")
	baseFlagSet.String(FlagToken, "", "authentication token to use for all api requests")
	baseFlagSet.StringP(FlagFormat, "o", "table", fmt.Sprintf("output format to use. allowed values: %s, %s or %s", FormatTable, FormatCSV, FormatJSON))

	_ = baseFlagSet.MarkHidden(FlagToken)

	Root.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.flow/config.json")
	Root.PersistentFlags().AddFlagSet(baseFlagSet)
}
