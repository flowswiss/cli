package cli

import "flag"

func initComputeServerCreate(parent *Command) *Command {
	testFlags := flag.NewFlagSet("create", flag.ExitOnError)
	test := testFlags.Int("test", 0, "hello there")

	return &Command{
		Name:   "create",
		Parent: parent,
		Flags:  testFlags,
		Handler: func() error {
			stdout.Printf("create server %d\n", *test)
			return nil
		},
	}
}

func initComputeServer(parent *Command) *Command {
	server := &Command{
		Name:        "server",
		Parent:      parent,
		SubCommands: []*Command{},
	}

	server.SubCommands = append(server.SubCommands, initComputeServerCreate(server))
	return server
}

func initCompute(parent *Command) *Command {
	compute := &Command{
		Name:        "compute",
		Parent:      parent,
		SubCommands: []*Command{},
	}

	compute.SubCommands = append(compute.SubCommands, initComputeServer(compute))
	return compute
}
