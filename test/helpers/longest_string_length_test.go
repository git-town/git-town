package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestLongestStringLength(t *testing.T) {
	t.Parallel()
	tests := map[int][]string{
		5: {"one", "two", "three"},
		0: {""},
	}
	for expected, input := range tests {
		assert.Equal(t, expected, helpers.LongestStringLength(input))
	}
}
