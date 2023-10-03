package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v9/test/helpers"
	"github.com/shoenig/test"
)

func TestLongestStringLength(t *testing.T) {
	t.Parallel()
	tests := map[int][]string{
		5: {"one", "two", "three"},
		0: {""},
	}
	for expected, input := range tests {
		test.EqOp(t, expected, helpers.LongestStringLength(input))
	}
}
