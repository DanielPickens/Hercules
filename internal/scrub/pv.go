package scrub

import (
	"context"

	"github.com/danielpickens/hercules/internal/cache"
	"github.com/danielpickens/hercules/internal/issues"
	"github.com/danielpickens/hercules/internal/sanitize"
)

// PersistentVolume represents a PersistentVolume scruber.
type PersistentVolume struct {
	*issues.Collector
	*cache.PersistentVolume
	*cache.Pod
}

// NewPersistentVolume return a new PersistentVolume scruber.
func NewPersistentVolume(ctx context.Context, c *Cache, codes *issues.Codes) Sanitizer {
	p := PersistentVolume{Collector: issues.NewCollector(codes, c.config)}

	var err error
	p.PersistentVolume, err = c.persistentvolumes()
	if err != nil {
		p.AddErr(ctx, err)
	}

	p.Pod, err = c.pods()
	if err != nil {
		p.AddErr(ctx, err)
	}

	return &p
}

// Sanitize all available PersistentVolumes.
func (s *PersistentVolume) Sanitize(ctx context.Context) error {
	return sanitize.NewPersistentVolume(s.Collector, s).Sanitize(ctx)
}
