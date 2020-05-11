package flow

import (
	"context"
	"fmt"
	"net/http"
)

type ImageService interface {
	List(ctx context.Context, options PaginationOptions) ([]*Image, *Response, error)
	Get(ctx context.Context, id Id) (*Image, *Response, error)
}

type Image struct {
	Id                 Id         `json:"id"`
	OperatingSystem    string     `json:"os"`
	Version            string     `json:"version"`
	Key                string     `json:"key"`
	Category           string     `json:"category"`
	Type               string     `json:"type"`
	MinRootDiskSize    int        `json:"min_root_disk_size"`
	Sorting            int        `json:"sorting"`
	RequiredLicenses   []*Product `json:"required_licenses"`
	AvailableLocations []Id       `json:"available_locations"`
}

func (i *Image) String() string {
	return fmt.Sprintf("%s %s", i.OperatingSystem, i.Version)
}

func (i *Image) AvailableAt(location *Location) bool {
	for _, available := range i.AvailableLocations {
		if available == location.Id {
			return true
		}
	}
	return false
}

type imageService struct {
	client *Client
}

func (s *imageService) List(ctx context.Context, options PaginationOptions) ([]*Image, *Response, error) {
	path, err := addOptions("/v3/entities/images", options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil, FlagNoAuthentication)
	if err != nil {
		return nil, nil, err
	}

	var val []*Image

	res, err := s.client.Do(req, &val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}

func (s *imageService) Get(ctx context.Context, id Id) (*Image, *Response, error) {
	path := fmt.Sprintf("/v3/entities/images/%d", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil, FlagNoAuthentication)
	if err != nil {
		return nil, nil, err
	}

	val := &Image{}

	res, err := s.client.Do(req, val)
	if err != nil {
		return nil, nil, err
	}

	return val, res, nil
}
