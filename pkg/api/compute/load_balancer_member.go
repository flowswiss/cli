package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/goclient"
	"github.com/flowswiss/goclient/compute"
)

type LoadBalancerMember compute.LoadBalancerMember

func (l LoadBalancerMember) Host() string {
	return fmt.Sprintf("%s:%d", l.Address, l.Port)
}

func (l LoadBalancerMember) String() string {
	return l.Name
}

func (l LoadBalancerMember) Keys() []string {
	return []string{fmt.Sprint(l.ID), l.Name, l.Host(), l.Status.Key, l.Status.Name}
}

func (l LoadBalancerMember) Columns() []string {
	return []string{"id", "name", "address", "status"}
}

func (l LoadBalancerMember) Values() map[string]interface{} {
	return map[string]interface{}{
		"id":      l.ID,
		"name":    l.Name,
		"address": l.Host(),
		"status":  l.Status.Name,
	}
}

type LoadBalancerMemberService struct {
	delegate compute.LoadBalancerMemberService
}

func NewLoadBalancerMemberService(client goclient.Client, loadBalancerID, poolID int) LoadBalancerMemberService {
	return LoadBalancerMemberService{
		delegate: compute.NewLoadBalancerMemberService(client, loadBalancerID, poolID),
	}
}

func (l LoadBalancerMemberService) List(ctx context.Context) ([]LoadBalancerMember, error) {
	res, err := l.delegate.List(ctx, goclient.Cursor{NoFilter: 1})
	if err != nil {
		return nil, err
	}

	items := make([]LoadBalancerMember, len(res.Items))
	for idx, item := range res.Items {
		items[idx] = LoadBalancerMember(item)
	}

	return items, nil
}

type LoadBalancerMemberCreate = compute.LoadBalancerMemberCreate

func (l LoadBalancerMemberService) Create(ctx context.Context, data LoadBalancerMemberCreate) (LoadBalancerMember, error) {
	res, err := l.delegate.Create(ctx, data)
	if err != nil {
		return LoadBalancerMember{}, err
	}

	return LoadBalancerMember(res), nil
}

func (l LoadBalancerMemberService) Delete(ctx context.Context, id int) error {
	return l.delegate.Delete(ctx, id)
}
