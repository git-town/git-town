package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/dchest/uniuri"
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

// InitGitRepository initializes a new Git repository in the given path.
// The given path must not exist.
func InitGitRepository(dir string, bare bool) (GitRepository, error) {

	// create the folder
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot create directory %q", dir)
	}

	// cd into the folder
	err = os.Chdir(dir)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot cd into dir %q", dir)
	}

	// initialize the repo in the folder
	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}
	result := GitRepository{dir: dir}
	_, err = result.Run("git", args...)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "error running git %s", strings.Join(args, " "))
	}
	return result, nil
}

// CloneGitRepository clones the given parent repo into a new GitRepository.
func CloneGitRepository(parentDir, childDir string) (GitRepository, error) {
	// clone the repo
	runner := ShellRunner{}
	_, err := runner.Run("git", "clone", parentDir, childDir)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot clone repo %s", parentDir)
	}

	// configure the repo
	result := GitRepository{dir: childDir}
	err = os.Chdir(childDir)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot cd into %s", childDir)
	}
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

// LoadGitRepository returns a GitRepository instance that manages the given existing folder
func LoadGitRepository(dir string) GitRepository {
	return GitRepository{dir: dir}
}

// CommitTableEntry contains the elements of a Gherkin table defining commit data.
type CommitTableEntry struct {
	branch      string
	location    string
	message     string
	fileName    string
	fileContent string
}

// NewCommitTableEntry provides a new CommitTableEntry with default values
func NewCommitTableEntry() CommitTableEntry {
	return CommitTableEntry{
		fileName:    "default_file_name_" + uniuri.NewLen(10),
		message:     "default commit message",
		location:    "local and remote",
		branch:      "main",
		fileContent: "default file content",
	}
}

// CreateCommits creates the commits described by the given Gherkin table in this Git repository.
func (gr *GitRepository) CreateCommits(table *gherkin.DataTable) error {
	err := os.Chdir(gr.dir)
	if err != nil {
		return errors.Wrapf(err, "cannot cd into root dir of repo: %s", gr.dir)
	}
	gr.originalCommits = gr.parseCommitsTable(table)
	for _, commit := range gr.originalCommits {
		err := gr.createCommit(commit)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gr *GitRepository) createCommit(commit CommitTableEntry) error {
	err := ioutil.WriteFile(commit.fileName, []byte(commit.fileContent), 0744)
	if err != nil {
		return errors.Wrapf(err, "cannot create file '%s' to commit", commit.fileName)
	}
	dir, err := os.Getwd()
	fmt.Println("CWD", dir, err)
	output, err := gr.Run("git", "add", commit.fileName)
	if err != nil {
		return errors.Wrapf(err, "cannot add file to commit: %s", output)
	}
	_, err = gr.Run("git", "commit", "-m", commit.message)
	if err != nil {
		return errors.Wrapf(err, "cannot commit")
	}
	return nil
}

func (gr *GitRepository) parseCommitsTable(table *gherkin.DataTable) []CommitTableEntry {
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
