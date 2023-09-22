package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestLocation(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		location := domain.NewLocation("branch-1")
		have, err := json.MarshalIndent(location, "", "  ")
		assert.Nil(t, err)
		want := `"branch-1"`
		assert.Equal(t, want, string(have))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"branch-1"`
		have := domain.Location{}
		err := json.Unmarshal([]byte(give), &have)
		assert.Nil(t, err)
		want := domain.NewLocation("branch-1")
		assert.Equal(t, want, have)
	})
}
