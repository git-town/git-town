package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		assert.Equal(t, i, len(RandomString(i)))
	}
}
