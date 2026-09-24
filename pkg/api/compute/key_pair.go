package compute

import (
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/generic"
	"github.com/flowswiss/goclient/v2/compute"
	"github.com/flowswiss/goclient/v2/core"
)

type (
	KeyPairCreate = compute.KeyPairCreateReq
	KeyPairList   = core.Cursor
	KeyPairDelete = compute.KeyPairDeleteReq
)

type GenericKeyPairService = generic.CFD[
	compute.KeyPair,
	KeyPair,
	KeyPairCreate,
	KeyPairList,
	KeyPairDelete,
]

func KeyPairService() GenericKeyPairService {
	return generic.NewCFD[
		compute.KeyPair,
		KeyPair,
		KeyPairCreate,
		KeyPairList,
		KeyPairDelete,
	](commands.Client.Compute.KeyPair, func(pair compute.KeyPair) KeyPair {
		return KeyPair(pair)
	})
}

type KeyPair compute.KeyPair

func (k KeyPair) String() string {
	return k.Name
}

func (k KeyPair) Keys() []string {
	return []string{fmt.Sprint(k.ID), k.Name, k.Fingerprint}
}

func (k KeyPair) Columns() []string {
	return []string{"id", "name", "fingerprint"}
}

func (k KeyPair) Values() map[string]any {
	return map[string]any{
		"id":          k.ID,
		"name":        k.Name,
		"fingerprint": k.Fingerprint,
	}
}
