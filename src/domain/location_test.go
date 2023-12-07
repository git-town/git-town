package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestLocation(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		location := domain.NewLocation("branch-1")
		have, err := json.MarshalIndent(location, "", "  ")
		must.NoError(t, err)
		want := `"branch-1"`
		must.EqOp(t, want, string(have))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"branch-1"`
		have := domain.EmptyLocation()
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := domain.NewLocation("branch-1")
		must.EqOp(t, want, have)
	})
}
