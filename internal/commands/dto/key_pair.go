package dto

import "github.com/flowswiss/cli/pkg/flow"

type KeyPair struct {
	*flow.KeyPair
}

func (k *KeyPair) Columns() []string {
	return []string{"id", "name", "fingerprint"}
}

func (k *KeyPair) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":          k.Id,
		"name":        k.Name,
		"fingerprint": k.Fingerprint,
	}
}
