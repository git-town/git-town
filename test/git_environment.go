package test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog/gherkin"
)

// GitEnvironment is the complete Git environment for a test scenario.
type GitEnvironment struct {

	// Dir is the directory that this environment is in.
	Dir string

	// OriginRepo is the Git repository that simulates the remote repo (on GitHub).
	OriginRepo GitRepository

	// DeveloperRepo is the Git repository that is locally checked out at the developer machine.
	DeveloperRepo GitRepository
}

// CloneGitEnvironment provides a GitEnvironment instance in the given directory,
// containing a copy of the given GitEnvironment.
func CloneGitEnvironment(original *GitEnvironment, dir string) (*GitEnvironment, error) {
	err := CopyDirectory(original.Dir, dir)
	if err != nil {
		return nil, fmt.Errorf("cannot clone GitEnvironment %q to folder %q: %w", original.Dir, dir, err)
	}
	result := GitEnvironment{
		Dir:           dir,
		DeveloperRepo: NewGitRepository(filepath.Join(dir, "developer"), dir),
		OriginRepo:    NewGitRepository(filepath.Join(dir, "origin"), dir),
	}
	// Since we copied the files from the memoized directory,
	// we have to set the "origin" remote to the copied origin repo here.
	err = result.DeveloperRepo.SetRemote(result.OriginRepo.Dir)
	return &result, err
}

// NewStandardGitEnvironment provides a GitEnvironment in the given directory,
// fully populated as a standardized setup for scenarios.
//
// The origin repo has the master branch checked out.
// Git repos cannot receive pushes of the currently checked out branch
// because that will change files in the current workspace.
// The tests don't use the master branch.
func NewStandardGitEnvironment(dir string) (gitEnv *GitEnvironment, err error) {
	// create the folder
	err = os.MkdirAll(dir, 0744)
	if err != nil {
		return gitEnv, fmt.Errorf("cannot create folder %q for Git environment: %w", dir, err)
	}
	// create the GitEnvironment
	gitEnv = &GitEnvironment{Dir: dir}
	// create the origin repo
	gitEnv.OriginRepo, err = InitGitRepository(gitEnv.originRepoPath(), gitEnv.Dir)
	if err != nil {
		return gitEnv, fmt.Errorf("cannot initialize origin directory at %q: %w", gitEnv.originRepoPath(), err)
	}
	err = gitEnv.OriginRepo.RunMany([][]string{
		{"git", "commit", "--allow-empty", "-m", "initial commit"},
		{"git", "checkout", "-b", "main"},
		{"git", "checkout", "master"},
	})
	if err != nil {
		return gitEnv, err
	}
	// clone the "developer" repo
	gitEnv.DeveloperRepo, err = CloneGitRepository(gitEnv.originRepoPath(), gitEnv.developerRepoPath(), gitEnv.Dir)
	if err != nil {
		return gitEnv, fmt.Errorf("cannot clone developer repo %q from origin %q: %w", gitEnv.originRepoPath(), gitEnv.developerRepoPath(), err)
	}
	err = gitEnv.DeveloperRepo.RunMany([][]string{
		{"git", "config", "git-town.main-branch-name", "main"},
		{"git", "config", "git-town.perennial-branch-names", ""},
		{"git", "checkout", "main"},
		// NOTE: the developer repo receives the master branch from origin
		//       but we don't want it here because it isn't used in tests.
		{"git", "branch", "-d", "master"},
	})
	return gitEnv, err
}

// CreateCommits creates the commits described by the given Gherkin table in this Git repository.
func (env *GitEnvironment) CreateCommits(table *gherkin.DataTable) error {
	commits, err := FromGherkinTable(table)
	if err != nil {
		return fmt.Errorf("cannot parse Gherkin table: %w", err)
	}
	for _, commit := range commits {
		var err error
		for _, location := range commit.Locations {
			switch location {
			case "local":
				err = env.DeveloperRepo.CreateCommit(commit)
			case "local, remote":
				err = env.DeveloperRepo.CreateCommit(commit)
				if err != nil {
					return fmt.Errorf("cannot create local commit: %w", err)
				}
				err = env.DeveloperRepo.PushBranch(commit.Branch)
				if err != nil {
					return fmt.Errorf("cannot push branch %q after creating commit: %w", commit.Branch, err)
				}
				// The developer repo has created and pushed the commit to origin already,
				// so all we need to do here is register the commit in the list of existing commits in origin.
				env.OriginRepo.RegisterOriginalCommit(commit)
			case "remote":
				err = env.OriginRepo.CreateCommit(commit)
			default:
				return fmt.Errorf("unknown commit location %q", commit.Locations)
			}
		}
		if err != nil {
			return err
		}
	}
	// after setting up the commits, check out the "master" branch in the origin repo so that we can git-push to it.
	err = env.OriginRepo.CheckoutBranch("master")
	if err != nil {
		return fmt.Errorf("cannot change origin repo back to master: %w", err)
	}
	return nil
}

// CommitTable provides a table for all commits in this Git environment containing only the given fields.
func (env GitEnvironment) CommitTable(fields []string) (result DataTable, err error) {
	builder := NewCommitTableBuilder()
	localCommits, err := env.DeveloperRepo.Commits(fields)
	if err != nil {
		return result, fmt.Errorf("cannot determine commits in the developer repo: %w", err)
	}
	for _, localCommit := range localCommits {
		builder.Add(localCommit, "local")
	}
	remoteCommits, err := env.OriginRepo.Commits(fields)
	if err != nil {
		return result, fmt.Errorf("cannot determine commits in the origin repo: %w", err)
	}
	for _, remoteCommit := range remoteCommits {
		builder.Add(remoteCommit, "remote")
	}
	return builder.Table(fields), nil
}

// developerRepoPath provides the full path to the Git repository with the given name.
func (env GitEnvironment) developerRepoPath() string {
	return filepath.Join(env.Dir, "developer")
}

// originRepoPath provides the full path to the Git repository with the given name.
func (env GitEnvironment) originRepoPath() string {
	return filepath.Join(env.Dir, "origin")
}

// Remove deletes all files used by this GitEnvironment from disk.
func (env GitEnvironment) Remove() error {
	return os.RemoveAll(env.Dir)
}
