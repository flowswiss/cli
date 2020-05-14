package commands

import (
	"github.com/spf13/cobra"
)

var computeCommand = &cobra.Command{
	Use:   "compute",
	Short: "Manage your compute server and networking",
}

func init() {
	computeCommand.AddCommand(serverCommand)
	computeCommand.AddCommand(keyPairCommand)
}


