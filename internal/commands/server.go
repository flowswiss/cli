package commands

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/flowswiss/cli/internal/commands/dto"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
	"time"
)

const maxFilterDepth = 2

var (
	serverCommand = &cobra.Command{
		Use:   "server",
		Short: "Manage your compute server",
	}

	serverListCommand = &cobra.Command{
		Use:   "list",
		Short: "List all server",
		RunE:  listVMs,
	}

	serverCreateCommand = &cobra.Command{
		Use:     "create",
		Short:   "Create new server",
		PreRunE: preCreateVM,
		RunE:    createVM,
	}

	serverUpdateCommand = &cobra.Command{
		Use:   "edit",
		Short: "Rename server",
		RunE:  updateVM,
	}

	serverDeleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "Delete server",
		RunE:  deleteVM,
	}
)

var (
	name             string
	locationFilter   string
	imageFilter      string
	productFilter    string
	networkFilter    string
	privateIp        string
	keyPairFilter    string
	password         string
	cloudInitFile    string
	attachExternalIp bool

	location *flow.Location
	image    *flow.Image
	product  *flow.Product
	network  *flow.Network
	keyPair  *flow.KeyPair
)

func init() {
	serverCommand.AddCommand(serverListCommand)
	serverCommand.AddCommand(serverCreateCommand)
	serverCommand.AddCommand(serverUpdateCommand)
	serverCommand.AddCommand(serverDeleteCommand)

	serverCreateCommand.Flags().StringVarP(&name, "name", "n", "", "Name of the new server")
	serverCreateCommand.Flags().StringVarP(&locationFilter, "location", "l", "", "Location of the server")
	serverCreateCommand.Flags().StringVarP(&imageFilter, "image", "i", "", "Name of the new server")
	serverCreateCommand.Flags().StringVarP(&productFilter, "product", "p", "", "Name of the new server")
	serverCreateCommand.Flags().StringVar(&networkFilter, "network", "", "Name of the new server")
	serverCreateCommand.Flags().StringVar(&privateIp, "private-ip", "", "Name of the new server")
	serverCreateCommand.Flags().StringVar(&keyPairFilter, "key-pair", "", "Name of the new server")
	serverCreateCommand.Flags().StringVar(&password, "windows-password", "", "Name of the new server")
	serverCreateCommand.Flags().StringVar(&cloudInitFile, "cloud-init", "", "Name of the new server")
	serverCreateCommand.Flags().BoolVar(&attachExternalIp, "attach-external-ip", true, "Name of the new server")
}

func listVMs(cmd *cobra.Command, args []string) error {
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

func preCreateVM(cmd *cobra.Command, args []string) error {
	// check required flags
	if name == "" {
		return errRequiredFlag("name")
	}

	if locationFilter == "" {
		return errRequiredFlag("location")
	}

	if imageFilter == "" {
		return errRequiredFlag("image")
	}

	if productFilter == "" {
		return errRequiredFlag("product")
	}

	// search location
	locations, _, err := client.Location.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	loc, err := findOne(locations, locationFilter, maxFilterDepth)
	if err != nil {
		return fmt.Errorf("location: %v", err)
	}
	location = loc.(*flow.Location)

	// search image
	images, _, err := client.Image.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	img, err := findOne(images, imageFilter, maxFilterDepth)
	if err != nil {
		return fmt.Errorf("image: %v", err)
	}
	image = img.(*flow.Image)

	// search product
	products, _, err := client.Product.ListByType(context.Background(), flow.PaginationOptions{NoFilter: 1}, "compute-engine-vm")
	if err != nil {
		return err
	}

	prod, err := findOne(products, productFilter, maxFilterDepth)
	if err != nil {
		return fmt.Errorf("product: %v", err)
	}
	product = prod.(*flow.Product)

	// search network
	if networkFilter != "" {
		networks, _, err := client.Network.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
		if err != nil {
			return err
		}

		net, err := findOne(networks, networkFilter, maxFilterDepth)
		if err != nil {
			return fmt.Errorf("network: %v", err)
		}
		network = net.(*flow.Network)
	}

	// search key pair
	if keyPairFilter != "" {
		keyPairs, _, err := client.KeyPair.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
		if err != nil {
			return err
		}

		key, err := findOne(keyPairs, keyPairFilter, maxFilterDepth)
		if err != nil {
			return fmt.Errorf("key-pair: %v", err)
		}
		keyPair = key.(*flow.KeyPair)
	}

	return nil
}

func createVM(cmd *cobra.Command, args []string) error {
	if !product.AvailableAt(location) {
		return fmt.Errorf("product is not available at the selected location")
	}

	if product.Type.Id != 4 {
		return fmt.Errorf("product is not a compute vm product")
	}

	if !image.AvailableAt(location) {
		return fmt.Errorf("image is not available at the selected location")
	}

	if strings.ToLower(image.Category) == "windows" && password == "" {
		return fmt.Errorf("windows images require password")
	}

	if strings.ToLower(image.Category) == "windows" && cloudInitFile != "" {
		return fmt.Errorf("windows images are not allowed to take cloud init scripts")
	}

	if strings.ToLower(image.Category) == "linux" && keyPair == nil {
		return fmt.Errorf("linux images require key pair")
	}

	cloudInit := ""
	if cloudInitFile != "" {
		data, err := ioutil.ReadFile(cloudInitFile)
		if err != nil {
			return err
		}

		cloudInit = base64.StdEncoding.EncodeToString(data)
	}

	networkId := flow.Id(0)
	if network != nil {
		networkId = network.Id
	}

	keyPairId := flow.Id(0)
	if keyPair != nil {
		keyPairId = keyPair.Id
	}

	data := &flow.ServerCreate{
		Name:             name,
		LocationId:       location.Id,
		ImageId:          image.Id,
		ProductId:        product.Id,
		AttachExternalIp: attachExternalIp,
		NetworkId:        networkId,
		PrivateIp:        privateIp,
		KeyPairId:        keyPairId,
		Password:         password,
		CloudInit:        cloudInit,
	}

	ordering, _, err := client.Server.Create(context.Background(), data)
	if err != nil {
		return err
	}

	id, err := ordering.Id()
	if err != nil {
		return err
	}

	done := make(chan string)
	sync := make(chan struct{})
	defer func() {
		<-sync
	}()

	go stderr.Progress("creating server", done, sync)

	for {
		order, _, err := client.Order.Get(context.Background(), id)
		if err != nil || order.Status.Id == 4 {
			done <- "creation failed"
			return err
		}

		if order.Status.Id == 3 {
			done <- "server created successfully"
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}

func updateVM(cmd *cobra.Command, args []string) error {
	return nil
}

func deleteVM(cmd *cobra.Command, args []string) error {
	return nil
}
