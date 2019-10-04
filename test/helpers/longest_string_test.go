package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLongestStringLength(t *testing.T) {
	testData := map[int][]string{
		5: {"one", "two", "three"},
		0: {},
	}
	for expected, input := range testData {
		assert.Equal(t, expected, LongestStringLength(input))
	}
}
