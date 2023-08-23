package fixture

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/slice"
	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/git-town/git-town/v9/test/datatable"
	"github.com/git-town/git-town/v9/test/git"
	"github.com/git-town/git-town/v9/test/helpers"
	"github.com/git-town/git-town/v9/test/testruntime"
)

// Fixture is a complete Git environment for a Cucumber scenario.
type Fixture struct {
	// Dir defines the local folder in which this Fixture is stored.
	// This folder also acts as the HOME directory for tests using this Fixture.
	// It contains the global Git configuration to use in this test.
	Dir string

	// CoworkerRepo is the optional Git repository that is locally checked out at the coworker machine.
	CoworkerRepo *testruntime.TestRuntime `exhaustruct:"optional"`

	// DevRepo is the Git repository that is locally checked out at the developer machine.
	DevRepo testruntime.TestRuntime `exhaustruct:"optional"`

	// OriginRepo is the Git repository that simulates the origin repo (on GitHub).
	// If this value is nil, the current test setup has no origin.
	OriginRepo *testruntime.TestRuntime `exhaustruct:"optional"`

	// SubmoduleRepo is the Git repository that simulates an external repo used as a submodule.
	// If this value is nil, the current test setup uses no submodules.
	SubmoduleRepo *testruntime.TestRuntime `exhaustruct:"optional"`

	// UpstreamRepo is the optional Git repository that contains the upstream for this environment.
	UpstreamRepo *testruntime.TestRuntime `exhaustruct:"optional"`
}

// CloneFixture provides a Fixture instance in the given directory,
// containing a copy of the given Fixture.
func CloneFixture(original Fixture, dir string) Fixture {
	helpers.CopyDirectory(original.Dir, dir)
	binDir := filepath.Join(dir, "bin")
	originDir := filepath.Join(dir, "origin")
	originRepo := testruntime.New(originDir, dir, "")
	developerDir := filepath.Join(dir, "developer")
	devRepo := testruntime.New(developerDir, dir, binDir)
	result := Fixture{
		Dir:        dir,
		DevRepo:    devRepo,
		OriginRepo: &originRepo,
	}
	// Since we copied the files from the memoized directory,
	// we have to set the "origin" remote to the copied origin repo here.
	result.DevRepo.MustRun("git", "remote", "remove", domain.OriginRemote.String())
	result.DevRepo.AddRemote(domain.OriginRemote, result.originRepoPath())
	result.DevRepo.Fetch()
	// and connect the main branches again
	result.DevRepo.ConnectTrackingBranch(domain.NewLocalBranchName("main"))
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
		{"git", "commit", "--allow-empty", "-m", "Initial commit"},
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

// AddSubmodule adds a submodule repository.
func (env *Fixture) AddSubmoduleRepo() {
	err := os.MkdirAll(env.submoduleRepoPath(), 0o744)
	if err != nil {
		log.Fatalf("cannot create directory %q: %v", env.submoduleRepoPath(), err)
	}
	submoduleRepo := testruntime.Initialize(env.submoduleRepoPath(), env.Dir, env.binPath())
	submoduleRepo.MustRunMany([][]string{
		{"git", "config", "--global", "protocol.file.allow", "always"},
		{"git", "commit", "--allow-empty", "-m", "Initial commit"},
	})
	env.SubmoduleRepo = &submoduleRepo
}

// AddUpstream adds an upstream repository.
func (env *Fixture) AddUpstream() {
	repo := testruntime.Clone(env.DevRepo.TestRunner, filepath.Join(env.Dir, "upstream"))
	env.UpstreamRepo = &repo
	env.DevRepo.AddRemote(domain.UpstreamRemote, env.UpstreamRepo.WorkingDir)
}

// AddCoworkerRepo adds a coworker repository.
func (env *Fixture) AddCoworkerRepo() {
	coworkerRepo := testruntime.Clone(env.OriginRepo.TestRunner, env.coworkerRepoPath())
	env.CoworkerRepo = &coworkerRepo
	env.initializeWorkspace(env.CoworkerRepo)
	env.CoworkerRepo.Debug = env.DevRepo.Debug
}

// binPath provides the full path of the folder containing the test tools for this Fixture.
func (env *Fixture) binPath() string {
	return filepath.Join(env.Dir, "bin")
}

// Branches provides a tabular list of all branches in this Fixture.
func (env *Fixture) Branches() datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("REPOSITORY", "BRANCHES")
	mainBranch := env.DevRepo.Config.MainBranch()
	localBranches, err := env.DevRepo.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	initialBranch := domain.NewLocalBranchName("initial")
	localBranches = slice.Remove(localBranches, initialBranch)
	localBranchesJoined := localBranches.Join(", ")
	if env.OriginRepo == nil {
		result.AddRow("local", localBranchesJoined)
		return result
	}
	originBranches, err := env.OriginRepo.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	originBranches = slice.Remove(originBranches, initialBranch)
	originBranchesJoined := originBranches.Join(", ")
	if localBranchesJoined == originBranchesJoined {
		result.AddRow("local, origin", localBranchesJoined)
	} else {
		result.AddRow("local", localBranchesJoined)
		result.AddRow("origin", originBranchesJoined)
	}
	return result
}

