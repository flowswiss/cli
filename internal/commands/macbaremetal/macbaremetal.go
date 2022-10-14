package macbaremetal

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
)

func Module(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mac-bare-metal",
		Aliases: []string{"macbaremetal"},
		Short:   "Manage your mac bare metal devices and networking",
	}

	cmd.AddCommand(
		NetworkCommand(app),
		RouterCommand(app),
		ElasticIPCommand(app),
		SecurityGroupCommand(app),
		DeviceCommand(app),
	)

	return cmd
}
