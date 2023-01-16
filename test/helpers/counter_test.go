package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v7/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestUniqueString(t *testing.T) {
	t.Parallel()
	assert.NotEqual(t, "", helpers.Counter{})
}
