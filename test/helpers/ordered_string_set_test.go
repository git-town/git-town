package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedStringSet(t *testing.T) {
	set1 := NewOrderedStringSet("one")
	set2 := set1.Add("two")
	set2 = set2.Add("two")
	assert.Equal(t, []string{"one", "two"}, set2.Slice())
	assert.Equal(t, "one, two", set2.String())
}
