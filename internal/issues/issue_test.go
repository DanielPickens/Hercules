package issues

import (
	"testing"

	"github.com/danielpickens/hercules/internal/client"
	"github.com/danielpickens/hercules/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestIsSubIssues(t *testing.T) {
	uu := map[string]struct {
		i Issue
		e bool
	}{
		"root":  {New(client.NewGVR("dan"), Root, config.WarnLevel, "blah"), false},
		"rootf": {Newf(client.NewGVR("dan"), Root, config.WarnLevel, "blah %s", "blee"), false},
		"sub":   {New(client.NewGVR("dan"), "sub1", config.WarnLevel, "blah"), true},
		"subf":  {Newf(client.NewGVR("dan"), "sub1", config.WarnLevel, "blah %s", "blee"), true},
	}

	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			assert.Equal(t, u.e, u.i.IsSubIssue())
		})
	}
}

func TestBlank(t *testing.T) {
	uu := map[string]struct {
		i Issue
		e bool
	}{
		"blank":    {Issue{}, true},
		"notBlank": {New(client.NewGVR("dan"), Root, config.WarnLevel, "blah"), false},
	}

	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			assert.Equal(t, u.e, u.i.Blank())
		})
	}
}
