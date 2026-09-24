package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/flowswiss/cli/v2/pkg/console"
	"github.com/flowswiss/goclient/v2"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

var Client *goclient.Client

func Print(out console.Writer, val any) error {
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

func PrintStdout(val any) error {
	return Print(Stdout, val)
}

func loadConfig(app Application) {
	if err := initViper(app); err != nil {
		Stderr.Errorf("%v\n", err)
		os.Exit(1)
	}

	client, err := buildClient(app)
	if err != nil {
		Stderr.Errorf("%v\n", err)
		os.Exit(1)
	}

	Client = client
}

func buildClient(app Application) (*goclient.Client, error) {
	endpoint := viper.GetString(FlagEndpoint)
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	token := viper.GetString(FlagToken)
	if len(token) == 0 {
		return nil, fmt.Errorf("missing authentication token")
	}

	httpClient := http.DefaultClient

	if viper.GetBool(FlagDump) {
		httpClient.Transport = dumpRequestTransport{
			delegate: httpClient.Transport,
		}
	}

	if viper.GetBool(FlagDryRun) {
		httpClient.Transport = dryRunTransport{
			delegate: httpClient.Transport,
		}
	}

	coreClient := core.NewClient(core.ClientOpts{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
		UserAgent:  fmt.Sprintf("%s-cli/%s", app.Name, app.Version),
		Token:      token,
	})

	return goclient.WithClient(coreClient), nil
}

func setupFlags(app Application, root *cobra.Command) {
	baseFlagSet = pflag.NewFlagSet("base", pflag.ContinueOnError)
	baseFlagSet.String(FlagEndpoint, app.Endpoint, "base endpoint to use for all api requests")
	baseFlagSet.String(FlagToken, "", "authentication token to use for all api requests")
	baseFlagSet.Bool(FlagDump, false, "dump all requests and responses to stderr")
	baseFlagSet.Bool(FlagDryRun, false, "dry run mode, print requests to stdout instead of sending them to the server")
	baseFlagSet.StringP(FlagFormat, "o", "table", fmt.Sprintf("output format to use. allowed values: %s, %s or %s", FormatTable, FormatCSV, FormatJSON))

	_ = baseFlagSet.MarkHidden(FlagToken)

	root.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is $HOME/.%s/config.json", app.Name))
	root.PersistentFlags().AddFlagSet(baseFlagSet)
}

func initViper(app Application) error {
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		configDir = filepath.Join(home, "."+app.Name)
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

	viper.SetEnvPrefix(app.Name)
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
