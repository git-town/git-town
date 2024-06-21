package fixture

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"

	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/cache"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/test/asserts"
	"github.com/git-town/git-town/v14/test/commands"
	"github.com/git-town/git-town/v14/test/datatable"
	"github.com/git-town/git-town/v14/test/filesystem"
	testgit "github.com/git-town/git-town/v14/test/git"
	"github.com/git-town/git-town/v14/test/helpers"
	"github.com/git-town/git-town/v14/test/subshell"
	"github.com/git-town/git-town/v14/test/testruntime"
)

// Fixture is a complete Git environment for a Cucumber scenario.
type Fixture struct {
	// CoworkerRepo is the optional Git repository that is locally checked out at the coworker machine.
	CoworkerRepo OptionP[testruntime.TestRuntime]

	// DevRepo is the Git repository that is locally checked out at the developer machine.
	DevRepo testruntime.TestRuntime

	// Dir defines the local folder in which this Fixture is stored.
	// This folder also acts as the HOME directory for tests using this Fixture.
	// It contains the global Git configuration to use in this test.
	Dir string

	// OriginRepo is the Git repository that simulates the origin repo (on GitHub).
	// If this value is nil, the current test setup has no origin.
	OriginRepo OptionP[testruntime.TestRuntime]

	// SecondWorktree is the directory that contains an additional workspace.
	// If this value is nil, the current test setup has no additional workspace.
	SecondWorktree OptionP[testruntime.TestRuntime]

	// SubmoduleRepo is the Git repository that simulates an external repo used as a submodule.
	// If this value is nil, the current test setup uses no submodules.
	SubmoduleRepo OptionP[testruntime.TestRuntime]

	// UpstreamRepo is the optional Git repository that contains the upstream for this environment.
	UpstreamRepo OptionP[testruntime.TestRuntime]
}

// CloneFixture provides a Fixture instance in the given directory,
// containing a copy of the given Fixture.
func CloneFixture(original Fixture, dir string) Fixture {
	filesystem.CopyDirectory(original.Dir, dir)
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
	devRepo.ConnectTrackingBranch(gitdomain.NewLocalBranchName("main"))
	return Fixture{
		CoworkerRepo:   NoneP[testruntime.TestRuntime](),
		DevRepo:        devRepo,
		Dir:            dir,
		OriginRepo:     SomeP(&originRepo),
		SecondWorktree: NoneP[testruntime.TestRuntime](),
		SubmoduleRepo:  NoneP[testruntime.TestRuntime](),
		UpstreamRepo:   NoneP[testruntime.TestRuntime](),
	}
}

// NewStandardFixture provides a Fixture in the given directory,
// fully populated as a standardized setup for scenarios.
//
// The origin repo has the initial branch checked out.
// Git repos cannot receive pushes of the currently checked out branch
// because that will change files in the current workspace.
// The tests don't use the initial branch.
func NewStandardFixture(dir string) Fixture {
	// create the folder
	// create the fixture
	originPath := originRepoPath(dir)
	binPath := binPath(dir)
	devRepoPath := developerRepoPath(dir)
	// create the origin repo
	err := os.MkdirAll(originPath, 0o744)
	if err != nil {
		log.Fatalf("cannot create directory %q: %v", originPath, err)
	}
	// initialize the repo in the folder
	originRepo := testruntime.Initialize(originPath, dir, binPath)
	err = originRepo.Run("git", "commit", "--allow-empty", "-m", "initial commit")
	if err != nil {
		log.Fatalf("cannot initialize origin directory at %q: %v", originPath, err)
	}
	err = originRepo.Run("git", "branch", "main", "initial")
	if err != nil {
		log.Fatalf("cannot initialize origin directory at %q: %v", originPath, err)
	}
	// clone the "developer" repo
	devRepo := testruntime.Clone(originRepo.TestRunner, devRepoPath)
	initializeWorkspace(&devRepo)
	devRepo.RemoveUnnecessaryFiles()
	originRepo.RemoveUnnecessaryFiles()
	return Fixture{
		CoworkerRepo:   NoneP[testruntime.TestRuntime](),
		DevRepo:        devRepo,
		Dir:            dir,
		OriginRepo:     SomeP(&originRepo),
		SecondWorktree: NoneP[testruntime.TestRuntime](),
		SubmoduleRepo:  NoneP[testruntime.TestRuntime](),
		UpstreamRepo:   NoneP[testruntime.TestRuntime](),
	}
}

