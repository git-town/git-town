package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
)

func TestDisplayTypes(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		t.Run("All", func(t *testing.T) {
			t.Parallel()
			give := configdomain.DisplayTypes{}
		})
	})
}
