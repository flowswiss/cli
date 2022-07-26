package objectstorage

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/objectstorage"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type Credential objectstorage.Credential

func (c Credential) String() string {
	return c.Endpoint
}

func (c Credential) Keys() []string {
	keys := []string{fmt.Sprint(c.ID), c.Endpoint}
	keys = append(keys, common.Location(c.Location).Keys()...)
	return keys
}

func (c Credential) Columns() []string {
	return []string{"id", "location", "endpoint", "access key", "secret key"}
}

func (c Credential) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":         c.ID,
		"location":   common.Location(c.Location),
		"endpoint":   c.Endpoint,
		"access key": c.AccessKey,
		"secret key": c.SecretKey,
	}
}

type CredentialService struct {
	delegate objectstorage.CredentialService
}

func NewCredentialService(client goclient.Client) CredentialService {
	return CredentialService{
		delegate: objectstorage.NewCredentialService(client),
	}
}

func (c CredentialService) List(ctx context.Context) ([]Credential, error) {
	res, err := c.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]Credential, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = Credential(item)
	}

	return items, nil
}
