package objectstorage

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
)

func Module(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "object-storage",
		Aliases: []string{"objectstorage"},
		Short:   "Manage your object storage",
	}

	cmd.AddCommand(
		InstanceCommand(app),
	)

	return cmd
}
