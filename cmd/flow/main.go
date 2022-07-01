package main

import (
	"github.com/flowswiss/cli/v2/internal/commands"

	_ "github.com/flowswiss/cli/v2/internal/commands/common"
	_ "github.com/flowswiss/cli/v2/internal/commands/compute"
	_ "github.com/flowswiss/cli/v2/internal/commands/macbaremetal"
)

func main() {
	commands.Run()
}
