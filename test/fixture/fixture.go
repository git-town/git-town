package fixture

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"

	"github.com/git-town/git-town/v15/internal/git"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks/cache"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/test/asserts"
	"github.com/git-town/git-town/v15/test/commands"
	"github.com/git-town/git-town/v15/test/datatable"
	testgit "github.com/git-town/git-town/v15/test/git"
	"github.com/git-town/git-town/v15/test/helpers"
	"github.com/git-town/git-town/v15/test/subshell"
	"github.com/git-town/git-town/v15/test/testruntime"
)

// Fixture is a complete Git environment for a Cucumber scenario.
type Fixture struct {
	// CoworkerRepo is the optional Git repository that is locally checked out at the coworker machine.
	CoworkerRepo OptionP[testruntime.TestRuntime]

	// DevRepo is the Git repository that is locally checked out at the developer machine.
	DevRepo OptionP[testruntime.TestRuntime]

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

// AddCoworkerRepo adds a coworker repository.
func (self *Fixture) AddCoworkerRepo() {
	coworkerRepo := testruntime.Clone(self.OriginRepo.GetOrPanic().TestRunner, self.coworkerRepoPath())
	self.CoworkerRepo = SomeP(&coworkerRepo)
	initializeWorkspace(&coworkerRepo)
	coworkerRepo.Verbose = self.DevRepo.GetOrPanic().Verbose
}

func (self *Fixture) AddSecondWorktree(branch gitdomain.LocalBranchName) {
	workTreePath := filepath.Join(self.Dir, "development_worktree")
	devRepo := self.DevRepo.GetOrPanic()
	devRepo.AddWorktree(workTreePath, branch)
	runner := subshell.TestRunner{
		BinDir:     devRepo.BinDir,
		Verbose:    devRepo.Verbose,
		HomeDir:    devRepo.HomeDir,
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
			Config:     devRepo.Config,
		},
		Config: devRepo.Config,
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
	devRepo := self.DevRepo.GetOrPanic()
	upstreamRepo := testruntime.Clone(devRepo.TestRunner, filepath.Join(self.Dir, gitdomain.RemoteUpstream.String()))
	upstreamRepo.TestRunner.Verbose = devRepo.Verbose
	self.UpstreamRepo = SomeP(&upstreamRepo)
	devRepo.AddRemote(gitdomain.RemoteUpstream, upstreamRepo.WorkingDir)
}

// Branches provides a tabular list of all branches in this Fixture.
func (self *Fixture) Branches() datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("REPOSITORY", "BRANCHES")
	mainBranch := gitdomain.NewLocalBranchName("main")
	initialBranch := gitdomain.NewLocalBranchName("initial")
	localBranches, err := self.DevRepo.GetOrPanic().LocalBranches()
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
	localCommits := self.DevRepo.GetOrPanic().Commits(fields, gitdomain.NewLocalBranchName("main"))
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
	devRepo := self.DevRepo.GetOrPanic()
	for _, commit := range commits {
		switch {
		case commit.Locations.Matches(testgit.LocationCoworker):
			self.CoworkerRepo.GetOrPanic().CreateCommit(commit)
		case commit.Locations.Matches(testgit.LocationLocal):
			devRepo.CreateCommit(commit)
		case commit.Locations.Matches(testgit.LocationLocal, testgit.LocationOrigin):
			devRepo.CreateCommit(commit)
			devRepo.PushBranch()
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
	devRepo := self.DevRepo.GetOrPanic()
	for _, row := range table.Rows[1:] {
		name := row.Cells[0].Value
		location := row.Cells[1].Value
		switch location {
		case "local":
			devRepo.CreateTag(name)
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
	localTags := self.DevRepo.GetOrPanic().Tags()
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
