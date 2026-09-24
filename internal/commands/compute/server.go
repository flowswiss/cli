package compute

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"

	"github.com/flowswiss/cli/v2/pkg/optional"
	"github.com/spf13/cobra"

	"github.com/flowswiss/goclient/v2/core"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/console"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ServerCommand(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"servers"},
		Short:   "Manage your compute server",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # List all servers
      %[1]s compute server list
      
      # Create a new server
      %[1]s compute server create --name my-server --location ALP1 --image linux-ubuntu-20.04-lts --product b1.4x8 --key-pair my-keypair
      
      # Delete a server
      %[1]s compute server delete my-server
		`, app.Name)),
	}

	commands.Add(app, cmd,
		&serverListCommand{},
		&serverCreateCommand{},
		&serverUpdateCommand{},
		&serverUpgradeCommand{},
		&serverDeleteCommand{},
	)

	commands.Add(app, cmd,
		serverActionRunCommandPreset("start"),
		serverActionRunCommandPreset("stop"),
		serverActionRunCommandPreset("restart"),
	)

	cmd.AddCommand(NetworkInterfaceCommand(app), ServerActionCommand(app), ServerVolumeCommand(app))

	return cmd
}

type serverListCommand struct {
	filter string
}

func (s *serverListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.ServerService().List(cmd.Context(), core.CursorAll)
	if err != nil {
		return err
	}

	if len(s.filter) != 0 {
		items = filter.Find(items, s.filter)
	}

	return commands.PrintStdout(items)
}

func (s *serverListCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverListCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"show", "ls", "get"},
		Short:             "List all server",
		Long:              "Prints a table of all compute servers belonging to the current organization.",
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.filter, "filter", "", "") // TODO

	return cmd
}

type serverCreateCommand struct {
	name             string
	location         string
	image            string
	product          string
	network          string
	privateIP        optional.Optional[net.IP]
	keyPair          string
	password         optional.Optional[string]
	cloudInitFile    optional.Optional[string]
	attachExternalIP bool
}

func (s *serverCreateCommand) Run(cmd *cobra.Command, args []string) error {
	location, err := common.FindLocation(cmd.Context(), commands.Client, s.location)
	if err != nil {
		return err
	}

	images, err := compute.Images(cmd.Context())
	if err != nil {
		return fmt.Errorf("fetch images: %w", err)
	}

	image, err := filter.FindOne(images, s.image)
	if err != nil {
		return fmt.Errorf("find image: %w", err)
	}

	if !image.AvailableAt(location) {
		return fmt.Errorf("image %s is not available in location %s", image, location.Name)
	}

	products, err := common.ProductsByType(cmd.Context(), commands.Client, common.ProductTypeComputeServer)
	if err != nil {
		return fmt.Errorf("fetch products: %w", err)
	}

	product, err := filter.FindOne(products, s.product)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}

	networkID := 0
	if s.network != "" {
		networks, err := compute.NetworkService().List(cmd.Context(), core.CursorAll)
		if err != nil {
			return fmt.Errorf("fetch networks: %w", err)
		}

		network, err := filter.FindOne(networks, s.network)
		if err != nil {
			return fmt.Errorf("find network: %w", err)
		}

		if network.Location.ID != location.ID {
			return fmt.Errorf("network %s is not available in location %s", network.Name, location.Name)
		}

		_, cidr, err := net.ParseCIDR(network.CIDR)
		if err != nil {
			return fmt.Errorf("parse network cidr: %w", err)
		}

		if pIP := s.privateIP.Value(); pIP != nil {
			if !cidr.Contains(*pIP) {
				return fmt.Errorf("private ip %s is not in network %s", *pIP, network.CIDR)
			}
		}

		networkID = network.ID
	}

	var privateIP *string
	if pIP := s.privateIP.Value(); pIP != nil {
		privateIP = new(pIP.String())
	}

	var keyPairID *int
	if s.keyPair != "" {
		keyPairs, err := compute.KeyPairService().List(cmd.Context(), core.CursorAll)
		if err != nil {
			return fmt.Errorf("fetch key pairs: %w", err)
		}

		keyPair, err := filter.FindOne(keyPairs, s.keyPair)
		if err != nil {
			return fmt.Errorf("find key pair: %w", err)
		}

		keyPairID = &keyPair.ID
	}

	if !image.IsWindows() && keyPairID == nil {
		return fmt.Errorf("key pair is required for non-windows images")
	}

	var password *string
	if image.IsWindows() {
		pw := s.password.Value()
		if pw == nil || len(*pw) == 0 {
			enteredPW, err := console.Password(commands.Stderr, "Windows User Password", checkWindowsPassword)
			if err != nil {
				return fmt.Errorf("read user password: %w", err)
			}

			pw = &enteredPW
		}

		if err = checkWindowsPassword(*pw); err != nil {
			return fmt.Errorf("check user password: %w", err)
		}
	}

	var cloudInit *string
	if ci := s.cloudInitFile.Value(); ci != nil {
		data, err := os.ReadFile(*ci)
		if err != nil {
			return fmt.Errorf("read cloud init file: %w", err)
		}

		cloudInit = new(base64.StdEncoding.EncodeToString(data))
	}

	data := compute.ServerCreate{
		Name:             s.name,
		LocationID:       location.ID,
		ImageID:          image.ID,
		ProductID:        product.ID,
		AttachExternalIP: s.attachExternalIP,
		NetworkID:        networkID,
		PrivateIP:        privateIP,
		KeyPairID:        keyPairID,
		Password:         password,
		CloudInit:        cloudInit,
	}

	service := compute.ServerService()

	ordering, err := service.Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create server: %w", err)
	}

	order, err := commands.WaitForOrder(cmd.Context(), "Creating server", ordering)
	if err != nil {
		return fmt.Errorf("wait for order: %w", err)
	}

	server, err := service.Get(cmd.Context(), compute.ServerGet{ID: uint(order.Product.ID)})
	if err != nil {
		return fmt.Errorf("fetch server: %w", err)
	}

	return commands.PrintStdout(server)
}

func (s *serverCreateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverCreateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new server",
		Long:  "Creates a new compute server.",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # Create a new ubuntu server
      %[1]s compute server create --name my-server --location ALP1 --image linux-ubuntu-20.04-lts --product b1.4x8 --key-pair my-keypair
      
      # Create a new windows server
      %[1]s compute server create --name my-server --location ALP1 --image microsoft-windows-server-2019 --product b1.2x8
		`, app.Name)), // TODO select correct image names
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVarP(&s.name, "name", "n", "", "name of the new server (required)")
	cmd.Flags().StringVarP(&s.location, "location", "l", "", "location of the server (required)")
	cmd.Flags().StringVarP(&s.image, "image", "i", "", "operating system image to use for the new server (required)")
	cmd.Flags().StringVarP(&s.product, "product", "p", "", "product to use for the new server (required)")
	cmd.Flags().StringVar(&s.network, "network", "", "network in which the first network interface should be created")
	cmd.Flags().IPVar(s.privateIP.Configure(cmd, "private-ip", "ip address of the server in the selected network"))
	cmd.Flags().StringVar(&s.keyPair, "key-pair", "", "ssh key-pair for connecting to the server (required if image is linux)")
	cmd.Flags().StringVar(s.password.Configure(cmd, "windows-password", "password for the windows admin user  (required if image is windows)"))
	cmd.Flags().StringVar(s.cloudInitFile.Configure(cmd, "cloud-init", "cloud init script to customize creation of the server"))
	cmd.Flags().BoolVar(&s.attachExternalIP, "attach-external-ip", true, "whether to attach an elastic ip to the server")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("location")
	_ = cmd.MarkFlagRequired("image")
	_ = cmd.MarkFlagRequired("product")
	_ = cmd.MarkFlagFilename("cloud-init")

	return cmd
}

