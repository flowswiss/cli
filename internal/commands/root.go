package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/pkg/console"
)

var (
	Name        = "flow"
	Description = "flow is a command-line interface for managing the Flow Swiss cloud platform."
	Version     = "v2.0.0-beta.1"

	DefaultEndpoint = "https://api.flow.swiss/"
)

var (
	Stdout = console.NewConsoleOutput(os.Stdout)
	Stderr = console.NewConsoleOutput(os.Stderr)
)

var Root = cobra.Command{
	Use:           Name,
	Short:         Description,
	Version:       Version,
	SilenceErrors: true,
}

func Run() {
	err := Root.ExecuteContext(context.Background())
	if err != nil {
		Stderr.Errorf("%v\n", err)
		os.Exit(1)
	}
}
