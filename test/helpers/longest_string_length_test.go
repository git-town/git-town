package helpers_test

import (
	"strings"
	"testing"

	"github.com/git-town/git-town/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestLongestStringLength(t *testing.T) {
	tests := map[int][]string{
		5: {"one", "two", "three"},
		0: {""},
	}
	for expected, input := range tests {
		t.Run(strings.Join(input, "-"), func(t *testing.T) {
			assert.Equal(t, expected, helpers.LongestStringLength(input))
		})
	}
}
