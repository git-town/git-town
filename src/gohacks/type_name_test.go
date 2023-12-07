package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestTypeName(t *testing.T) {
	t.Parallel()
	sha := domain.NewSHA("123456")
	tests := map[any]string{
		"hello":                          "string",
		123:                              "int",
		domain.NewLocalBranchName("foo"): "LocalBranchName", // instance of a struct
		&sha:                             "SHA",             // pointer variable
		nil:                              "nil",
	}
	for give, want := range tests {
		have := gohacks.TypeName(give)
		must.EqOp(t, want, have)
	}
}
