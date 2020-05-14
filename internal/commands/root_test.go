package commands

import (
	"os"
)

func setupTests() {
	configDir = os.TempDir()
	_ = initClient(nil)
}
