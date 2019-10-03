package test

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/pkg/errors"
)

// GitRepository is a Git repository that exists inside a Git environment.
type GitRepository struct {

	// dir contains the path of the directory that this repository is in.
	dir string

	// originalCommits contains the commits in this repository before the test ran.
	originalCommits []CommitTableEntry

	// ShellRunner enables to run console commands in this repo.
	ShellRunner
}

// NewGitRepository provides a new GitRepository instance working in the given directory.
// The directory must contain an existing Git repo.
func NewGitRepository(dir string) GitRepository {
	result := GitRepository{dir: dir}
	result.ShellRunner = NewShellRunner(dir)
	return result
}

// InitGitRepository initializes a new Git repository in the given path.
// Creates missing folders as needed.
func InitGitRepository(dir string, bare bool) (GitRepository, error) {

	// create the folder
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot create directory %q", dir)
	}

	// initialize the repo in the folder
	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}
	result := NewGitRepository(dir)
	_, err = result.Run("git", args...)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "error running git %s", strings.Join(args, " "))
	}
	return result, nil
}

// CloneGitRepository clones the given parent repo into a new GitRepository.
func CloneGitRepository(parentDir, childDir string) (GitRepository, error) {
	runner := NewShellRunner(".")
	_, err := runner.Run("git", "clone", parentDir, childDir)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot clone repo %q", parentDir)
	}
	result := NewGitRepository(childDir)
	userName := strings.Replace(path.Base(childDir), "_secondary", "", 1)
	err = runner.RunMany([][]string{
		[]string{"git", "config", "user.name", userName},
		[]string{"git", "config", "user.email", userName + "@example.com"},
		[]string{"git", "config", "push.default", "simple"},
		[]string{"git", "config", "core.editor", "vim"},
		[]string{"git", "config", "git-town.main-branch-name", "main"},
		[]string{"git", "config", "git-town.perennial-branch-names", ""},
	})
	return result, err
}

// CreateCommits creates the commits described by the given Gherkin table in this Git repository.
func (repo *GitRepository) CreateCommits(table *gherkin.DataTable) error {
	repo.originalCommits = repo.parseCommitsTable(table)
	for _, commit := range repo.originalCommits {
		err := repo.createCommit(commit)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *GitRepository) createCommit(commit CommitTableEntry) error {
	err := repo.createFile(path.Join(repo.dir, commit.fileName), commit.fileContent)
	if err != nil {
		return err
	}
	output, err := repo.Run("git", "add", commit.fileName)
	if err != nil {
		return errors.Wrapf(err, "cannot add file to commit: %s", output)
	}
	_, err = repo.Run("git", "commit", "-m", commit.message)
	if err != nil {
		return errors.Wrapf(err, "cannot commit")
	}
	return nil
}

// createFile creates a file with the given name and content in this repository.
func (repo *GitRepository) createFile(name, content string) error {
	err := ioutil.WriteFile(path.Join(repo.dir, name), []byte(content), 0744)
	if err != nil {
		return errors.Wrapf(err, "cannot create file %q", name)
	}
	return nil
}

func (repo *GitRepository) parseCommitsTable(table *gherkin.DataTable) []CommitTableEntry {
	result := []CommitTableEntry{}
	columnNames := []string{}
	for _, cell := range table.Rows[0].Cells {
		columnNames = append(columnNames, cell.Value)
	}
	for _, row := range table.Rows[1:] {
		commit := NewCommitTableEntry()
		for i, cell := range row.Cells {
			switch columnNames[i] {
			case "BRANCH":
				commit.branch = cell.Value
			case "LOCATION":
				commit.location = cell.Value
			case "MESSAGE":
				commit.message = cell.Value
			default:
				log.Fatalf("GitRepository.parseCommitsTable: unknown column name: %s", columnNames[i])
			}
		}
		result = append(result, commit)
	}
	return result
}
