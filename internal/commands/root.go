package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/pkg/console"
)

var (
	Stdout = console.NewConsoleOutput(os.Stdout)
	Stderr = console.NewConsoleOutput(os.Stderr)
)

var Root = cobra.Command{}

type ModuleFactory func(app Application) *cobra.Command

type Application struct {
	Name        string
	Description string
	Version     string
	Endpoint    string

	Modules []ModuleFactory
}

func Run(app Application) {
	root := cobra.Command{
		Use:           app.Name,
		Short:         app.Description,
		Version:       app.Version,
		SilenceErrors: true,
	}

	for _, module := range app.Modules {
		cmd := module(app)
		root.AddCommand(cmd)
	}

	setupFlags(app, &root)
	cobra.OnInitialize(func() {
		loadConfig(app)
	})

	err := root.ExecuteContext(context.Background())
	if err != nil {
		Stderr.Errorf("%v\n", err)
		os.Exit(1)
	}
}
