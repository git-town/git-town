package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedStringSetContains(t *testing.T) {
	set := OrderedStringSet{}
	set.Add("one")
	assert.True(t, set.Contains("one"), "should contain one")
	assert.False(t, set.Contains("two"), "should not contain two")
}

func TestOrderedStringSetSlice(t *testing.T) {
	set := OrderedStringSet{}
	set.Add("one")
	set.Add("two")
	set.Add("one")
	assert.Equal(t, []string{"one", "two"}, set.Slice())
}
