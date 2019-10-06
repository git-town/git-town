package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueString(t *testing.T) {
	assert.NotEqual(t, "", UniqueString())
}
