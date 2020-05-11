package commands

import (
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/flowswiss/cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
)

type CliConfig struct {
	Verbosity int    `mapstructure:"verbosity"`
	Endpoint  string `mapstructure:"endpoint_url"`

	TwoFactorCode string
}

const configType = "json"

var (
	stdout = output.NewConsoleOutput(os.Stdout)
	stderr = output.NewConsoleOutput(os.Stderr)

	client *flow.Client
	config = &CliConfig{}

	configFile string
	configDir  string

	root = &cobra.Command{
		Use:     "flow",
		Short:   "Command line interface for the Flow Platform",
		Version: "1.0.0",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			base, err := url.Parse(config.Endpoint)
			if err != nil {
				return err
			}

			initClient(base)
			return nil
		},
	}
)

func Do() {
	handleError(root.Execute())
}

func handleError(err error) {
	if err != nil {
		stderr.Errorf("%s\n", err.Error())
		os.Exit(1)
	}
}

func errRequiredFlag(flag string) error {
	return fmt.Errorf("%s is required", flag)
}

func init() {
	cobra.OnInitialize(initConfig)

	root.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.flow.json")

	root.PersistentFlags().CountP("verbosity", "v", "enable a more verbose output (repeat up to 3 times to see entire output)")
	root.PersistentFlags().String("endpoint-url", "https://api.flow.swiss/", "base endpoint to use for all api requests")
	root.PersistentFlags().String("username", "", "name of the user to authenticate with")
	root.PersistentFlags().String("password", "", "password of the user to authenticate with")
	root.PersistentFlags().StringVar(&config.TwoFactorCode, "two-factor-code", "", "two factor code")

	handleError(viper.BindPFlag("verbosity", root.PersistentFlags().Lookup("verbosity")))
	handleError(viper.BindPFlag("endpoint_url", root.PersistentFlags().Lookup("endpoint-url")))
	handleError(authConfig.BindPFlag("username", root.PersistentFlags().Lookup("username")))
	handleError(authConfig.BindPFlag("password", root.PersistentFlags().Lookup("password")))

	root.AddCommand(authCommand)
	root.AddCommand(computeCommand)
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

func initClient(base *url.URL) {
	client = flow.NewClient(base)
	client.CredentialsProvider = &CommandLineCredentialsProvider{}
	client.TokenStorage = &flow.MemoryTokenStorage{}
	client.Organization = 2

	client.OnRequest = func(req *http.Request) {
		if config.Verbosity >= 1 {
			stderr.Printf("Requesting %s %s\n", req.Method, req.URL)
		}

		if config.Verbosity >= 3 {
			dump, err := httputil.DumpRequestOut(req, true)
			if err == nil {
				stderr.Color(output.AnsiBright+output.AnsiBlack).Printf("%s\n", dump).Reset()
			}
		}
	}

	client.OnResponse = func(res *http.Response) {
		if config.Verbosity >= 2 {
			dump, err := httputil.DumpResponse(res, true)
			if err == nil {
				stderr.Color(output.AnsiBright+output.AnsiBlack).Printf("%s\n\n", dump).Reset()
			}
		}
	}
}
