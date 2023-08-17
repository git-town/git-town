package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v9/test/helpers"
)

func TestOrderedSet(t *testing.T) {
	t.Parallel()
	t.Run("strings", func(t *testing.T) {
		t.Parallel()
		set := helpers.NewOrderedSet("one", "two")
		set.Add("three")
	})
}
