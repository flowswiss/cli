package compute

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
)

func Module(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Manage your compute server and networking",
	}

	cmd.AddCommand(
		CertificateCommand(app),
		ElasticIPCommand(app),
		ImageCommand(app),
		KeyPairCommand(app),
		LoadBalancerCommand(app),
		NetworkCommand(app),
		RouterCommand(app),
		SecurityGroupCommand(app),
		ServerCommand(app),
		SnapshotCommand(app),
		VolumeCommand(app),
	)

	return cmd
}
