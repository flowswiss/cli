package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ElasticIPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elastic-ip",
		Aliases: []string{"elastic-ips", "elasticip", "elasticips"},
		Short:   "Manage compute elastic ips",
		Example: commands.FormatExamples(fmt.Sprintf(`
  			# List all compute elastic ips
	  		%[1]s compute elastic-ip list	

			# Create a new compute elastic ip
			%[1]s compute elastic-ip create --location=ZRH1

			# Attach a compute elastic ip to a server
			%[1]s compute elastic-ip attach 1.1.1.1 my-server
		`, commands.Name)),
	}

	commands.Add(cmd, &elasticIPListCommand{}, &elasticIPCreateCommand{}, &elasticIPDeleteCommand{}, &elasticIPAttachCommand{}, &elasticIPDetachCommand{})

	return cmd
}

type elasticIPListCommand struct {
	filter string
}

func (e *elasticIPListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.NewElasticIPService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch elastic ips: %w", err)
	}

	if len(e.filter) != 0 {
		items = filter.Find(items, e.filter)
	}

	return commands.PrintStdout(items)
}

func (e *elasticIPListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (e *elasticIPListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List elastic ips",
		Long:              "Lists all compute elastic ips.",
		ValidArgsFunction: e.CompleteArg,
		RunE:              e.Run,
	}

	cmd.Flags().StringVar(&e.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type elasticIPCreateCommand struct {
	location string
}

func (e *elasticIPCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Config.Client, e.location)
	if err != nil {
		return err
	}

	data := compute.ElasticIPCreate{
		LocationID: location.ID,
	}

	item, err := compute.NewElasticIPService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create elastic ip: %w", err)
	}

	return commands.PrintStdout(item)
}

func (e *elasticIPCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (e *elasticIPCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Aliases:           []string{"add", "new"},
		Short:             "Create new elastic ip",
		Long:              "Creates a new compute elastic ip.",
		ValidArgsFunction: e.CompleteArg,
		RunE:              e.Run,
	}

	cmd.Flags().StringVar(&e.location, "location", "", "location where the elastic ip will be created")
	_ = cmd.MarkFlagRequired("location")

	return cmd
}

type elasticIPDeleteCommand struct {
	force bool
}

func (e *elasticIPDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	service := compute.NewElasticIPService(commands.Config.Client)

	elasticIPs, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch elastic ips: %w", err)
	}

	elasticIP, err := filter.FindOne(elasticIPs, args[0])
	if err != nil {
		return fmt.Errorf("find elastic ip: %w", err)
	}

	if elasticIP.Attachment.ID != 0 {
		commands.Stderr.Errorf("WARNING: The elastic ip is still attached to a server. Active connections to the server might get disturbed.\n")
	}

	if !e.force && !commands.ConfirmDeletion("elastic ip", elasticIP) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	if elasticIP.Attachment.ID != 0 {
		err = service.Detach(cmd.Context(), elasticIP.Attachment.ID, elasticIP.ID)
		if err != nil {
			return fmt.Errorf("detach elastic ip: %w", err)
		}
	}

	err = service.Delete(cmd.Context(), elasticIP.ID)
	if err != nil {
		return fmt.Errorf("delete elastic ip: %w", err)
	}

	return nil
}

func (e *elasticIPDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeElasticIP(cmd.Context(), toComplete, nil)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (e *elasticIPDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete ELASTIC-IP",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete elastic ip",
		Long:    "Deletes a compute elastic ip.",
		Example: commands.FormatExamples(fmt.Sprintf(`
	  		# Delete a compute elastic ip
			%[1]s compute elastic-ip delete 1.1.1.1

			# Force the deletion a compute elastic ip without confirmation
			%[1]s compute elastic-ip delete 1.1.1.1 --force
		`, commands.Name)),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: e.CompleteArg,
		RunE:              e.Run,
	}

	cmd.Flags().BoolVar(&e.force, "force", false, "force the deletion of the elastic ip without asking for confirmation")

	return cmd
}

type elasticIPAttachCommand struct {
	networkInterface string
}

