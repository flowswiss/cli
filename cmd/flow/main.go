package main

import (
	"github.com/flowswiss/cli/internal/commands"

	_ "github.com/flowswiss/cli/internal/commands/common"
)

func main() {
	commands.Run()
}
