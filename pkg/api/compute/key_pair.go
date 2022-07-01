package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

type KeyPair compute.KeyPair

func (k KeyPair) Keys() []string {
	return []string{fmt.Sprint(k.ID), k.Name, k.Fingerprint}
}

func (k KeyPair) Columns() []string {
	return []string{"id", "name", "fingerprint"}
}

func (k KeyPair) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":          k.ID,
		"name":        k.Name,
		"fingerprint": k.Fingerprint,
	}
}

type KeyPairService struct {
	delegate compute.KeyPairService
}

func NewKeyPairService(client goclient.Client) KeyPairService {
	return KeyPairService{
		delegate: compute.NewKeyPairService(client),
	}
}

func (k KeyPairService) List(ctx context.Context) ([]KeyPair, error) {
	res, err := k.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]KeyPair, len(res.Items))
	for i, item := range res.Items {
		items[i] = KeyPair(item)
	}

	return items, nil
}

type KeyPairCreate = compute.KeyPairCreate

func (k KeyPairService) Create(ctx context.Context, data KeyPairCreate) (KeyPair, error) {
	res, err := k.delegate.Create(ctx, data)
	if err != nil {
		return KeyPair{}, err
	}

	return KeyPair(res), nil
}

func (k KeyPairService) Delete(ctx context.Context, id int) error {
	return k.delegate.Delete(ctx, id)
}
