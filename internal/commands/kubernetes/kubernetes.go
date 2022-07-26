package compute

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
)

func AddCommands(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "kubernetes",
		Short: "Manage your kubernetes clusters",
	}

	cmd.AddCommand(
		ClusterCommand(),
	)

	parent.AddCommand(cmd)
}

func init() {
	AddCommands(&commands.Root)
}
