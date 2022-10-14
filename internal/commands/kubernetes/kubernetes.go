package kubernetes

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
)

func Module(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes",
		Short: "Manage your kubernetes clusters",
	}

	cmd.AddCommand(
		ClusterCommand(app),
	)

	return cmd
}
