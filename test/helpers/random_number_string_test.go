package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	testData := []int{0, 1, 10, 100}
	for _, input := range testData {
		assert.Equal(t, input, len(RandomNumberString(input)))
	}
}
