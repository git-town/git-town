package fixture_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v14/test/fixture"
	"github.com/shoenig/test/must"
)

func TestFixtureFactory(t *testing.T) {
	t.Parallel()

	t.Run("CreateFixture", func(t *testing.T) {
		t.Parallel()
		gm := fixture.CreateFactory()
		defer gm.Remove()
		result := gm.CreateFixture("foo")
		_, err := os.Stat(result.DevRepo.WorkingDir)
		must.False(t, os.IsNotExist(err))
	})
}
