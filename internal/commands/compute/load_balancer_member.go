package compute

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func LoadBalancerMemberCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "member",
		Aliases: []string{"members"},
		Short:   "Manage load balancer members",
	}

	commands.Add(cmd,
		&loadBalancerMemberListCommand{},
		&loadBalancerMemberCreateCommand{},
		&loadBalancerMemberDeleteCommand{},
	)

	return cmd
}

type loadBalancerMemberListCommand struct {
	filter string
}

func (r *loadBalancerMemberListCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	pool, err := findLoadBalancerPool(cmd.Context(), loadBalancer.ID, args[1])
	if err != nil {
		return err
	}

	service := compute.NewLoadBalancerMemberService(commands.Config.Client, loadBalancer.ID, pool.ID)

	items, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch loadBalancerMembers: %w", err)
	}

	if len(r.filter) != 0 {
		items = filter.Find(items, r.filter)
	}

	return commands.PrintStdout(items)
}

func (r *loadBalancerMemberListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list LOAD-BALANCER POOL",
		Aliases: []string{"show", "ls", "get"},
		Short:   "List load balancer members",
		Long:    "Lists all load balancer member of the selected load balancer pool.",
		Args:    cobra.ExactArgs(2),
		RunE:    r.Run,
	}

	cmd.Flags().StringVar(&r.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type loadBalancerMemberCreateCommand struct {
	name    string
	address net.IP
	port    int
}

func (r *loadBalancerMemberCreateCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	pool, err := findLoadBalancerPool(cmd.Context(), loadBalancer.ID, args[1])
	if err != nil {
		return err
	}

	service := compute.NewLoadBalancerMemberService(commands.Config.Client, loadBalancer.ID, pool.ID)

	data := compute.LoadBalancerMemberCreate{
		Name:    r.name,
		Address: r.address.String(),
		Port:    r.port,
	}

	item, err := service.Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create load balancer member: %w", err)
	}

	return commands.PrintStdout(item)
}

func (r *loadBalancerMemberCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create LOAD-BALANCER POOL",
		Short: "Create a load balancer member",
		Long:  "Creates a new load balancer member.",
		Args:  cobra.ExactArgs(2),
		RunE:  r.Run,
	}

	cmd.Flags().StringVar(&r.name, "name", "", "name of the load balancer member")
	cmd.Flags().IPVar(&r.address, "address", net.IP{}, "ip address of the load balancer member")
	cmd.Flags().IntVar(&r.port, "port", 0, "port of the load balancer member")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("address")
	_ = cmd.MarkFlagRequired("port")

	return cmd
}

type loadBalancerMemberDeleteCommand struct {
	force bool
}

func (r *loadBalancerMemberDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	loadBalancer, err := findLoadBalancer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	pool, err := findLoadBalancerPool(cmd.Context(), loadBalancer.ID, args[1])
	if err != nil {
		return err
	}

	service := compute.NewLoadBalancerMemberService(commands.Config.Client, loadBalancer.ID, pool.ID)

	members, err := service.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch load balancer members: %w", err)
	}

	member, err := filter.FindOne(members, args[2])
	if err != nil {
		return fmt.Errorf("find load balancer member: %w", err)
	}

	if !r.force && !commands.ConfirmDeletion("load balancer member", member) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = service.Delete(cmd.Context(), member.ID)
	if err != nil {
		return fmt.Errorf("delete load balancer member: %w", err)
	}

	return nil
}

func (r *loadBalancerMemberDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete LOAD-BALANCER POOL MEMBER",
		Short: "Delete load balancer member",
		Long:  "Deletes a load balancer member.",
		Args:  cobra.ExactArgs(3),
		RunE:  r.Run,
	}

	cmd.Flags().BoolVar(&r.force, "force", false, "force the deletion of the loadBalancerMember without asking for confirmation")

	return cmd
}
