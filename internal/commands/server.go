package commands

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/flowswiss/cli/internal/commands/dto"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/flowswiss/cli/pkg/output"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const maxFilterDepth = 2

const (
	flagName             = "name"
	flagLocation         = "location"
	flagImage            = "image"
	flagProduct          = "product"
	flagNetwork          = "network"
	flagPrivateIp        = "private-ip"
	flagKeyPair          = "key-pair"
	flagWindowsPassword  = "windows-password"
	flagCloudInit        = "cloud-init"
	flagAttachExternalIp = "attach-external-ip"

	flagForce      = "force"
	flagDetachOnly = "detach-only"
)

var (
	serverCommand = &cobra.Command{
		Use:   "server",
		Short: "Manage your compute server",
	}

	serverListCommand = &cobra.Command{
		Use:   "list",
		Short: "List all server",
		RunE:  listServer,
	}

	serverCreateCommand = &cobra.Command{
		Use:   "create",
		Short: "Create new server",
		RunE:  createServer,
	}

	serverUpdateCommand = &cobra.Command{
		Use:   "edit",
		Short: "Rename server",
		RunE:  updateServer,
	}

	serverDeleteCommand = &cobra.Command{
		Use:   "delete <server>",
		Short: "Delete server",
		Args:  cobra.ExactArgs(1),
		RunE:  deleteServer,
	}
)

func init() {
	serverCommand.AddCommand(serverListCommand)
	serverCommand.AddCommand(serverCreateCommand)
	serverCommand.AddCommand(serverUpdateCommand)
	serverCommand.AddCommand(serverDeleteCommand)

	serverCreateCommand.Flags().StringP(flagName, "n", "", "name of the new server (required)")
	serverCreateCommand.Flags().StringP(flagLocation, "l", "", "location of the server (required)")
	serverCreateCommand.Flags().StringP(flagImage, "i", "", "operating system image to use for the new server (required)")
	serverCreateCommand.Flags().StringP(flagProduct, "p", "", "product to use for the new server (required)")
	serverCreateCommand.Flags().String(flagNetwork, "", "network in which the first network interface should be created")
	serverCreateCommand.Flags().String(flagPrivateIp, "", "ip address of the server in the selected network")
	serverCreateCommand.Flags().String(flagKeyPair, "", "ssh key-pair for connecting to the server (required if image is linux)")
	serverCreateCommand.Flags().String(flagWindowsPassword, "", "password for the windows admin user  (required if image is windows)")
	serverCreateCommand.Flags().String(flagCloudInit, "", "cloud init script to customize creation of the server")
	serverCreateCommand.Flags().Bool(flagAttachExternalIp, true, "whether to attach an elastic ip to the server")

	handleError(serverCreateCommand.MarkFlagRequired(flagName))
	handleError(serverCreateCommand.MarkFlagRequired(flagLocation))
	handleError(serverCreateCommand.MarkFlagRequired(flagImage))
	handleError(serverCreateCommand.MarkFlagRequired(flagProduct))

	serverDeleteCommand.Flags().Bool(flagForce, false, "forces deletion of the server without asking for confirmation")
	serverDeleteCommand.Flags().Bool(flagDetachOnly, false, "specifies whether elastic ips should only be detached without getting deleted")
}

func findServer(filter string) (*flow.Server, error) {
	servers, _, err := client.Server.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	srv, err := findOne(servers, filter, 2)
	if err != nil {
		return nil, fmt.Errorf("server: %v", err)
	}

	return srv.(*flow.Server), nil
}

func listServer(cmd *cobra.Command, args []string) error {
	server, _, err := client.Server.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	var displayable []*dto.Server
	for _, server := range server {
		displayable = append(displayable, &dto.Server{Server: server})
	}

	return display(displayable)
}

