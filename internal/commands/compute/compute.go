package compute

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
)

func AddCommands(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Manage your compute server and networking",
	}

	cmd.AddCommand(
		ElasticIPCommand(),
		ImageCommand(),
		KeyPairCommand(),
		NetworkCommand(),
		RouterCommand(),
		SecurityGroupCommand(),
		ServerCommand(),
		SnapshotCommand(),
		VolumeCommand(),
	)

	parent.AddCommand(cmd)
}

func init() {
	AddCommands(&commands.Root)
}
