package compute

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
)

func AddCommands(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "mac-bare-metal",
		Short: "Manage your mac bare metal devices and networking",
	}

	cmd.AddCommand(
		NetworkCommand(),
		RouterCommand(),
		ElasticIPCommand(),
		SecurityGroupCommand(),
		DeviceCommand(),
	)

	parent.AddCommand(cmd)
}

func init() {
	AddCommands(&commands.Root)
}
