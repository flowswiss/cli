//go:build !cloudbit

package build

const (
	Name        = "flow"
	Description = "flow is a command-line interface for managing the Flow Swiss cloud platform."
	Version     = "dev"

	DefaultEndpoint = "https://api.flow.swiss/"
)