func createServer(cmd *cobra.Command, args []string) error {
	data, err := parseCreateServerData(cmd)
	if err != nil {
		return err
	}

	ordering, _, err := client.Server.Create(context.Background(), data)
	if err != nil {
		return err
	}

	id, err := ordering.Id()
	if err != nil {
		return err
	}

	progress := output.NewProgress("creating server")
	go progress.Display(stderr)

	for {
		order, _, err := client.Order.Get(context.Background(), id)
		if err != nil || order.Status.Id == 4 {
			progress.Complete("creation failed")
			return err
		}

		if order.Status.Id == 3 {
			progress.Complete("server created successfully")
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}

func updateServer(cmd *cobra.Command, args []string) error {
	return nil
}

func deleteServer(cmd *cobra.Command, args []string) error {
	server, err := findServer(args[0])
	if err != nil {
		return err
	}

	force, err := cmd.Flags().GetBool(flagForce)
	if err != nil {
		return err
	}

	detachOnly, err := cmd.Flags().GetBool(flagDetachOnly)
	if err != nil {
		return err
	}

	if !force {
		stderr.Printf("Are you sure you want to delete the server \"%v\" (y/N): ", server)

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return err
		}

		if strings.ToLower(input) != "y\n" {
			return nil
		}
	}

	elasticIps, _, err := client.ServerAttachment.ListAttachedElasticIps(context.Background(), server.Id, flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	for _, elasticIp := range elasticIps {
		_, err := client.ServerAttachment.DetachElasticIp(context.Background(), server.Id, elasticIp.Id)
		if err != nil {
			return err
		}

		if !detachOnly {
			_, err := client.ElasticIp.Delete(context.Background(), elasticIp.Id)
			if err != nil {
				return err
			}
		}
	}

	_, err = client.Server.Delete(context.Background(), server.Id)
	if err != nil {
		return err
	}

	return nil
}

func parseCreateServerData(cmd *cobra.Command) (*flow.ServerCreate, error) {
	var err error
	result := &flow.ServerCreate{}

	// validate name
	result.Name, err = cmd.Flags().GetString(flagName)
	if err != nil {
		return nil, err
	}

	// validate location
	locationFilter, err := cmd.Flags().GetString(flagLocation)
	if err != nil {
		return nil, err
	}

	location, err := findLocation(locationFilter)
	if err != nil {
		return nil, err
	}

	result.LocationId = location.Id

	// validate image
	imageFilter, err := cmd.Flags().GetString(flagImage)
	if err != nil {
		return nil, err
	}

	image, err := findImage(imageFilter)
	if err != nil {
		return nil, err
	}

	if !image.AvailableAt(location) {
		return nil, fmt.Errorf("image is not available at the selected location")
	}

	result.ImageId = image.Id

	// validate product
	productFilter, err := cmd.Flags().GetString(flagProduct)
	if err != nil {
		return nil, err
	}

	product, err := findProduct(productFilter)
	if err != nil {
		return nil, err
	}

	if !product.AvailableAt(location) {
		return nil, fmt.Errorf("product is not available at the selected location")
	}

	disk := product.FindItem(3)
	if product.Type.Id != 4 || disk == nil {
		return nil, fmt.Errorf("product is not a compute vm product")
	}

	if image.MinRootDiskSize > disk.Amount {
		return nil, fmt.Errorf("the %s image requires at least %d GB of storage", image.Key, image.MinRootDiskSize)
	}

	result.ProductId = product.Id

	// validate cloud init
	cloudInitFile, err := cmd.Flags().GetString(flagCloudInit)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(image.Category) == "windows" && cloudInitFile != "" {
		return nil, fmt.Errorf("windows images are not allowed to take cloud init scripts")
	}

	if cloudInitFile != "" {
		data, err := ioutil.ReadFile(cloudInitFile)
		if err != nil {
			return nil, err
		}

		result.CloudInit = base64.StdEncoding.EncodeToString(data)
	}

	// validate network
	networkFilter, err := cmd.Flags().GetString(flagNetwork)
	if err != nil {
		return nil, err
	}

	if networkFilter != "" {
		network, err := findNetwork(networkFilter)
		if err != nil {
			return nil, err
		}

		result.NetworkId = network.Id
	}

	// validate private ip
	result.PrivateIp, err = cmd.Flags().GetString(flagPrivateIp)
	if err != nil {
		return nil, err
	}

	// validate key pair
	keyPairFilter, err := cmd.Flags().GetString(flagKeyPair)

	if keyPairFilter != "" {
		keyPair, err := findKeyPair(keyPairFilter)
		if err != nil {
			return nil, err
		}

		result.KeyPairId = keyPair.Id
	}

	if strings.ToLower(image.Category) == "linux" && result.KeyPairId == 0 {
		return nil, fmt.Errorf("linux images require key pair")
	}

	// validate windows password
	result.Password, err = cmd.Flags().GetString(flagWindowsPassword)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(image.Category) == "windows" && result.Password == "" {
		return nil, fmt.Errorf("windows images require password")
	}

	// validate external ip
	result.AttachExternalIp, err = cmd.Flags().GetBool(flagAttachExternalIp)
	if err != nil {
		return nil, err
	}

	return result, nil
}
