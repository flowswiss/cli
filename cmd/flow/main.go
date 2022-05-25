package main

import (
	"github.com/flowswiss/cli/internal/commands"

	_ "github.com/flowswiss/cli/internal/commands/common"
	_ "github.com/flowswiss/cli/internal/commands/macbaremetal"
)

func main() {
	commands.Run()
}
