package fixture

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/cache"
	"github.com/git-town/git-town/v21/internal/test/commands"
	"github.com/git-town/git-town/v21/internal/test/datatable"
	"github.com/git-town/git-town/v21/internal/test/helpers"
	"github.com/git-town/git-town/v21/internal/test/subshell"
	"github.com/git-town/git-town/v21/internal/test/testgit"
	"github.com/git-town/git-town/v21/internal/test/testruntime"
	"github.com/git-town/git-town/v21/pkg/asserts"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Fixture is a complete Git environment for a Cucumber scenario.
type Fixture struct {
	// CoworkerRepo is the optional Git repository that is locally checked out at the coworker machine.
	CoworkerRepo OptionalMutable[commands.TestCommands]

	// DevRepo is the Git repository that is locally checked out at the developer machine.
	DevRepo OptionalMutable[commands.TestCommands]

	// Dir defines the local folder in which this Fixture is stored.
	// This folder also acts as the HOME directory for tests using this Fixture.
	// It contains the global Git configuration to use in this test.
	Dir string

	// OriginRepo is the Git repository that simulates the origin repo (on GitHub).
	// If this value is nil, the current test setup has no origin.
	OriginRepo OptionalMutable[commands.TestCommands]

	// SecondWorktree is the directory that contains an additional workspace.
	// If this value is nil, the current test setup has no additional workspace.
	SecondWorktree OptionalMutable[commands.TestCommands]

	// SubmoduleRepo is the Git repository that simulates an external repo used as a submodule.
	// If this value is nil, the current test setup uses no submodules.
	SubmoduleRepo OptionalMutable[commands.TestCommands]

	// UpstreamRepo is the optional Git repository that contains the upstream for this environment.
	UpstreamRepo OptionalMutable[commands.TestCommands]
}

// AddCoworkerRepo adds a coworker repository.
func (self *Fixture) AddCoworkerRepo() {
	coworkerRepo := testruntime.Clone(self.OriginRepo.GetOrPanic().TestRunner, self.coworkerRepoPath())
	self.CoworkerRepo = MutableSome(&coworkerRepo)
	initializeWorkspace(&coworkerRepo)
	coworkerRepo.Verbose = self.DevRepo.GetOrPanic().Verbose
}

func (self *Fixture) AddSecondWorktree(branch gitdomain.LocalBranchName) {
	workTreePath := filepath.Join(self.Dir, "development_worktree")
	devRepo := self.DevRepo.GetOrPanic()
	devRepo.AddWorktree(workTreePath, branch)
	runner := subshell.TestRunner{
		BinDir:           devRepo.BinDir,
		HomeDir:          devRepo.HomeDir,
		ProposalOverride: None[string](),
		Verbose:          devRepo.Verbose,
		WorkingDir:       workTreePath,
	}
	gitCommands := git.Commands{
		CurrentBranchCache: &cache.WithPrevious[gitdomain.LocalBranchName]{},
		RemotesCache:       &cache.Cache[gitdomain.Remotes]{},
	}
	self.SecondWorktree = MutableSome(&commands.TestCommands{
		TestRunner: &runner,
		Git:        &gitCommands,
		Config:     devRepo.Config,
		SnapShots:  map[configdomain.ConfigScope]configdomain.SingleSnapshot{},
	})
}

// AddSubmodule adds a submodule repository.
func (self *Fixture) AddSubmoduleRepo() {
	err := os.MkdirAll(self.submoduleRepoPath(), 0o744)
	asserts.NoError(err)
	submoduleRepo := testruntime.Initialize(self.submoduleRepoPath(), self.Dir, self.binPath())
	submoduleRepo.MustRun("git", "config", "--global", "protocol.file.allow", "always")
	self.SubmoduleRepo = MutableSome(&submoduleRepo)
}

// AddUpstream adds an upstream repository.
func (self *Fixture) AddUpstream() {
	devRepo := self.DevRepo.GetOrPanic()
	upstreamRepo := testruntime.Clone(devRepo.TestRunner, filepath.Join(self.Dir, gitdomain.RemoteUpstream.String()))
	upstreamRepo.TestRunner.Verbose = devRepo.Verbose
	self.UpstreamRepo = MutableSome(&upstreamRepo)
	devRepo.AddRemote(gitdomain.RemoteUpstream, upstreamRepo.WorkingDir)
}

