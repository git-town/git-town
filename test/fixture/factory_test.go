package fixture_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v19/test/fixture"
	"github.com/shoenig/test/must"
)

func TestFixtureFactory(t *testing.T) {
	t.Parallel()

	t.Run("CreateEmptyFixture", func(t *testing.T) {
		t.Parallel()
		factory := fixture.CreateFactory()
		defer factory.Remove()
		result := factory.CreateEmptyFixture("foo")
		_, err := os.Stat(result.DevRepo.GetOrPanic().WorkingDir)
		must.False(t, os.IsNotExist(err))
	})

	t.Run("CreateFixture", func(t *testing.T) {
		t.Parallel()
		factory := fixture.CreateFactory()
		defer factory.Remove()
		result := factory.CreateFixture("foo")
		_, err := os.Stat(result.DevRepo.GetOrPanic().WorkingDir)
		must.False(t, os.IsNotExist(err))
	})
}
