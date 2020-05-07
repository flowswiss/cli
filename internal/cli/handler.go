package cli

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Parent      *Command
	SubCommands []*Command
	Flags       *flag.FlagSet

	Handler func() error
}

func (c *Command) Path() string {
	if c.Parent == nil {
		return c.Name
	}

	return fmt.Sprintf("%s %s", c.Parent.Path(), c.Name)
}

func (c *Command) HasSubCommands() bool {
	return c.SubCommands != nil && len(c.SubCommands) > 0
}

func (c *Command) HasOptions() bool {
	return c.Flags.NFlag() > 0
}

func (c *Command) PrintUsage() {
	stderr.Printf("Usage:\n  %s", c.Path())
	if c.HasSubCommands() {
		fmt.Printf(" <subcommand>")
	}
	if c.HasOptions() {
		fmt.Println(" [options...]")
	}
	stderr.Printf("\n")

	if c.SubCommands != nil && len(c.SubCommands) > 0 {
		stderr.Bold("\nAvailable sub commands:\n")
		for _, cmd := range c.SubCommands {
			stderr.Printf("  - %s\n", cmd.Name)
		}
	}

	stderr.Bold("\nAvailable options:\n")
	c.Flags.SetOutput(stderr.Writer)
	c.Flags.PrintDefaults()
}

func (c *Command) Handle(args []string) error {
	if c.Flags == nil {
		c.Flags = flag.NewFlagSet(c.Name, flag.ExitOnError)
	}

	c.Flags.Usage = c.PrintUsage

	err := c.Flags.Parse(args)
	if err != nil {
		return err
	}

	if c.Handler != nil {
		return c.Handler()
	}

	if c.Flags.NArg() == 0 {
		return errorTooFewArguments
	}

	for _, cmd := range c.SubCommands {
		if cmd.Name == c.Flags.Arg(0) {
			return cmd.Do(c.Flags.Args())
		}
	}

	return errorInvalidCommand
}

func (c *Command) Do(args []string) error {
	err := c.Handle(args[1:])
	if _, ok := err.(cliError); err != nil && ok {
		if err != errorPrintUsage {
			stderr.Printf("%s\n", err.Error())
		}

		c.PrintUsage()
		os.Exit(1)
	}
	return err
}
