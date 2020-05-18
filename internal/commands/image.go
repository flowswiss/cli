package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/internal/commands/dto"
	"github.com/flowswiss/cli/pkg/flow"
	"github.com/spf13/cobra"
	"sort"
)

var (
	imageCommand = &cobra.Command{
		Use:   "image",
		Short: "",
	}

	imageListCommand = &cobra.Command{
		Use:   "list",
		Short: "List all images",
		RunE:  listImage,
	}
)

func init() {
	imageCommand.AddCommand(imageListCommand)

	imageListCommand.Flags().String(flagLocation, "", "filter for availability at location")
}

func findImage(filter string) (*flow.Image, error) {
	images, _, err := client.Image.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	image, err := findOne(images, filter, 2)
	if err != nil {
		return nil, fmt.Errorf("image: %v", err)
	}

	return image.(*flow.Image), nil
}

func listImage(cmd *cobra.Command, args []string) error {
	images, _, err := client.Image.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	locations, _, err := client.Location.List(context.Background(), flow.PaginationOptions{NoFilter: 1})
	if err != nil {
		return err
	}

	locationFilter, err := cmd.Flags().GetString(flagLocation)
	if err != nil {
		return err
	}

	var location *flow.Location
	if locationFilter != "" {
		location, err = findLocation(locationFilter)
		if err != nil {
			return err
		}
	}

	sort.Sort(imageBySorting{images})

	var displayable []*dto.Image
	for _, image := range images {
		if location != nil && !image.AvailableAt(location) {
			continue
		}

		displayable = append(displayable, &dto.Image{Image: image, Locations: locations})
	}

	return display(displayable)
}

type imageBySorting struct {
	Images []*flow.Image
}

func (s imageBySorting) Len() int {
	return len(s.Images)
}

func (s imageBySorting) Swap(i, j int) {
	s.Images[i], s.Images[j] = s.Images[j], s.Images[i]
}

func (s imageBySorting) Less(i, j int) bool {
	return s.Images[i].Sorting < s.Images[j].Sorting
}
