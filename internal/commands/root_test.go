package commands

import (
	"os"
)

func setupTests() {
	configDir = os.TempDir()
	initClient(nil)
}
