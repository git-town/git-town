package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, want, have)
	}
}
