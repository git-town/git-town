package fixture

import (
	"github.com/git-town/git-town/v22/internal/test/commands"
	"github.com/git-town/git-town/v22/internal/test/testruntime"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// An empty Git repo for testing.
// This is useful for scenarios that require testing the behavior of Git Town in a fresh repository.
type Fresh struct {
	Dir string
}

// allows using this fresh environment as a Fixture
func (self Fresh) AsFixture() Fixture {
	binDir := binPath(self.Dir)
	developerDir := developerRepoPath(self.Dir)
	devRepo := testruntime.New(developerDir, self.Dir, binDir)
	return Fixture{
		CoworkerRepo:   MutableNone[commands.TestCommands](),
		DevRepo:        MutableSome(&devRepo),
		Dir:            self.Dir,
		OriginRepo:     MutableNone[commands.TestCommands](),
		SecondWorktree: MutableNone[commands.TestCommands](),
		SubmoduleRepo:  MutableNone[commands.TestCommands](),
		UpstreamRepo:   MutableNone[commands.TestCommands](),
	}
}