// Branches provides a tabular list of all branches in this Fixture.
func (self *Fixture) Branches() datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("REPOSITORY", "BRANCHES")
	mainBranch := gitdomain.NewLocalBranchName("main")
	initialBranch := gitdomain.NewLocalBranchName("initial")
	localBranches, _ := asserts.NoError2(self.DevRepo.GetOrPanic().LocalBranches())
	localBranchesJoined := localBranches.Remove(initialBranch).Hoist(mainBranch).Join(", ")
	originRepo, hasOriginRepo := self.OriginRepo.Get()
	if !hasOriginRepo {
		result.AddRow("local", localBranchesJoined)
		return result
	}
	originBranches, _ := asserts.NoError2(originRepo.LocalBranches())
	originBranchesJoined := originBranches.Remove(initialBranch).Hoist(mainBranch).Join(", ")
	originName := self.DevRepo.Value.Config.NormalConfig.DevRemote.String()
	if localBranchesJoined == originBranchesJoined {
		result.AddRow("local, "+originName, localBranchesJoined)
	} else {
		result.AddRow("local", localBranchesJoined)
		result.AddRow(originName, originBranchesJoined)
	}
	return result
}

// CommitTable provides a table for all commits in this Git environment containing only the given fields.
func (self *Fixture) CommitTable(fields []string) datatable.DataTable {
	builder := datatable.NewCommitTableBuilder()
	lineage := self.DevRepo.Value.Config.NormalConfig.Lineage
	var mainBranch gitdomain.BranchName
	mainIsLocal := self.DevRepo.Value.Git.BranchExists(self.DevRepo.Value, "main")
	if mainIsLocal {
		mainBranch = gitdomain.NewBranchName("main")
	} else {
		mainBranch = gitdomain.NewBranchName("origin/main")
	}
	localCommits := self.DevRepo.GetOrPanic().Commits(fields, mainBranch, lineage)
	builder.AddMany(localCommits, "local")
	if coworkerRepo, hasCoworkerRepo := self.CoworkerRepo.Get(); hasCoworkerRepo {
		coworkerCommits := coworkerRepo.Commits(fields, "main", lineage)
		builder.AddMany(coworkerCommits, "coworker")
	}
	if originRepo, hasOriginRepo := self.OriginRepo.Get(); hasOriginRepo {
		originCommits := originRepo.Commits(fields, "main", lineage)
		builder.AddMany(originCommits, self.DevRepo.Value.Config.NormalConfig.DevRemote.String())
	}
	if upstreamRepo, hasUpstreamRepo := self.UpstreamRepo.Get(); hasUpstreamRepo {
		upstreamCommits := upstreamRepo.Commits(fields, "main", lineage)
		builder.AddMany(upstreamCommits, "upstream")
	}
	if secondWorkTree, hasSecondWorkTree := self.SecondWorktree.Get(); hasSecondWorkTree {
		secondWorktreeCommits := secondWorkTree.Commits(fields, "main", lineage)
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
		originRepo.CheckoutBranch("initial")
	}
}

// CreateTags creates tags from the given gherkin table.
func (self *Fixture) CreateTags(table *godog.Table) {
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

func (self *Fixture) Delete() {
	os.RemoveAll(self.Dir)
}

// TagTable provides a table for all tags in this Git environment.
func (self *Fixture) TagTable() datatable.DataTable {
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
func (self *Fixture) coworkerRepoPath() string {
	return filepath.Join(self.Dir, "coworker")
}

func developerRepoPath(rootDir string) string {
	return filepath.Join(rootDir, "developer")
}

func initializeWorkspace(repo *commands.TestCommands) {
	asserts.NoError(repo.Config.SetMainBranch("main", repo.TestRunner))
	asserts.NoError(gitconfig.SetPerennialBranches(repo.TestRunner, gitdomain.LocalBranchNames{}, configdomain.ConfigScopeLocal))
	repo.MustRun("git", "checkout", "main")
	// NOTE: the developer repos receives the initial branch from origin
	//       but we don't want it here because it isn't used in tests.
	repo.MustRun("git", "branch", "-d", "initial")
}

func originRepoPath(rootDir string) string {
	return filepath.Join(rootDir, gitdomain.RemoteOrigin.String())
}

// submoduleRepoPath provides the full path to the Git repository with the given name.
func (self *Fixture) submoduleRepoPath() string {
	return filepath.Join(self.Dir, "submodule")
}
