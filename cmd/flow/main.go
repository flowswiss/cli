package main

import (
	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/internal/commands/common"
	"github.com/flowswiss/cli/v2/internal/commands/compute"
	"github.com/flowswiss/cli/v2/internal/commands/kubernetes"
	"github.com/flowswiss/cli/v2/internal/commands/macbaremetal"
	"github.com/flowswiss/cli/v2/internal/commands/objectstorage"
)

var Version = "dev"

func main() {
	app := commands.Application{
		Name:        "flow",
		Description: "flow is a command-line interface for managing the Flow Swiss cloud platform.",
		Version:     Version,
		Endpoint:    "https://api.flow.swiss/",

		Modules: []commands.ModuleFactory{
			common.Location,
			common.Module,
			common.Product,

			compute.Module,
			kubernetes.Module,
			macbaremetal.Module,
			objectstorage.Module,
		},
	}

	commands.Run(app)
}
