package test

import (
	"fmt"
	"os"
	"path"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/pkg/errors"
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
		return nil, errors.Wrapf(err, "cannot clone GitEnvironment %q to folder %q", original.Dir, dir)
	}
	result := GitEnvironment{
		Dir:           dir,
		DeveloperRepo: NewGitRepository(path.Join(dir, "developer")),
		OriginRepo:    NewGitRepository(path.Join(dir, "origin")),
	}
	// Since we copied the files from the memoized directory,
	// we have to set the "origin" remote to the copied origin repo here.
	err = result.DeveloperRepo.SetRemote(result.OriginRepo.Dir)
	return &result, err
}

// NewGitEnvironment provides a Git environment instance located in the given directory path.
// Missing directories are created as needed.
func NewGitEnvironment(baseDir string) (*GitEnvironment, error) {
	err := os.MkdirAll(baseDir, 0744)
	return &GitEnvironment{Dir: baseDir}, err
}

// NewStandardGitEnvironment provides a GitEnvironment in the given directory,
// fully populated as a standardized setup for scenarios.
func NewStandardGitEnvironment(dir string) (gitEnv *GitEnvironment, err error) {
	gitEnv, err = NewGitEnvironment(dir)
	if err != nil {
		return gitEnv, errors.Wrapf(err, "cannot create a new standard environment")
	}

	// create the origin repo
	gitEnv.OriginRepo, err = InitGitRepository(gitEnv.originRepoPath())
	if err != nil {
		return gitEnv, errors.Wrapf(err, "cannot initialize origin directory at %q", gitEnv.originRepoPath())
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
	gitEnv.DeveloperRepo, err = CloneGitRepository(gitEnv.originRepoPath(), gitEnv.developerRepoPath())
	if err != nil {
		return gitEnv, errors.Wrapf(err, "cannot clone developer repo %q from origin %q", gitEnv.originRepoPath(), gitEnv.developerRepoPath())
	}
	err = gitEnv.DeveloperRepo.RunMany([][]string{
		{"git", "checkout", "main"},
	})
	return gitEnv, err
}

// CreateCommits creates the commits described by the given Gherkin table in this Git repository.
func (env *GitEnvironment) CreateCommits(table *gherkin.DataTable) error {
	commits, err := FromGherkinTable(table)
	if err != nil {
		return errors.Wrap(err, "cannot parse Gherkin table")
	}
	for _, commit := range commits {
		var err error
		for _, location := range commit.Locations {
			switch location {
			case "local":
				err = env.DeveloperRepo.CreateCommit(commit, false)
			case "local, remote":
				err = env.DeveloperRepo.CreateCommit(commit, true)
				if err != nil {
					return err
				}
				// The developer repo has created and pushed the commit to origin already,
				// so all we need to do here is register the commit in the list of existing commits in origin.
				env.OriginRepo.RegisterOriginalCommit(commit)
			case "remote":
				err = env.OriginRepo.CreateCommit(commit, false)
			default:
				return fmt.Errorf("unknown commit location %q", commit.Locations)
			}
		}
		if err != nil {
			return err
		}
	}
	// after setting up the commits, check out the "master" branch in the origin repo so that we can git-push to it.
	env.OriginRepo.CheckoutBranch("master")
	return nil
}

// Commits provides a table for all commits in this Git environment containing only the given fields.
func (env GitEnvironment) Commits(fields []string) (result DataTable, err error) {
	builder := NewCommitTableBuilder()

	localCommits, err := env.DeveloperRepo.Commits(fields)
	if err != nil {
		return result, errors.Wrap(err, "cannot determine commits in the developer repo")
	}
	for _, localCommit := range localCommits {
		builder.Add(localCommit, "local")
	}

	remoteCommits, err := env.OriginRepo.Commits(fields)
	if err != nil {
		return result, errors.Wrap(err, "cannot determine commits in the origin repo")
	}
	for _, remoteCommit := range remoteCommits {
		builder.Add(remoteCommit, "remote")
	}

	return builder.Table(fields), nil
}

// Remove deletes all files used by this GitEnvironment from disk.
func (env GitEnvironment) Remove() error {
	return os.RemoveAll(env.Dir)
}

// developerRepoPath provides the full path to the Git repository with the given name.
func (env GitEnvironment) developerRepoPath() string {
	return path.Join(env.Dir, "developer")
}

// originRepoPath provides the full path to the Git repository with the given name.
func (env GitEnvironment) originRepoPath() string {
	return path.Join(env.Dir, "origin")
}
