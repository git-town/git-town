package fixture

import (
	"os"

	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/git-town/git-town/v17/test/asserts"
	"github.com/git-town/git-town/v17/test/commands"
	"github.com/git-town/git-town/v17/test/testruntime"
)

// An empty Git repo for testing.
// This is useful for scenarios that require testing the behavior of Git Town in a fresh repository.
type Fresh struct {
	Dir string
}

// NewFresh provides a Fresh instance in the given directory.
//
// The repo has no branches.
func NewFresh(dir string) Fresh {
	binPath := binPath(dir)
	devRepoPath := developerRepoPath(dir)
	// create the "developer" repo
	err := os.MkdirAll(devRepoPath, 0o744)
	asserts.NoError(err)
	// initialize the repo in the folder
	devRepo := testruntime.InitializeNoInitialCommit(devRepoPath, dir, binPath)
	devRepo.RemoveUnnecessaryFiles()
	return Fresh{dir}
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
