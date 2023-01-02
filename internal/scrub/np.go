package scrub

import (
	"context"

	"github.com/danielpickens/hercules/internal/cache"
	"github.com/danielpickens/hercules/internal/issues"
	"github.com/danielpickens/hercules/internal/sanitize"
	"github.com/danielpickens/hercules/pkg/config"
	"github.com/danielpickens/hercules/types"
)

// NetworkPolicy represents a NetworkPolicy scruber.
type NetworkPolicy struct {
	*issues.Collector
	*cache.NetworkPolicy
	*cache.Namespace
	*cache.Pod
	*config.Config

	client types.Connection
}

// NewNetworkPolicy return a new NetworkPolicy scruber.
func NewNetworkPolicy(ctx context.Context, c *Cache, codes *issues.Codes) Sanitizer {
	n := NetworkPolicy{
		client:    c.factory.Client(),
		Config:    c.config,
		Collector: issues.NewCollector(codes, c.config),
	}

	var err error
	n.NetworkPolicy, err = c.networkpolicies()
	if err != nil {
		n.AddErr(ctx, err)
	}

	n.Namespace, err = c.namespaces()
	if err != nil {
		n.AddCode(ctx, 402, err)
	}

	n.Pod, err = c.pods()
	if err != nil {
		n.AddErr(ctx, err)
	}

	return &n
}

// Sanitize all available NetworkPolicys.
func (n *NetworkPolicy) Sanitize(ctx context.Context) error {
	return sanitize.NewNetworkPolicy(n.Collector, n).Sanitize(ctx)
}
