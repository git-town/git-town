package main

import (
	"fmt"
	"testing"

	"github.com/shoenig/test/must"
)

func TestCamelToKebab(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"SyncStrategy":            "sync-strategy",
		"PushNewBranches":         "push-new-branches",
		"CreatePrototypeBranches": "create-prototype-branches",
		"BranchType":              "branch-type",
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := CamelToKebab(give)
			must.Eq(t, want, have)
		})
	}
}
