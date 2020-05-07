package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestUniqueString(t *testing.T) {
	assert.NotEqual(t, "", helpers.UniqueString())
}
