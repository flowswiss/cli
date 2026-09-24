package objectstorage

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/core"
	"github.com/flowswiss/goclient/v2/objectstorage"

	"github.com/flowswiss/cli/v2/pkg/api/common"
)

type (
	CredentialList = core.Cursor
)

type GenericCredentialService = generic.ListService[CredentialList, objectstorage.Credential, Credential]

func CredentialService() GenericCredentialService {
	client := commands.Client.ObjectStorage.Credential
	cast := func(credential objectstorage.Credential) Credential {
		return Credential(credential)
	}

	return generic.NewList[CredentialList, objectstorage.Credential, Credential](client, cast)
}

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

func (c Credential) Values() map[string]any {
	return map[string]any{
		"id":         c.ID,
		"location":   common.Location(c.Location),
		"endpoint":   c.Endpoint,
		"access key": c.AccessKey,
		"secret key": c.SecretKey,
	}
}
