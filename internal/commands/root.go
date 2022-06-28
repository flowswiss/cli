package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands/build"
	"github.com/flowswiss/cli/v2/pkg/console"
)

var (
	Name            = build.Name
	Version         = build.Version
	DefaultEndpoint = build.DefaultEndpoint
)

var (
	Stdout = console.NewConsoleOutput(os.Stdout)
	Stderr = console.NewConsoleOutput(os.Stderr)
)

var Root = cobra.Command{
	Use:           Name,
	Short:         build.Description,
	Version:       build.Version,
	SilenceErrors: true,
}

func Run() {
	err := Root.ExecuteContext(context.Background())
	if err != nil {
		Stderr.Errorf("%v\n", err)
		os.Exit(1)
	}
}
