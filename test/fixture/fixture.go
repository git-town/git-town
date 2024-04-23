package fixture

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/cache"
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
	CoworkerRepo *testruntime.TestRuntime `exhaustruct:"optional"`

	// DevRepo is the Git repository that is locally checked out at the developer machine.
	DevRepo testruntime.TestRuntime `exhaustruct:"optional"`

	// Dir defines the local folder in which this Fixture is stored.
	// This folder also acts as the HOME directory for tests using this Fixture.
	// It contains the global Git configuration to use in this test.
	Dir string

	// OriginRepo is the Git repository that simulates the origin repo (on GitHub).
	// If this value is nil, the current test setup has no origin.
	OriginRepo *testruntime.TestRuntime `exhaustruct:"optional"`

	// SecondWorktree is the directory that contains an additional workspace.
	// If this value is nil, the current test setup has no additional workspace.
	SecondWorktree *testruntime.TestRuntime `exhaustruct:"optional"`

	// SubmoduleRepo is the Git repository that simulates an external repo used as a submodule.
	// If this value is nil, the current test setup uses no submodules.
	SubmoduleRepo *testruntime.TestRuntime `exhaustruct:"optional"`

	// UpstreamRepo is the optional Git repository that contains the upstream for this environment.
	UpstreamRepo *testruntime.TestRuntime `exhaustruct:"optional"`
}

// CloneFixture provides a Fixture instance in the given directory,
// containing a copy of the given Fixture.
func CloneFixture(original Fixture, dir string) Fixture {
	filesystem.CopyDirectory(original.Dir, dir)
	binDir := filepath.Join(dir, "bin")
	originDir := filepath.Join(dir, gitdomain.RemoteOrigin.String())
	originRepo := testruntime.New(originDir, dir, "")
	developerDir := filepath.Join(dir, "developer")
	devRepo := testruntime.New(developerDir, dir, binDir)
	result := Fixture{
		DevRepo:    devRepo,
		Dir:        dir,
		OriginRepo: &originRepo,
	}
	// Since we copied the files from the memoized directory,
	// we have to set the "origin" remote to the copied origin repo here.
	result.DevRepo.MustRun("git", "remote", "remove", gitdomain.RemoteOrigin.String())
	result.DevRepo.AddRemote(gitdomain.RemoteOrigin, result.originRepoPath())
	result.DevRepo.Fetch()
	// and connect the main branches again
	result.DevRepo.ConnectTrackingBranch(gitdomain.NewLocalBranchName("main"))
	return result
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
	gitEnv := Fixture{Dir: dir}
	// create the origin repo
	err := os.MkdirAll(gitEnv.originRepoPath(), 0o744)
	if err != nil {
		log.Fatalf("cannot create directory %q: %v", gitEnv.originRepoPath(), err)
	}
	// initialize the repo in the folder
	originRepo := testruntime.Initialize(gitEnv.originRepoPath(), gitEnv.Dir, gitEnv.binPath())
	err = originRepo.RunMany([][]string{
		{"git", "commit", "--allow-empty", "-m", "initial commit"},
		{"git", "branch", "main", "initial"},
	})
	if err != nil {
		log.Fatalf("cannot initialize origin directory at %q: %v", gitEnv.originRepoPath(), err)
	}
	gitEnv.OriginRepo = &originRepo
	// clone the "developer" repo
	gitEnv.DevRepo = testruntime.Clone(originRepo.TestRunner, gitEnv.developerRepoPath())
	gitEnv.initializeWorkspace(&gitEnv.DevRepo)
	gitEnv.DevRepo.RemoveUnnecessaryFiles()
	gitEnv.OriginRepo.RemoveUnnecessaryFiles()
	return gitEnv
}

// AddCoworkerRepo adds a coworker repository.
func (self *Fixture) AddCoworkerRepo() {
	coworkerRepo := testruntime.Clone(self.OriginRepo.TestRunner, self.coworkerRepoPath())
	self.CoworkerRepo = &coworkerRepo
	self.initializeWorkspace(self.CoworkerRepo)
	self.CoworkerRepo.Verbose = self.DevRepo.Verbose
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
	backendCommands := git.BackendCommands{
		Runner:             &runner,
		DryRun:             false,
		CurrentBranchCache: &cache.LocalBranchWithPrevious{},
		RemotesCache:       &cache.Remotes{},
	}
	self.SecondWorktree = &testruntime.TestRuntime{
		TestCommands: commands.TestCommands{
			TestRunner:      &runner,
			BackendCommands: &backendCommands,
		},
		Backend: backendCommands,
		Config:  config,
	}
}

// AddSubmodule adds a submodule repository.
func (self *Fixture) AddSubmoduleRepo() {
	err := os.MkdirAll(self.submoduleRepoPath(), 0o744)
	if err != nil {
		log.Fatalf("cannot create directory %q: %v", self.submoduleRepoPath(), err)
	}
	submoduleRepo := testruntime.Initialize(self.submoduleRepoPath(), self.Dir, self.binPath())
	submoduleRepo.MustRunMany([][]string{
		{"git", "config", "--global", "protocol.file.allow", "always"},
		{"git", "commit", "--allow-empty", "-m", "initial commit"},
	})
	self.SubmoduleRepo = &submoduleRepo
}

