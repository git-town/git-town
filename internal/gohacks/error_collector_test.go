package gohacks_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/shoenig/test/must"
)

func TestErrorCollector(t *testing.T) {
	t.Parallel()

	t.Run("Check", func(t *testing.T) {
		t.Parallel()
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.ErrorCollector{}
			fc.Check(nil)
			must.Nil(t, fc.Err)
			fc.Check(errors.New("first"))
			fc.Check(errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
		t.Run("indicates whether it received an error", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.ErrorCollector{}
			must.False(t, fc.Check(nil))
			must.True(t, fc.Check(errors.New("")))
			must.True(t, fc.Check(nil))
		})
	})

	t.Run("captures the first error it receives", func(t *testing.T) {
		t.Parallel()
		fc := gohacks.ErrorCollector{}
		fc.Check(errors.New("first"))
		fc.Check(errors.New("second"))
		must.ErrorContains(t, fc.Err, "first")
	})
}