// AddCoworkerRepo adds a coworker repository.
func (self *Fixture) AddCoworkerRepo() {
	coworkerRepo := testruntime.Clone(self.OriginRepo.GetOrPanic().TestRunner, self.coworkerRepoPath())
	self.CoworkerRepo = SomeP(&coworkerRepo)
	initializeWorkspace(&coworkerRepo)
	coworkerRepo.Verbose = self.DevRepo.Verbose
}

func (self *Fixture) AddSecondWorktree(branch gitdomain.LocalBranchName) {
	workTreePath := filepath.Join(self.Dir, "development_worktree")
	self.DevRepo.AddWorktree(workTreePath, branch)
	runner := subshell.TestRunner{
		BinDir:     self.DevRepo.BinDir,
		Verbose:    self.DevRepo.Verbose,
		HomeDir:    self.DevRepo.HomeDir,
		WorkingDir: workTreePath,
	}
	gitCommands := git.Commands{
		CurrentBranchCache: &cache.LocalBranchWithPrevious{},
		RemotesCache:       &cache.Remotes{},
	}
	self.SecondWorktree = SomeP(&testruntime.TestRuntime{
		TestCommands: commands.TestCommands{
			TestRunner: &runner,
			Commands:   &gitCommands,
			Config:     self.DevRepo.Config,
		},
		Config: self.DevRepo.Config,
	})
}

// AddSubmodule adds a submodule repository.
func (self *Fixture) AddSubmoduleRepo() {
	err := os.MkdirAll(self.submoduleRepoPath(), 0o744)
	if err != nil {
		log.Fatalf("cannot create directory %q: %v", self.submoduleRepoPath(), err)
	}
	submoduleRepo := testruntime.Initialize(self.submoduleRepoPath(), self.Dir, self.binPath())
	submoduleRepo.MustRun("git", "config", "--global", "protocol.file.allow", "always")
	self.SubmoduleRepo = SomeP(&submoduleRepo)
}

// AddUpstream adds an upstream repository.
func (self *Fixture) AddUpstream() {
	upstreamRepo := testruntime.Clone(self.DevRepo.TestRunner, filepath.Join(self.Dir, gitdomain.RemoteUpstream.String()))
	self.UpstreamRepo = SomeP(&upstreamRepo)
	self.DevRepo.AddRemote(gitdomain.RemoteUpstream, upstreamRepo.WorkingDir)
}

// Branches provides a tabular list of all branches in this Fixture.
func (self *Fixture) Branches() datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("REPOSITORY", "BRANCHES")
	mainBranch := gitdomain.NewLocalBranchName("main")
	initialBranch := gitdomain.NewLocalBranchName("initial")
	localBranches, err := self.DevRepo.LocalBranches()
	asserts.NoError(err)
	localBranchesJoined := localBranches.RemoveWorktreeMarkers().Remove(initialBranch).Hoist(mainBranch).Join(", ")
	originRepo, hasOriginRepo := self.OriginRepo.Get()
	if !hasOriginRepo {
		result.AddRow("local", localBranchesJoined)
		return result
	}
	originBranches, err := originRepo.LocalBranches()
	asserts.NoError(err)
	originBranchesJoined := originBranches.Remove(initialBranch).Hoist(mainBranch).Join(", ")
	if localBranchesJoined == originBranchesJoined {
		result.AddRow("local, origin", localBranchesJoined)
	} else {
		result.AddRow("local", localBranchesJoined)
		result.AddRow("origin", originBranchesJoined)
	}
	return result
}

// CommitTable provides a table for all commits in this Git environment containing only the given fields.
func (self Fixture) CommitTable(fields []string) datatable.DataTable {
	builder := datatable.NewCommitTableBuilder()
	localCommits := self.DevRepo.Commits(fields, gitdomain.NewLocalBranchName("main"))
	builder.AddMany(localCommits, "local")
	if coworkerRepo, hasCoworkerRepo := self.CoworkerRepo.Get(); hasCoworkerRepo {
		coworkerCommits := coworkerRepo.Commits(fields, gitdomain.NewLocalBranchName("main"))
		builder.AddMany(coworkerCommits, "coworker")
	}
	if originRepo, hasOriginRepo := self.OriginRepo.Get(); hasOriginRepo {
		originCommits := originRepo.Commits(fields, gitdomain.NewLocalBranchName("main"))
		builder.AddMany(originCommits, gitdomain.RemoteOrigin.String())
	}
	if upstreamRepo, hasUpstreamRepo := self.UpstreamRepo.Get(); hasUpstreamRepo {
		upstreamCommits := upstreamRepo.Commits(fields, gitdomain.NewLocalBranchName("main"))
		builder.AddMany(upstreamCommits, "upstream")
	}
	if secondWorkTree, hasSecondWorkTree := self.SecondWorktree.Get(); hasSecondWorkTree {
		secondWorktreeCommits := secondWorkTree.Commits(fields, gitdomain.NewLocalBranchName("main"))
		builder.AddMany(secondWorktreeCommits, "worktree")
	}
	return builder.Table(fields)
}

