package cli

import (
	"github.com/flowswiss/cli/pkg/output"
	"os"
)

const (
	errorTooFewArguments = cliError("too few arguments")
	errorInvalidCommand  = cliError("selected command does not exist")

	errorPrintUsage = cliError("print usage")
)

var (
	stdout *output.Output
	stderr *output.Output

	root *Command
)

type cliError string

func (e cliError) Error() string {
	return string(e)
}

func init() {
	stdout = output.NewConsoleOutput(os.Stdout)
	stderr = output.NewConsoleOutput(os.Stderr)

	root = &Command{
		Name:        os.Args[0],
		SubCommands: []*Command{},
	}

	root.SubCommands = append(root.SubCommands, initCompute(root))
}

func Do() {
	err := root.Do(os.Args)
	if err != nil {
		stderr.Errorf("%s\n", err.Error())
		os.Exit(1)
	}
}
