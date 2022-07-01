package compute

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

const (
	ImageCategoryLinux   = "linux"
	ImageCategoryWindows = "windows"
	ImageCategoryBSD     = "bsd"
)

type Image struct {
	compute.Image

	Availability []common.Location `json:"-"`
}

func (i Image) IsWindows() bool {
	return i.Category == ImageCategoryWindows
}

func (i Image) AvailableAt(location common.Location) bool {
	for _, available := range i.AvailableLocations {
		if available == location.ID {
			return true
		}
	}

	return false
}

func (i Image) String() string {
	return fmt.Sprint(i.OperatingSystem, " ", i.Version)
}

func (i Image) Keys() []string {
	return []string{
		strconv.FormatInt(int64(i.ID), 10),
		i.Key,
		fmt.Sprint(i.OperatingSystem, " ", i.Version),
	}
}

func (i Image) Columns() []string {
	return []string{"id", "operating system", "version", "key", "type", "availability", "license"}
}

func (i Image) Values() map[string]interface{} {
	availabilityBuf := &strings.Builder{}
	for idx, location := range i.Availability {
		if idx != 0 {
			availabilityBuf.WriteString(", ")
		}

		availabilityBuf.WriteString(location.Name)
	}

	licenseBuf := &strings.Builder{}
	for idx, license := range i.RequiredLicenses {
		licenseBuf.WriteString(fmt.Sprintf("%s (%s)", license.Name, common.Product(license).PricePerHour()))

		if idx+1 < len(i.RequiredLicenses) {
			licenseBuf.WriteString(", ")
		}
	}

	imageType := regexp.MustCompile("(?:_|^)[a-z]").ReplaceAllStringFunc(i.Type, func(s string) string {
		if s[0] == '_' {
			s = " " + s[1:]
		}

		return strings.ToUpper(s)
	})

	return map[string]interface{}{
		"id":               i.ID,
		"operating system": i.OperatingSystem,
		"version":          i.Version,
		"key":              i.Key,
		"type":             imageType,
		"availability":     availabilityBuf.String(),
		"license":          licenseBuf.String(),
	}
}

func Images(ctx context.Context, client goclient.Client) ([]Image, error) {
	locations, err := common.Locations(ctx, client)
	if err != nil {
		return nil, err
	}

	res, err := compute.NewImageService(client).List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Image, len(res.Items))
	for idx, item := range res.Items {
		availability := make([]common.Location, len(item.AvailableLocations))
		for idx, id := range item.AvailableLocations {
			for _, location := range locations {
				if location.ID == id {
					availability[idx] = location
					break
				}
			}
		}

		items[idx] = Image{
			Image:        item,
			Availability: availability,
		}
	}

	return items, nil
}