func (e *elasticIPAttachCommand) Run(cmd *cobra.Command, args []string) error {
	elasticIP, err := findElasticIP(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if elasticIP.Attachment.ID != 0 {
		return fmt.Errorf("elastic ip is already attached to a server")
	}

	server, err := findServer(cmd.Context(), args[1])
	if err != nil {
		return err
	}

	data := compute.ElasticIPAttach{
		ElasticIPID: elasticIP.ID,
	}

	if e.networkInterface != "" {
		networkInterface, err := findNetworkInterface(cmd.Context(), server.ID, e.networkInterface)
		if err != nil {
			return err
		}

		data.NetworkInterfaceID = networkInterface.ID
	} else {
	searchFreeNetworkInterface:
		for _, network := range server.Networks {
			for _, iface := range network.Interfaces {
				if iface.PublicIP == "" {
					data.NetworkInterfaceID = iface.ID
					break searchFreeNetworkInterface
				}
			}
		}

		if data.NetworkInterfaceID == 0 {
			return fmt.Errorf("server has no free network interface to attach the elastic ip to")
		}
	}

	_, err = compute.NewElasticIPService(commands.Config.Client).Attach(cmd.Context(), server.ID, data)
	if err != nil {
		return fmt.Errorf("attach elastic ip: %w", err)
	}

	return nil
}

func (e *elasticIPAttachCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeElasticIP(cmd.Context(), toComplete, func(item compute.ElasticIP) bool {
			return item.Attachment.ID == 0
		})
	}

	if len(args) == 1 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (e *elasticIPAttachCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "attach ELASTIC-IP SERVER",
		Short:             "Attach elastic ip to server",
		Long:              "Attaches a compute elastic ip to a server.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: e.CompleteArg,
		RunE:              e.Run,
	}

	return cmd
}

type elasticIPDetachCommand struct {
	force bool
}

func (e *elasticIPDetachCommand) Run(cmd *cobra.Command, args []string) error {
	elasticIP, err := findElasticIP(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	server, err := findServer(cmd.Context(), args[1])
	if err != nil {
		return err
	}

	if elasticIP.Attachment.ID != server.ID {
		return fmt.Errorf("elastic ip not attached to the selected server")
	}

	if !e.force && !commands.Confirm(fmt.Sprintf("Are you sure you want to detach the elastic ip %q? Active connection to the server might get disturbed.", elasticIP)) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.NewElasticIPService(commands.Config.Client).Detach(cmd.Context(), server.ID, elasticIP.ID)
	if err != nil {
		return fmt.Errorf("detach elastic ip: %w", err)
	}

	return nil
}

func (e *elasticIPDetachCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeElasticIP(cmd.Context(), toComplete, func(item compute.ElasticIP) bool {
			return item.Attachment.ID != 0
		})
	}

	if len(args) == 1 {
		elasticIP, err := findElasticIP(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return []string{elasticIP.Attachment.Name}, cobra.ShellCompDirectiveNoFileComp
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (e *elasticIPDetachCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "detach ELASTIC-IP SERVER",
		Short:             "Detach elastic ip from server",
		Long:              "Detaches a compute elastic ip from a server.",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: e.CompleteArg,
		RunE:              e.Run,
	}

	cmd.Flags().BoolVar(&e.force, "force", false, "force the detachment of the elastic ip without asking for confirmation")

	return cmd
}

func completeElasticIP(ctx context.Context, term string, itemFilter func(ip compute.ElasticIP) bool) ([]string, cobra.ShellCompDirective) {
	elasticIPs, err := compute.NewElasticIPService(commands.Config.Client).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.FindWithCustomFilter(elasticIPs, term, itemFilter)

	names := make([]string, len(filtered))
	for i, ip := range filtered {
		names[i] = ip.PublicIP
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findElasticIP(ctx context.Context, term string) (compute.ElasticIP, error) {
	elasticIPs, err := compute.NewElasticIPService(commands.Config.Client).List(ctx)
	if err != nil {
		return compute.ElasticIP{}, fmt.Errorf("fetch elastic ips: %w", err)
	}

	elasticIP, err := filter.FindOne(elasticIPs, term)
	if err != nil {
		return compute.ElasticIP{}, fmt.Errorf("find elastic ip: %w", err)
	}

	return elasticIP, nil
}
