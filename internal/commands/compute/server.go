package compute

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/cobra"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/common"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/console"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func ServerCommand() *cobra.Command {
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
		`, commands.Name)),
	}

	commands.Add(cmd, &serverListCommand{}, &serverCreateCommand{}, &serverUpdateCommand{}, &serverDeleteCommand{})
	cmd.AddCommand(NetworkInterfaceCommand())

	return cmd
}

type serverListCommand struct {
	filter string
}

func (s *serverListCommand) Run(cmd *cobra.Command, args []string) error {
	items, err := compute.NewServerService(commands.Config.Client).List(cmd.Context())
	if err != nil {
		return err
	}

	if len(s.filter) != 0 {
		items = filter.Find(items, s.filter)
	}

	return commands.PrintStdout(items)
}

func (s *serverListCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all server",
		Long:  "Prints a table of all compute servers belonging to the current organization.",
		RunE:  s.Run,
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
	privateIP        net.IP
	keyPair          string
	password         string
	cloudInitFile    string
	attachExternalIP bool
}

func (s *serverCreateCommand) Run(cmd *cobra.Command, args []string) error {
	locations, err := common.Locations(cmd.Context(), commands.Config.Client)
	if err != nil {
		return fmt.Errorf("fetch locations: %w", err)
	}

	location, err := filter.FindOne(locations, s.location)
	if err != nil {
		return fmt.Errorf("find location: %w", err)
	}

	images, err := compute.Images(cmd.Context(), commands.Config.Client)
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

	products, err := common.ProductsByType(cmd.Context(), commands.Config.Client, common.ProductTypeComputeServer)
	if err != nil {
		return fmt.Errorf("fetch products: %w", err)
	}

	product, err := filter.FindOne(products, s.product)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}

	networkID := 0
	if s.network != "" {
		networks, err := compute.NewNetworkService(commands.Config.Client).List(cmd.Context())
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

		if !cidr.Contains(s.privateIP) {
			return fmt.Errorf("private ip %s is not in network %s", s.privateIP, network.CIDR)
		}

		networkID = network.ID
	}

	privateIP := ""
	if len(s.privateIP) != 0 {
		privateIP = s.privateIP.String()
	}

	keyPairID := 0
	if s.keyPair != "" {
		keyPairs, err := compute.NewKeyPairService(commands.Config.Client).List(cmd.Context())
		if err != nil {
			return fmt.Errorf("fetch key pairs: %w", err)
		}

		keyPair, err := filter.FindOne(keyPairs, s.keyPair)
		if err != nil {
			return fmt.Errorf("find key pair: %w", err)
		}

		keyPairID = keyPair.ID
	}

	if !image.IsWindows() && keyPairID == 0 {
		return fmt.Errorf("key pair is required for non-windows images")
	}

	password := s.password
	if image.IsWindows() {
		if len(password) == 0 {
			password, err = console.Password(commands.Stderr, "Windows User Password", checkWindowsPassword)
			if err != nil {
				return fmt.Errorf("read user password: %w", err)
			}
		}

		if err = checkWindowsPassword(password); err != nil {
			return fmt.Errorf("check user password: %w", err)
		}
	}

	cloudInit := ""
	if len(s.cloudInitFile) != 0 {
		data, err := os.ReadFile(s.cloudInitFile)
		if err != nil {
			return fmt.Errorf("read cloud init file: %w", err)
		}

		cloudInit = base64.StdEncoding.EncodeToString(data)
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

	ordering, err := compute.NewServerService(commands.Config.Client).Create(cmd.Context(), data)
	if err != nil {
		return fmt.Errorf("create server: %w", err)
	}

	progress := console.NewProgress("Creating server")
	go progress.Display(commands.Stderr)

	err = common.WaitForOrder(cmd.Context(), commands.Config.Client, ordering)
	if err != nil {
		progress.Complete("Order failed")

		return fmt.Errorf("wait for order: %w", err)
	}

	progress.Complete("Order completed")

	return nil
}

func (s *serverCreateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new server",
		Long:  "Creates a new compute server.",
		Example: commands.FormatExamples(fmt.Sprintf(`
			# Create a new ubuntu server
	  		%[1]s compute server create --name my-server --location ALP1 --image linux-ubuntu-20.04-lts --product b1.4x8 --key-pair my-keypair

			# Create a new windows server
	  		%[1]s compute server create --name my-server --location ALP1 --image microsoft-windows-server-2019 --product b1.2x8
		`, commands.Name)), // TODO select correct image names
		RunE: s.Run,
	}

	cmd.Flags().StringVarP(&s.name, "name", "n", "", "name of the new server (required)")
	cmd.Flags().StringVarP(&s.location, "location", "l", "", "location of the server (required)")
	cmd.Flags().StringVarP(&s.image, "image", "i", "", "operating system image to use for the new server (required)")
	cmd.Flags().StringVarP(&s.product, "product", "p", "", "product to use for the new server (required)")
	cmd.Flags().StringVar(&s.network, "network", "", "network in which the first network interface should be created")
	cmd.Flags().IPVar(&s.privateIP, "private-ip", nil, "ip address of the server in the selected network")
	cmd.Flags().StringVar(&s.keyPair, "key-pair", "", "ssh key-pair for connecting to the server (required if image is linux)")
	cmd.Flags().StringVar(&s.password, "windows-password", "", "password for the windows admin user  (required if image is windows)")
	cmd.Flags().StringVar(&s.cloudInitFile, "cloud-init", "", "cloud init script to customize creation of the server")
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

	data := compute.ServerUpdate{
		Name: s.name,
	}

	server, err = compute.NewServerService(commands.Config.Client).Update(cmd.Context(), server.ID, data)
	if err != nil {
		return fmt.Errorf("update server: %w", err)
	}

	return commands.PrintStdout(server)
}

func (s *serverUpdateCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update SERVER",
		Short: "Update server",
		Long:  "Updates a compute server.",
		Args:  cobra.ExactArgs(1),
		RunE:  s.Run,
	}

	cmd.Flags().StringVar(&s.name, "name", "", "new name of the server")

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

	if !s.force {
		if !console.Confirm(commands.Stderr, fmt.Sprintf("Are you sure you want to delete server %q?", server.Name)) {
			commands.Stderr.Println("aborted.")
			return nil
		}
	}

	err = compute.NewServerService(commands.Config.Client).Delete(cmd.Context(), server.ID, !s.detachOnly)
	if err != nil {
		return fmt.Errorf("delete server: %w", err)
	}

	return nil
}

func (s *serverDeleteCommand) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete SERVER",
		Short: "Delete server",
		Long:  "Deletes a compute server.",
		Example: commands.FormatExamples(fmt.Sprintf(`
			# Delete a server and elastic ips attached to it
			%[1]s compute server delete my-server

			# Delete a server, but keep elastic ips
			%[1]s compute server delete my-server --detach-only
		`, commands.Name)),
		Args: cobra.ExactArgs(1),
		RunE: s.Run,
	}

	cmd.Flags().BoolVar(&s.force, "force", false, "forces deletion of the server without asking for confirmation")
	cmd.Flags().BoolVar(&s.detachOnly, "detach-only", false, "specifies whether elastic ips should only be detached without getting deleted")

	return cmd
}

func findServer(ctx context.Context, term string) (compute.Server, error) {
	servers, err := compute.NewServerService(commands.Config.Client).List(ctx)
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
