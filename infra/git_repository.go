package infra

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

	// dir contains the path of the directory that this repository is in
	dir string

	// originalCommits contains the commits in this repository before the test ran
	originalCommits []CommitTableEntry

	// Runner enables to run console commands in this repo
	ShellRunner
}

// InitGitRepository initializes a new Git repository in the given folder.
func InitGitRepository(dir string, bare bool) (GitRepository, error) {

	// create the folder
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot create directory %s", dir)
	}

	// cd into the folder
	err = os.Chdir(dir)
	if err != nil {
		return GitRepository{}, errors.Wrapf(err, "cannot cd into dir %s", dir)
	}

	// initialize the repo in the folder
	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}
	runner := ShellRunner{}
	result := runner.Run("git", args...)
	if result.Err != nil {
		return GitRepository{}, errors.Wrap(result.Err, "error running git "+strings.Join(args, " "))
	}
	return GitRepository{dir: dir}, nil
}

// CloneGitRepository clones the given parent repo into a new GitRepository.
func CloneGitRepository(parentDir, childDir string) (GitRepository, error) {
	fmt.Printf("cloning parent '%s' to '%s'", parentDir, childDir)

	// clone the repo
	runner := ShellRunner{}
	result := runner.Run("git", "clone", parentDir, childDir)
	if result.Err != nil {
		return GitRepository{}, errors.Wrapf(result.Err, "cannot clone repo %s", parentDir)
	}

	// configure the repo
	err := os.Chdir(childDir)
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
	return GitRepository{dir: childDir}, err
}

// LoadGitRepository returns a GitRepository instance that manages the given existing folder
func LoadGitRepository(dir string) GitRepository {
	return GitRepository{dir: dir}
}

type CommitTableEntry struct {
	branch      string
	location    string
	message     string
	fileName    string
	fileContent string
}

// NewCommitTableEntry creates a new CommitTableEntry with default values
func NewCommitTableEntry() CommitTableEntry {
	return CommitTableEntry{
		fileName:    "default_file_name_" + uniuri.NewLen(10),
		message:     "default commit message",
		location:    "local and remote",
		branch:      "main",
		fileContent: "default file content",
	}
}

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
	runResult := gr.Run("git", "add", commit.fileName)
	if runResult.Err != nil {
		return errors.Wrapf(runResult.Err, "cannot add file to commit: %s", runResult.Output)
	}
	runResult = gr.Run("git", "commit", "-m", commit.message)
	if runResult.Err != nil {
		return errors.Wrapf(runResult.Err, "cannot commit")
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
