package optional

import "github.com/spf13/cobra"

type Optional[T any] struct {
	value T

	cmd  *cobra.Command
	name string
}

func (o *Optional[T]) Configure(cmd *cobra.Command, name string, usage string) (*T, string, T, string) {
	o.cmd = cmd
	o.name = name

	var def T
	return &o.value, name, def, usage
}

func (o *Optional[T]) Value() *T {
	if !o.cmd.Flags().Changed(o.name) {
		return nil
	}

	return &o.value
}
