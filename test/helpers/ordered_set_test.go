package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/v9/test/helpers"
)

func TestOrderedSet(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		set := helpers.NewOrderedSet("one", "two")
		set.Add("three")
	})
}
