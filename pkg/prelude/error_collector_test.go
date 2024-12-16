package prelude_test

import (
	"errors"
	"testing"

	"github.com/shoenig/test/must"
)

func TestErrorCollector(t *testing.T) {
	t.Parallel()

	t.Run("captures the first error it receives", func(t *testing.T) {
		t.Parallel()
		fc := ErrorCollector{}
		fc.Check(errors.New("first"))
		fc.Check(errors.New("second"))
		must.ErrorContains(t, fc.Err, "first")
	})

	t.Run("Check", func(t *testing.T) {
		t.Parallel()
		t.Run("captures the first error it receives", func(t *testing.T) {
			fc := ErrorCollector{}
			fc.Check(nil)
			must.Nil(t, fc.Err)
			fc.Check(errors.New("first"))
			fc.Check(errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
		t.Run("indicates whether it received an error", func(t *testing.T) {
			fc := ErrorCollector{}
			must.False(t, fc.Check(nil))
			must.True(t, fc.Check(errors.New("")))
			must.True(t, fc.Check(nil))
		})
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		t.Run("registers the given error", func(t *testing.T) {
			fc := ErrorCollector{}
			fc.Fail("failed %s", "reason")
			must.ErrorContains(t, fc.Err, "failed reason")
		})
	})
}
