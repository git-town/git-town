package fixture

import (
	"os"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/test/commands"
	"github.com/git-town/git-town/v22/internal/test/filesystem"
	"github.com/git-town/git-town/v22/internal/test/testruntime"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// A fully populated Git repos template for testing.
// This is just the template that can be efficiently cloned.
// To perform Git operations, clone or derive a Fixture from it.
type Memoized struct {
	Dir string
}

// NewMemoized provides a Memoized instance in the given directory.
//
// The origin repo has the initial branch checked out.
// Git repos cannot receive pushes of the currently checked out branch
// because that will change files in the current workspace.
// The tests don't use the initial branch.
func NewMemoized(dir string) Memoized {
	originPath := originRepoPath(dir)
	binPath := binPath(dir)
	devRepoPath := developerRepoPath(dir)
	// create the origin repo
	err := os.MkdirAll(originPath, 0o744)
	asserts.NoError(err)
	// initialize the repo in the folder
	originRepo := testruntime.Initialize(originPath, dir, binPath)
	err = originRepo.Run("git", "branch", "main", "initial")
	asserts.NoError(err)
	// clone the "developer" repo
	devRepo := testruntime.Clone(originRepo.TestRunner, devRepoPath)
	initializeWorkspace(&devRepo)
	devRepo.RemoveUnnecessaryFiles()
	originRepo.RemoveUnnecessaryFiles()
	return Memoized{dir}
}

// allows using this memoized environment as a Fixture
func (self Memoized) AsFixture() Fixture {
	binDir := binPath(self.Dir)
	developerDir := developerRepoPath(self.Dir)
	originDir := originRepoPath(self.Dir)
	originRepo := testruntime.New(originDir, self.Dir, "")
	devRepo := testruntime.New(developerDir, self.Dir, binDir)
	return Fixture{
		CoworkerRepo:   MutableNone[commands.TestCommands](),
		DevRepo:        MutableSome(&devRepo),
		Dir:            self.Dir,
		OriginRepo:     MutableSome(&originRepo),
		SecondWorktree: MutableNone[commands.TestCommands](),
		SubmoduleRepo:  MutableNone[commands.TestCommands](),
		UpstreamRepo:   MutableNone[commands.TestCommands](),
	}
}

// provides a copy of this Memoized in the given directory
func (self Memoized) CloneInto(dir string) Fixture {
	filesystem.CopyDirectory(self.Dir, dir)
	binDir := binPath(dir)
	originDir := originRepoPath(dir)
	originRepo := testruntime.New(originDir, dir, "")
	developerDir := developerRepoPath(dir)
	devRepo := testruntime.New(developerDir, dir, binDir)
	// Since we copied the files from the memoized directory,
	// we have to set the "origin" remote to the copied origin repo here.
	devRepo.MustRun("git", "remote", "remove", gitdomain.RemoteOrigin.String())
	devRepo.AddRemote(gitdomain.RemoteOrigin, originDir)
	devRepo.Fetch()
	// and connect the main branches again
	devRepo.ConnectTrackingBranch("main")
	return Fixture{
		CoworkerRepo:   MutableNone[commands.TestCommands](),
		DevRepo:        MutableSome(&devRepo),
		Dir:            dir,
		OriginRepo:     MutableSome(&originRepo),
		SecondWorktree: MutableNone[commands.TestCommands](),
		SubmoduleRepo:  MutableNone[commands.TestCommands](),
		UpstreamRepo:   MutableNone[commands.TestCommands](),
	}
}