// AddUpstream adds an upstream repository.
func (self *Fixture) AddUpstream() {
	repo := testruntime.Clone(self.DevRepo.TestRunner, filepath.Join(self.Dir, gitdomain.RemoteUpstream.String()))
	self.UpstreamRepo = &repo
	self.DevRepo.AddRemote(gitdomain.RemoteUpstream, self.UpstreamRepo.WorkingDir)
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
	if self.OriginRepo == nil {
		result.AddRow("local", localBranchesJoined)
		return result
	}
	originBranches, err := self.OriginRepo.LocalBranches()
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
	if self.CoworkerRepo != nil {
		coworkerCommits := self.CoworkerRepo.Commits(fields, gitdomain.NewLocalBranchName("main"))
		builder.AddMany(coworkerCommits, "coworker")
	}
	if self.OriginRepo != nil {
		originCommits := self.OriginRepo.Commits(fields, gitdomain.NewLocalBranchName("main"))
		builder.AddMany(originCommits, gitdomain.RemoteOrigin.String())
	}
	if self.UpstreamRepo != nil {
		upstreamCommits := self.UpstreamRepo.Commits(fields, gitdomain.NewLocalBranchName("main"))
		builder.AddMany(upstreamCommits, "upstream")
	}
	if self.SecondWorktree != nil {
		secondWorktreeCommits := self.SecondWorktree.Commits(fields, gitdomain.NewLocalBranchName("main"))
		builder.AddMany(secondWorktreeCommits, "worktree")
	}
	return builder.Table(fields)
}

// CreateCommits creates the commits described by the given Gherkin table in this Git repository.
func (self *Fixture) CreateCommits(commits []testgit.Commit) {
	for _, commit := range commits {
		switch {
		case commit.Locations.Matches(testgit.LocationCoworker):
			self.CoworkerRepo.CreateCommit(commit)
		case commit.Locations.Matches(testgit.LocationLocal):
			self.DevRepo.CreateCommit(commit)
		case commit.Locations.Matches(testgit.LocationLocal, testgit.LocationOrigin):
			self.DevRepo.CreateCommit(commit)
			self.DevRepo.PushBranch()
		case commit.Locations.Matches(testgit.LocationOrigin):
			self.OriginRepo.CreateCommit(commit)
		case commit.Locations.Matches(testgit.LocationUpstream):
			self.UpstreamRepo.CreateCommit(commit)
		default:
			log.Fatalf("unknown commit locations %q", commit.Locations)
		}
	}
	// after setting up the commits, check out the "initial" branch in the origin repo so that we can git-push to it.
	if self.OriginRepo != nil {
		self.OriginRepo.CheckoutBranch(gitdomain.NewLocalBranchName("initial"))
	}
}

// CreateTags creates tags from the given gherkin table.
func (self Fixture) CreateTags(table *messages.PickleStepArgument_PickleTable) {
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
			self.OriginRepo.CreateTag(name)
		default:
			log.Fatalf("tag table LOCATION must be 'local' or 'origin'")
		}
	}
}

// TagTable provides a table for all tags in this Git environment.
func (self Fixture) TagTable() datatable.DataTable {
	builder := datatable.NewTagTableBuilder()
	localTags := self.DevRepo.Tags()
	builder.AddMany(localTags, "local")
	if self.OriginRepo != nil {
		originTags := self.OriginRepo.Tags()
		builder.AddMany(originTags, gitdomain.RemoteOrigin.String())
	}
	return builder.Table()
}

// binPath provides the full path of the folder containing the test tools for this Fixture.
func (self *Fixture) binPath() string {
	return filepath.Join(self.Dir, "bin")
}

// coworkerRepoPath provides the full path to the Git repository with the given name.
func (self Fixture) coworkerRepoPath() string {
	return filepath.Join(self.Dir, "coworker")
}

// developerRepoPath provides the full path to the Git repository with the given name.
func (self Fixture) developerRepoPath() string {
	return filepath.Join(self.Dir, "developer")
}

func (self Fixture) initializeWorkspace(repo *testruntime.TestRuntime) {
	asserts.NoError(repo.Config.SetMainBranch(gitdomain.NewLocalBranchName("main")))
	asserts.NoError(repo.Config.SetPerennialBranches(gitdomain.LocalBranchNames{}))
	repo.MustRunMany([][]string{
		{"git", "checkout", "main"},
		// NOTE: the developer repos receives the initial branch from origin
		//       but we don't want it here because it isn't used in tests.
		{"git", "branch", "-d", "initial"},
	})
}

// originRepoPath provides the full path to the Git repository with the given name.
func (self Fixture) originRepoPath() string {
	return filepath.Join(self.Dir, gitdomain.RemoteOrigin.String())
}

// submoduleRepoPath provides the full path to the Git repository with the given name.
func (self Fixture) submoduleRepoPath() string {
	return filepath.Join(self.Dir, "submodule")
}
