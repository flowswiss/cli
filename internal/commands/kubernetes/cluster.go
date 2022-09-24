package compute

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/api/kubernetes"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster",
		Aliases: []string{"clusters"},
		Short:   "Manage your kubernetes cluster",
	}

	commands.Add(cmd,
		&clusterListCommand{},
		&clusterCreateCommand{},
		&clusterUpdateCommand{},
		&clusterDeleteCommand{},
		&clusterUpgradeCommand{},
	)

	cmd.AddCommand(
		ClusterActionCommand(),
		LoadBalancerCommand(),
		NodeCommand(),
		VolumeCommand(),
	)

	return cmd
}

type clusterListCommand struct {
	filter string
}

func (c *clusterListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := kubernetes.NewClusterService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return err
	}

	if len(c.filter) != 0 {
		items = filter.Find(items, c.filter)
	}

	return commands.PrintStdout(items)
}

func (c *clusterListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *clusterListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List all cluster",
		Long:              "Prints a table of all kubernetes clusters belonging to the current organization.",
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().StringVar(&c.filter, "filter", "", "custom term to filter the results")

	return cmd
}

type clusterCreateCommand struct {
	name             string
	location         string
	network          string
	workerProduct    string
	workerCount      int
	attachExternalIP bool
}

func (c *clusterCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Config.Client, c.location)
	if err != nil {
		return err
	}

	products, err := common.ProductsByType(cmd.Context(), commands.Config.Client, common.ProductTypeKubernetesNode)
	if err != nil {
		return fmt.Errorf("fetch products: %w", err)
	}

	workerProduct, err := filter.FindOne(products, c.workerProduct)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}

	networkID := 0
	if c.network != "" {
		networks, err := compute.NewNetworkService(commands.Config.Client).List(cmd.Context())
		if err != nil {
			return fmt.Errorf("fetch networks: %w", err)
		}

		network, err := filter.FindOne(networks, c.network)
		if err != nil {
			return fmt.Errorf("find network: %w", err)
		}

		if network.Location.ID != location.ID {
			return fmt.Errorf("network %s is not available in location %s", network.Name, location.Name)
		}

		networkID = network.ID
	}

	data := kubernetes.ClusterCreate{
		Name:       c.name,
		LocationID: location.ID,
		NetworkID:  networkID,
		Worker: kubernetes.ClusterWorkerCreate{
			ProductID: workerProduct.ID,
			Count:     c.workerCount,
		},
		AttachExternalIP: c.attachExternalIP,
	}

	service := kubernetes.NewClusterService(commands.Config.Client)

	ordering, err := service.Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create cluster: %w", err)
	}

	order, err := commands.WaitForOrder(cmd.Context(), "Creating cluster", ordering)
	if err != nil {
		return fmt.Errorf("wait for order: %w", err)
	}

	cluster, err := service.Get(cmd.Context(), order.Product.ID)
	if err != nil {
		return fmt.Errorf("fetch cluster: %w", err)
	}

	return commands.PrintStdout(cluster)
}

func (c *clusterCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *clusterCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new cluster",
		Long:  "Creates a new kubernetes cluster.",
		Example: commands.FormatExamples(fmt.Sprintf(`
			# Create a new cluster
      %[1]s kubernetes cluster create --name my-cluster --location ALP1 --worker-count 3 --worker-product k1.1x2
		`, commands.Name)),
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().StringVarP(&c.name, "name", "n", "", "name of the cluster (required)")
	cmd.Flags().StringVarP(&c.location, "location", "l", "", "location of the cluster (required)")
	cmd.Flags().StringVar(&c.network, "network", "", "network in which the cluster will be created")
	cmd.Flags().StringVar(&c.workerProduct, "worker-product", "", "product for the worker nodes (required)")
	cmd.Flags().IntVar(&c.workerCount, "worker-count", 3, "number of worker nodes")
	cmd.Flags().BoolVar(&c.attachExternalIP, "attach-external-ip", true, "whether to attach an elastic ip to the cluster")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("location")
	_ = cmd.MarkFlagRequired("worker-product")

	return cmd
}

type clusterUpdateCommand struct {
	name string
}

func (c *clusterUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	data := kubernetes.ClusterUpdate{
		Name: c.name,
	}

	cluster, err = kubernetes.NewClusterService(commands.Config.Client).Update(cmd.Context(), cluster.ID, data)
	if err != nil {
		return fmt.Errorf("update cluster: %w", err)
	}

	return commands.PrintStdout(cluster)
}

func (c *clusterUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeCluster(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *clusterUpdateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update CLUSTER",
		Short:             "Update cluster",
		Long:              "Updates a kubernetes cluster.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().StringVar(&c.name, "name", "", "new name of the cluster")

	return cmd
}

type clusterDeleteCommand struct {
	force bool
}

func (c *clusterDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !c.force && !commands.ConfirmDeletion("kubernetes cluster", cluster) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = kubernetes.NewClusterService(commands.Config.Client).Delete(cmd.Context(), cluster.ID)
	if err != nil {
		return fmt.Errorf("delete cluster: %w", err)
	}

	return nil
}

func (c *clusterDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeCluster(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *clusterDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete CLUSTER",
		Short:             "Delete cluster",
		Long:              "Deletes a kubernetes cluster.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().BoolVar(&c.force, "force", false, "forces deletion of the cluster without asking for confirmation")

	return cmd
}

type clusterUpgradeCommand struct {
	workerProduct string
	workerCount   int
}

func (c *clusterUpgradeCommand) Run(cmd *cobra.Command, args []string) error {
	cluster, err := findCluster(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	products, err := common.ProductsByType(cmd.Context(), commands.Config.Client, common.ProductTypeKubernetesNode)
	if err != nil {
		return fmt.Errorf("fetch products: %w", err)
	}

	workerProduct, err := filter.FindOne(products, c.workerProduct)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}

	data := kubernetes.ClusterUpdateFlavor{
		Worker: kubernetes.ClusterWorkerUpdate{
			ProductID: workerProduct.ID,
			Count:     c.workerCount,
		},
	}

	cluster, err = kubernetes.NewClusterService(commands.Config.Client).UpdateFlavor(cmd.Context(), cluster.ID, data)
	if err != nil {
		return fmt.Errorf("upgrade cluster: %w", err)
	}

	return commands.PrintStdout(cluster)
}

func (c *clusterUpgradeCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeCluster(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (c *clusterUpgradeCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "upgrade CLUSTER",
		Short:             "Upgrade cluster",
		Long:              "Upgrades a kubernetes cluster.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: c.CompleteArg,
		RunE:              c.Run,
	}

	cmd.Flags().StringVar(&c.workerProduct, "worker-product", "", "product for the worker nodes (required)")
	cmd.Flags().IntVar(&c.workerCount, "worker-count", 0, "number of worker nodes (required)")

	_ = cmd.MarkFlagRequired("worker-product")
	_ = cmd.MarkFlagRequired("worker-count")

	return cmd
}

func completeCluster(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	clusters, err := kubernetes.NewClusterService(commands.Config.Client).List(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(clusters, term)

	names := make([]string, len(filtered))
	for i, cluster := range filtered {
		names[i] = cluster.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findCluster(ctx context.Context, term string) (kubernetes.Cluster, error) {
	clusters, err := kubernetes.NewClusterService(commands.Config.Client).List(ctx)
	if err != nil {
		return kubernetes.Cluster{}, fmt.Errorf("fetch clusters: %w", err)
	}

	cluster, err := filter.FindOne(clusters, term)
	if err != nil {
		return kubernetes.Cluster{}, fmt.Errorf("find cluster: %w", err)
	}

	return cluster, nil
}
