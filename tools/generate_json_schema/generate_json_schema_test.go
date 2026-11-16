package main

import (
	"fmt"
	"testing"

	"github.com/shoenig/test/must"
)

func TestCamelToKebab(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"CreatePrototypeBranches": "create-prototype-branches",
		"Single":                  "single",
		"SyncStrategy":            "sync-strategy",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := CamelToKebab(give)
			must.Eq(t, want, have)
		})
	}
}
