package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestOrderedStringSet(t *testing.T) {
	t.Parallel()
	set1 := helpers.NewOrderedStringSet("one")
	set2 := set1.Add("two")
	set2 = set2.Add("two")
	assert.Equal(t, []string{"one", "two"}, set2.Slice())
	assert.Equal(t, "one, two", set2.String())
}
