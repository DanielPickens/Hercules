package scrub

import (
	"context"

	"github.com/danielpickens/hercules/internal/cache"
	"github.com/danielpickens/hercules/internal/issues"
	"github.com/danielpickens/hercules/internal/sanitize"
	"github.com/danielpickens/hercules/pkg/config"
	"github.com/danielpickens/hercules/types"
)

// Ingress represents a Ingress scruber.
type Ingress struct {
	*issues.Collector
	*cache.Ingress
	*config.Config

	client types.Connection
}

// NewIngress return a new Ingress scruber.
func NewIngress(ctx context.Context, c *Cache, codes *issues.Codes) Sanitizer {
	d := Ingress{
		client:    c.factory.Client(),
		Config:    c.config,
		Collector: issues.NewCollector(codes, c.config),
	}

	var err error
	d.Ingress, err = c.ingresses()
	if err != nil {
		d.AddErr(ctx, err)
	}

	return &d
}

// Sanitize all available Ingresss.
func (i *Ingress) Sanitize(ctx context.Context) error {
	return sanitize.NewIngress(i.Collector, i).Sanitize(ctx)
}
