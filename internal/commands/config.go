package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/flowswiss/goclient"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/term"

	"github.com/flowswiss/cli/v2/pkg/console"
)

const (
	FlagEndpoint = "endpoint"
	FlagToken    = "token"
	FlagDump     = "dump"
	FlagDryRun   = "dry-run"
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

var Config config

type config struct {
	Client   goclient.Client
	Terminal bool
}

func Print(out console.Writer, val interface{}) error {
	format := viper.GetString(FlagFormat)
	if format == FormatJSON {
		return json.NewEncoder(out).Encode(val)
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

	table.Format(out, separator, pretty)

	Stderr.Printf("Found a total of %d items\n", len(table.Rows))
	return nil
}

func PrintStdout(val interface{}) error {
	return Print(Stdout, val)
}

func load() {
	cfg, err := loadConfig()
	if err != nil {
		Stderr.Errorf("%v\n", err)
		os.Exit(1)
	}

	Config = cfg
}

func loadConfig() (config, error) {
	if err := initViper(); err != nil {
		return config{}, err
	}

	endpoint := viper.GetString(FlagEndpoint)
	token := viper.GetString(FlagToken)

	if len(token) == 0 {
		return config{}, fmt.Errorf("missing authentication token")
	}

	opts := []goclient.Option{
		goclient.WithBase(endpoint),
		goclient.WithToken(token),
		goclient.WithUserAgent(fmt.Sprintf("%s-cli/%s", Name, Version)),
	}

	if viper.GetBool(FlagDump) {
		opts = append(opts, goclient.WithHTTPClientOption(func(client *http.Client) {
			client.Transport = dumpRequestTransport{
				delegate: client.Transport,
			}
		}))
	}

	if viper.GetBool(FlagDryRun) {
		opts = append(opts, goclient.WithHTTPClientOption(func(client *http.Client) {
			client.Transport = dryRunTransport{
				delegate: client.Transport,
			}
		}))
	}

	return config{
		Client:   goclient.NewClient(opts...),
		Terminal: term.IsTerminal(int(os.Stdin.Fd())),
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
	baseFlagSet.Bool(FlagDump, false, "dump all requests and responses to stderr")
	baseFlagSet.Bool(FlagDryRun, false, "dry run mode, print requests to stdout instead of sending them to the server")
	baseFlagSet.StringP(FlagFormat, "o", "table", fmt.Sprintf("output format to use. allowed values: %s, %s or %s", FormatTable, FormatCSV, FormatJSON))

	_ = baseFlagSet.MarkHidden(FlagToken)

	Root.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.flow/config.json")
	Root.PersistentFlags().AddFlagSet(baseFlagSet)

	cobra.OnInitialize(load)
}
