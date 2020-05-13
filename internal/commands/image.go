package commands

import (
	"context"
	"fmt"
	"github.com/flowswiss/cli/pkg/flow"
)

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
