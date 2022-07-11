package compute

import (
	"context"
	"fmt"

	"github.com/flowswiss/cli/v2/internal/commands"
	"github.com/flowswiss/cli/v2/pkg/api/compute"
	"github.com/flowswiss/cli/v2/pkg/filter"
)

func findSecurityGroup(ctx context.Context, term string) (compute.SecurityGroup, error) {
	securityGroups, err := compute.NewSecurityGroupService(commands.Config.Client).List(ctx)
	if err != nil {
		return compute.SecurityGroup{}, fmt.Errorf("fetch security groups: %w", err)
	}

	securityGroup, err := filter.FindOne(securityGroups, term)
	if err != nil {
		return compute.SecurityGroup{}, fmt.Errorf("find security group: %w", err)
	}

	return securityGroup, nil
}
