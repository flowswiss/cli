package common

import (
	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/internal/commands"
)

func AddCommands(parent *cobra.Command) {
	parent.AddCommand(
		ModuleCommand(),
		LocationCommand(),
		ProductCommand(),
	)
}

func init() {
	AddCommands(&commands.Root)
}
