package scrub

import (
	"context"

	"github.com/danielpickens/hercules/internal/issues"
	"github.com/danielpickens/hercules/pkg/config"
)

// Sanitizer represents a resource sanitizer.
type Sanitizer interface {
	Collector
	Sanitize(context.Context) error
}

// Collector collects sanitization issues.
type Collector interface {
	MaxSeverity(res string) config.Level
	Outcome() issues.Outcome
}
