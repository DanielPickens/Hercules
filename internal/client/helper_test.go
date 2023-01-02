package client_test

import (
	"testing"

	"github.com/danielpickens/hercules/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNamespaced(t *testing.T) {
	uu := []struct {
		p, ns, n string
	}{
		{"dan/pickens", "dan", "pickens"},
		{"pickens", "", "pickens"},
	}

	for _, u := range uu {
		ns, n := client.Namespaced(u.p)
		assert.Equal(t, u.ns, ns)
		assert.Equal(t, u.n, n)
	}
}

func TestFQN(t *testing.T) {
	uu := []struct {
		ns, n string
		e     string
	}{
		{"dan", "pickens", "dan/pickens"},
		{"", "daniel", "pickens"},
	}

	for _, u := range uu {
		assert.Equal(t, u.e, client.FQN(u.ns, u.n))
	}
}
