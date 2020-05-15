package commands

import (
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/flowswiss/cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	stdout = output.NewConsoleOutput(os.Stdout)
	stderr = output.NewConsoleOutput(os.Stderr)

	client *flow.Client

	root = &cobra.Command{
		Use:           "flow",
		Short:         "Command line interface for the Flow Platform",
		Version:       "1.0.0",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			base, err := url.Parse(config.Endpoint)
			if err != nil {
				return err
			}

			return initClient(base)
		},
	}
)

func Do() {
	handleError(root.Execute())
}

func handleError(err error) {
	if err != nil {
		stderr.Errorf("%v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	root.AddCommand(authCommand)
	root.AddCommand(computeCommand)
	root.AddCommand(productsCommand)
}

func initClient(base *url.URL) error {
	client = flow.NewClient(base)
	client.CredentialsProvider = &CommandLineCredentialsProvider{}
	client.TokenStorage = &flow.MemoryTokenStorage{}

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

	organizationFilter := viper.GetString("organization")
	if organizationFilter != "" {
		organization, err := findOrganization(organizationFilter)
		if err != nil {
			return err
		}

		client.SelectedOrganization = organization.Id
	}

	return nil
}
