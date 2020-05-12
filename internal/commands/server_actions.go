package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/flowswiss/cli/pkg/output"
	"github.com/spf13/cobra"
	"time"
)

const (
	serverStartAction  = "start"
	serverStopAction   = "stop"
	serverRebootAction = "reboot"
)

var (
	serverActionStart = &cobra.Command{
		Use:   "start",
		Short: "Start the selected server",
		RunE:  startServer,
	}

	serverActionStop = &cobra.Command{
		Use:   "stop",
		Short: "Stop the selected server",
		RunE:  stopServer,
	}

	serverActionReboot = &cobra.Command{
		Use:   "reboot",
		Short: "Reboot the selected server",
		RunE:  rebootServer,
	}
)

func init() {
	serverCommand.AddCommand(serverActionStart)
	serverCommand.AddCommand(serverActionStop)
	serverCommand.AddCommand(serverActionReboot)

	serverActionStart.Flags().String("server", "", "identification for the server to execute this action on")
	serverActionStop.Flags().String("server", "", "identification for the server to execute this action on")
	serverActionReboot.Flags().String("server", "", "identification for the server to execute this action on")
}

func findServer(cmd *cobra.Command) (*flow.Server, error) {
	serverFilter, err := cmd.Flags().GetString("server")
	if err != nil {
		return nil, err
	}

	if serverFilter == "" {
		return nil, errRequiredFlag("server")
	}

	servers, _, err := client.Server.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	srv, err := findOne(servers, serverFilter, 2)
	if err != nil {
		return nil, fmt.Errorf("server: %v", err)
	}

	return srv.(*flow.Server), nil
}

func isActionAllowed(command string, server *flow.Server) bool {
	for _, action := range server.Status.Actions {
		if action.Command == command {
			return true
		}
	}
	return false
}

func waitForStatus(id flow.Id, destination flow.Id, allowedStates []flow.Id) error {
	for {
		server, _, err := client.Server.Get(context.Background(), id)
		if err != nil {
			return err
		}

		if server.Status.Id == destination {
			return nil
		}

		found := false
		for _, allowed := range allowedStates {
			if server.Status.Id == allowed {
				found = true
				continue
			}
		}

		if !found {
			return fmt.Errorf("status of the server does not match expectation: %v", server.Status)
		}

		time.Sleep(time.Second)
	}
}

func startServer(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd)
	if err != nil {
		return err
	}

	if !isActionAllowed(serverStartAction, server) {
		return fmt.Errorf("action is not allowed in %s state", server.Status.Key)
	}

	server, _, err = client.Server.RunAction(context.Background(), server.Id, serverStartAction)
	if err != nil {
		return err
	}

	progress := output.NewProgress("starting server")
	go progress.Display(stderr)

	err = waitForStatus(server.Id, 1, []flow.Id{4})
	if err != nil {
		progress.Complete("server failed to start")
		return err
	}

	progress.Complete("server started successfully")
	return nil
}

func stopServer(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd)
	if err != nil {
		return err
	}

	if !isActionAllowed(serverStopAction, server) {
		return fmt.Errorf("action is not allowed in %s state", server.Status.Key)
	}

	server, _, err = client.Server.RunAction(context.Background(), server.Id, serverStopAction)
	if err != nil {
		return err
	}

	progress := output.NewProgress("stopping server")
	go progress.Display(stderr)

	err = waitForStatus(server.Id, 2, []flow.Id{5})
	if err != nil {
		progress.Complete("server failed to stop")
		return err
	}

	progress.Complete("server stopped successfully")
	return nil
}

func rebootServer(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd)
	if err != nil {
		return err
	}

	if !isActionAllowed(serverRebootAction, server) {
		return fmt.Errorf("action is not allowed in %s state", server.Status.Key)
	}

	server, _, err = client.Server.RunAction(context.Background(), server.Id, serverRebootAction)
	if err != nil {
		return err
	}

	progress := output.NewProgress("rebooting server")
	go progress.Display(stderr)

	err = waitForStatus(server.Id, 1, []flow.Id{4, 5})
	if err != nil {
		progress.Complete("server failed to reboot")
		return err
	}

	progress.Complete("server rebooted successfully")
	return nil
}