type serverUpdateCommand struct {
	name string
}

func (s *serverUpdateCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	data := compute.ServerUpdate{ID: uint(server.ID),
		Name: s.name,
	}

	server, err = compute.ServerService().Update(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("update server: %w", err)
	}

	return commands.PrintStdout(server)
}

func (s *serverUpdateCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverUpdateCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update SERVER",
		Short:             "Update server",
		Long:              "Updates a compute server.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.name, "name", "", "new name of the server")

	return cmd
}

type serverUpgradeCommand struct {
	product string
}

func (s *serverUpgradeCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	products, err := common.ProductsByType(cmd.Context(), commands.Client, common.ProductTypeComputeServer)
	if err != nil {
		return fmt.Errorf("fetch products: %w", err)
	}

	product, err := filter.FindOne(products, s.product)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}

	data := compute.ServerUpgrade{ID: uint(server.ID),
		ProductID: product.ID,
	}

	service := compute.ServerService()

	ordering, err := service.Upgrade(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("upgrade server: %w", err)
	}

	order, err := commands.WaitForOrder(cmd.Context(), "Upgrading server", ordering)
	if err != nil {
		return fmt.Errorf("wait for order: %w", err)
	}

	server, err = service.Get(cmd.Context(), compute.ServerGet{ID: uint(order.Product.ID)})
	if err != nil {
		return fmt.Errorf("fetch server: %w", err)
	}

	return commands.PrintStdout(server)
}

