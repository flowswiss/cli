package commands

import (
	"context"

	"github.com/spf13/cobra"
)

type Command interface {
	Run(ctx context.Context, cfg Config, args []string) error
	Desc() *cobra.Command
}

type CommandFunc func(ctx context.Context, cfg Config) error

func Build(cmd Command) *cobra.Command {
	base := cmd.Desc()
	base.RunE = wrapCommand(cmd)
	return base
}

func Add(parent *cobra.Command, children ...Command) {
	for _, child := range children {
		parent.AddCommand(Build(child))
	}
}

func wrapCommand(c Command) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		return c.Run(cmd.Context(), cfg, args)
	}
}