// CreateCommits creates the commits described by the given Gherkin table in this Git repository.
func (env *Fixture) CreateCommits(commits []git.Commit) {
	for _, commit := range commits {
		for _, location := range commit.Locations {
			switch location {
			case "coworker":
				env.CoworkerRepo.CreateCommit(commit)
			case "local":
				env.DevRepo.CreateCommit(commit)
			case "local, origin":
				env.DevRepo.CreateCommit(commit)
				env.DevRepo.PushBranch()
			case "origin":
				env.OriginRepo.CreateCommit(commit)
			case "upstream":
				env.UpstreamRepo.CreateCommit(commit)
			default:
				log.Fatalf("unknown commit location %q", commit.Locations)
			}
		}
	}
	// after setting up the commits, check out the "initial" branch in the origin repo so that we can git-push to it.
	if env.OriginRepo != nil {
		env.OriginRepo.CheckoutBranch(domain.NewLocalBranchName("initial"))
	}
}

// CreateOriginBranch creates a branch with the given name only in the origin directory.
func (env Fixture) CreateOriginBranch(name, parent string) {
	env.OriginRepo.CreateBranch(domain.NewLocalBranchName(name), domain.NewLocalBranchName(parent))
}

// CreateTags creates tags from the given gherkin table.
func (env Fixture) CreateTags(table *messages.PickleStepArgument_PickleTable) {
	columnNames := helpers.TableFields(table)
	if columnNames[0] != "NAME" && columnNames[1] != "LOCATION" {
		log.Fatalf("tag table must have columns NAME and LOCATION")
	}
	for _, row := range table.Rows[1:] {
		name := row.Cells[0].Value
		location := row.Cells[1].Value
		switch location {
		case "local":
			env.DevRepo.CreateTag(name)
		case "origin":
			env.OriginRepo.CreateTag(name)
		default:
			log.Fatalf("tag table LOCATION must be 'local' or 'origin'")
		}
	}
}

// CommitTable provides a table for all commits in this Git environment containing only the given fields.
func (env Fixture) CommitTable(fields []string) datatable.DataTable {
	builder := datatable.NewCommitTableBuilder()
	localCommits := env.DevRepo.Commits(fields, domain.NewLocalBranchName("main"))
	builder.AddMany(localCommits, "local")
	if env.CoworkerRepo != nil {
		coworkerCommits := env.CoworkerRepo.Commits(fields, domain.NewLocalBranchName("main"))
		builder.AddMany(coworkerCommits, "coworker")
	}
	if env.OriginRepo != nil {
		originCommits := env.OriginRepo.Commits(fields, domain.NewLocalBranchName("main"))
		builder.AddMany(originCommits, domain.OriginRemote.String())
	}
	if env.UpstreamRepo != nil {
		upstreamCommits := env.UpstreamRepo.Commits(fields, domain.NewLocalBranchName("main"))
		builder.AddMany(upstreamCommits, "upstream")
	}
	return builder.Table(fields)
}

// TagTable provides a table for all tags in this Git environment.
func (env Fixture) TagTable() datatable.DataTable {
	builder := datatable.NewTagTableBuilder()
	localTags := env.DevRepo.Tags()
	builder.AddMany(localTags, "local")
	if env.OriginRepo != nil {
		originTags := env.OriginRepo.Tags()
		builder.AddMany(originTags, domain.OriginRemote.String())
	}
	return builder.Table()
}

func (env Fixture) initializeWorkspace(repo *testruntime.TestRuntime) {
	asserts.NoError(repo.Config.SetMainBranch(domain.NewLocalBranchName("main")))
	asserts.NoError(repo.Config.SetPerennialBranches(domain.LocalBranchNames{}))
	repo.MustRunMany([][]string{
		{"git", "checkout", "main"},
		// NOTE: the developer repos receives the initial branch from origin
		//       but we don't want it here because it isn't used in tests.
		{"git", "branch", "-d", "initial"},
	})
}

// coworkerRepoPath provides the full path to the Git repository with the given name.
func (env Fixture) coworkerRepoPath() string {
	return filepath.Join(env.Dir, "coworker")
}

// developerRepoPath provides the full path to the Git repository with the given name.
func (env Fixture) developerRepoPath() string {
	return filepath.Join(env.Dir, "developer")
}

// originRepoPath provides the full path to the Git repository with the given name.
func (env Fixture) originRepoPath() string {
	return filepath.Join(env.Dir, domain.OriginRemote.String())
}

// submoduleRepoPath provides the full path to the Git repository with the given name.
func (env Fixture) submoduleRepoPath() string {
	return filepath.Join(env.Dir, "submodule")
}

// Remove deletes all files used by this Fixture from disk.
func (env Fixture) Remove() {
	asserts.NoError(os.RemoveAll(env.Dir))
}
