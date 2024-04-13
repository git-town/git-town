package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestTypeName(t *testing.T) {
	t.Parallel()
	sha := gitdomain.NewSHA("123456")
	tests := map[any]string{
		"hello":                             "string",
		123:                                 "int",
		gitdomain.NewLocalBranchName("foo"): "LocalBranchName", // instance of a struct
		&sha:                                "SHA",             // pointer variable
		nil:                                 "nil",
	}
	for give, want := range tests {
		have := gohacks.TypeName(give)
		must.EqOp(t, want, have)
	}
}