func (s *serverUpgradeCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverUpgradeCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "upgrade SERVER",
		Short:             "Upgrade server",
		Long:              "Upgrades a compute server.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().StringVar(&s.product, "product", "", "product to use for the new server")

	_ = cmd.MarkFlagRequired("product")

	return cmd
}

type serverDeleteCommand struct {
	force      bool
	detachOnly bool
}

func (s *serverDeleteCommand) Run(cmd *cobra.Command, args []string) error {
	server, err := findServer(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	if !s.force && !commands.ConfirmDeletion("server", server) {
		commands.Stderr.Println("aborted.")
		return nil
	}

	err = compute.ServerService().Delete(cmd.Context(), compute.ServerDelete{ID: uint(server.ID), DeleteElasticIP: !s.detachOnly})
	if err != nil {
		return fmt.Errorf("delete server: %w", err)
	}

	return nil
}

func (s *serverDeleteCommand) CompleteArg(cmd *cobra.Command, args []string, toComplete string) (
	[]string,
	cobra.ShellCompDirective,
) {
	if len(args) == 0 {
		return completeServer(cmd.Context(), toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func (s *serverDeleteCommand) Build(app commands.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete SERVER",
		Short: "Delete server",
		Long:  "Deletes a compute server.",
		Example: commands.FormatExamples(fmt.Sprintf(`
      # Delete a server and elastic ips attached to it
      %[1]s compute server delete my-server
      
      # Delete a server, but keep elastic ips
      %[1]s compute server delete my-server --detach-only
		`, app.Name)),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: s.CompleteArg,
		RunE:              s.Run,
	}

	cmd.Flags().BoolVar(&s.force, "force", false, "forces deletion of the server without asking for confirmation")
	cmd.Flags().BoolVar(&s.detachOnly, "detach-only", false, "specifies whether elastic ips should only be detached without getting deleted")

	return cmd
}

func completeServer(ctx context.Context, term string) ([]string, cobra.ShellCompDirective) {
	servers, err := compute.ServerService().List(ctx, core.CursorAll)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	filtered := filter.Find(servers, term)

	names := make([]string, len(filtered))
	for i, server := range filtered {
		names[i] = server.Name
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func findServer(ctx context.Context, term string) (compute.Server, error) {
	servers, err := compute.ServerService().List(ctx, core.CursorAll)
	if err != nil {
		return compute.Server{}, fmt.Errorf("fetch servers: %w", err)
	}

	server, err := filter.FindOne(servers, term)
	if err != nil {
		return compute.Server{}, fmt.Errorf("find server: %w", err)
	}

	return server, nil
}

const specialChars = "~!@#$%^&*_-+=`|\\(){}[]:;\"'<>,.?/"

func checkWindowsPassword(password string) error {
	// https://docs.microsoft.com/en-us/windows/security/threat-protection/security-policy-settings/password-must-meet-complexity-requirements

	// 1. Passwords may not contain the user's samAccountName (Account Name) value or entire displayName (Full
	//	  Name value). Both checks aren't case-sensitive.
	if strings.Contains(strings.ToLower(password), "administrator") {
		return fmt.Errorf("windows user password cannot contain the username")
	}

	// 2. The password contains characters from three of the following categories:
	//    	- Uppercase letters of European languages (A through Z, with diacritic marks, Greek and Cyrillic
	//   		characters)
	//   	- Lowercase letters of European languages (a through z, with diacritic marks, Greek and Cyrillic
	//   		characters)
	//   	- Base 10 digits (0 through 9)
	//		- Non-alphanumeric characters (special characters): (~!@#$%^&*_-+=`|\(){}[]:;"'<>,.?/) Currency symbols such
	//			as the Euro or British Pound aren't counted as special characters for this policy setting.
	// 		- Any Unicode character that's categorized as an alphabetic character but isn't uppercase or lowercase. This
	//			group includes Unicode characters from Asian languages. (NOTE: not implemented)
	hasUppercase := false
	hasLowercase := false
	hasDigit := false
	hasSpecial := false
	count := 0

	for _, char := range password {
		if unicode.IsUpper(char) && !hasUppercase {
			hasUppercase = true
			count++
		}

		if unicode.IsLower(char) && !hasLowercase {
			hasLowercase = true
			count++
		}

		if char >= '0' && char <= '9' && !hasDigit {
			hasDigit = true
			count++
		}

		if strings.ContainsRune(specialChars, char) && !hasSpecial {
			hasSpecial = true
			count++
		}
	}

	if count < 3 {
		return fmt.Errorf("windows user password must contain at least 3 of the following categories: uppercase letters, lowercase letters, digits, and non-alphanumeric characters")
	}

	return nil
}