// CreateCommits creates the commits described by the given Gherkin table in this Git repository.
func (self *Fixture) CreateCommits(commits []testgit.Commit) {
	for _, commit := range commits {
		switch {
		case commit.Locations.Matches(testgit.LocationCoworker):
			self.CoworkerRepo.GetOrPanic().CreateCommit(commit)
		case commit.Locations.Matches(testgit.LocationLocal):
			self.DevRepo.CreateCommit(commit)
		case commit.Locations.Matches(testgit.LocationLocal, testgit.LocationOrigin):
			self.DevRepo.CreateCommit(commit)
			self.DevRepo.PushBranch()
		case commit.Locations.Matches(testgit.LocationOrigin):
			self.OriginRepo.GetOrPanic().CreateCommit(commit)
		case commit.Locations.Matches(testgit.LocationUpstream):
			self.UpstreamRepo.GetOrPanic().CreateCommit(commit)
		default:
			log.Fatalf("unknown commit locations %q", commit.Locations)
		}
	}
	// after setting up the commits, check out the "initial" branch in the origin repo so that we can git-push to it.
	if originRepo, hasOriginRepo := self.OriginRepo.Get(); hasOriginRepo {
		originRepo.CheckoutBranch(gitdomain.NewLocalBranchName("initial"))
	}
}

// CreateTags creates tags from the given gherkin table.
func (self Fixture) CreateTags(table *godog.Table) {
	columnNames := helpers.TableFields(table)
	if columnNames[0] != "NAME" && columnNames[1] != "LOCATION" {
		log.Fatalf("tag table must have columns NAME and LOCATION")
	}
	for _, row := range table.Rows[1:] {
		name := row.Cells[0].Value
		location := row.Cells[1].Value
		switch location {
		case "local":
			self.DevRepo.CreateTag(name)
		case "origin":
			self.OriginRepo.GetOrPanic().CreateTag(name)
		default:
			log.Fatalf("tag table LOCATION must be 'local' or 'origin'")
		}
	}
}

func (self Fixture) Delete() {
	os.RemoveAll(self.Dir)
}

// TagTable provides a table for all tags in this Git environment.
func (self Fixture) TagTable() datatable.DataTable {
	builder := datatable.NewTagTableBuilder()
	localTags := self.DevRepo.Tags()
	builder.AddMany(localTags, "local")
	if originRepo, hasOriginRepo := self.OriginRepo.Get(); hasOriginRepo {
		originTags := originRepo.Tags()
		builder.AddMany(originTags, gitdomain.RemoteOrigin.String())
	}
	return builder.Table()
}

// binPath provides the full path of the folder containing the test tools for this Fixture.
func (self *Fixture) binPath() string {
	return binPath(self.Dir)
}

func binPath(rootDir string) string {
	return filepath.Join(rootDir, "bin")
}

// coworkerRepoPath provides the full path to the Git repository with the given name.
func (self Fixture) coworkerRepoPath() string {
	return filepath.Join(self.Dir, "coworker")
}

func developerRepoPath(rootDir string) string {
	return filepath.Join(rootDir, "developer")
}

func initializeWorkspace(repo *testruntime.TestRuntime) {
	asserts.NoError(repo.Config.SetMainBranch(gitdomain.NewLocalBranchName("main")))
	asserts.NoError(repo.Config.SetPerennialBranches(gitdomain.LocalBranchNames{}))
	repo.MustRun("git", "checkout", "main")
	// NOTE: the developer repos receives the initial branch from origin
	//       but we don't want it here because it isn't used in tests.
	repo.MustRun("git", "branch", "-d", "initial")
}

func originRepoPath(rootDir string) string {
	return filepath.Join(rootDir, gitdomain.RemoteOrigin.String())
}

// submoduleRepoPath provides the full path to the Git repository with the given name.
func (self Fixture) submoduleRepoPath() string {
	return filepath.Join(self.Dir, "submodule")
}
