package scrub

import (
	"context"

	"github.com/danielpickens/hercules/pkg/config"
	"github.com/danielpickens/hercules/types"
	"github.com/derailed/danielpickens/hercules/cache"
	"github.com/derailed/danielpickens/hercules/issues"
	"github.com/derailed/danielpickens/hercules/sanitize"
)

// ClusterRole represents a ClusterRole scruber.
type ClusterRole struct {
	client types.Connection
	*config.Config
	*issues.Collector

	*cache.ClusterRole
	*cache.ClusterRoleBinding
	*cache.RoleBinding
}

// NewClusterRole return a new ClusterRole scruber.
func NewClusterRole(ctx context.Context, c *Cache, codes *issues.Codes) Sanitizer {
	crb := ClusterRole{
		client:    c.factory.Client(),
		Config:    c.config,
		Collector: issues.NewCollector(codes, c.config),
	}

	var err error
	crb.ClusterRole, err = c.clusterroles()
	if err != nil {
		crb.AddErr(ctx, err)
	}

	crb.ClusterRoleBinding, err = c.clusterrolebindings()
	if err != nil {
		crb.AddCode(ctx, 402, err)
	}

	crb.RoleBinding, err = c.rolebindings()
	if err != nil {
		crb.AddErr(ctx, err)
	}

	return &crb
}

// Sanitize all available ClusterRoles.
func (c *ClusterRole) Sanitize(ctx context.Context) error {
	return sanitize.NewClusterRole(c.Collector, c).Sanitize(ctx)
}
