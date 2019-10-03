package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLongestString(t *testing.T) {
	testData := map[int][]string{
		5: []string{"one", "two", "three"},
	}
	for expected, input := range testData {
		actual := LongestString(input)
		assert.Equal(t, expected, actual)
	}
}
