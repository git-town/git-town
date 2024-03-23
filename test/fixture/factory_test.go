package fixture_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v13/test/fixture"
	"github.com/shoenig/test/must"
)

func TestFixtureFactory(t *testing.T) {
	t.Parallel()

	t.Run("NewFixtureFactory", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		_ = fixture.NewFactory(dir)
		memoizedPath := filepath.Join(dir, "memoized")
		_, err := os.Stat(memoizedPath)
		must.False(t, os.IsNotExist(err))
	})

	t.Run("CreateFixture", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		gm := fixture.NewFactory(dir)
		result := gm.CreateFixture("foo")
		_, err := os.Stat(result.DevRepo.WorkingDir)
		must.False(t, os.IsNotExist(err))
	})
}
