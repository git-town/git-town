package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/util"
	"github.com/git-town/git-town/test/helpers"
	"github.com/stretchr/testify/assert"
)

// GitRepo is a Git repository that exists inside a Git environment.
type GitRepo struct {

	// Dir contains the path of the directory that this repository is in.
	Dir string

	// Shell runs console commands in this repo.
	Shell command.Shell

	// configCache contains the Git Town configuration to use.
	// This value is lazy loaded. Please use Configuration() to access it.
	configCache *git.Configuration
}

// CloneGitRepo clones a Git repo in originDir into a new GitRepository in workingDir.
// The cloning operation is using the given homeDir as the $HOME.
func CloneGitRepo(originDir, targetDir, homeDir, binDir string) (GitRepo, error) {
	res, err := command.Run("git", "clone", originDir, targetDir)
	if err != nil {
		return GitRepo{}, fmt.Errorf("cannot clone repo %q: %w\n%s", originDir, err, res.Output())
	}
	return NewGitRepository(targetDir, homeDir, NewMockingShell(targetDir, homeDir, binDir)), nil
}

// CreateTestRepo creates a GitRepo for use in tests
func CreateTestRepo(t *testing.T) GitRepo {
	dir := createTempDir(t)
	repo, err := InitGitRepository(dir, dir, "")
	assert.Nil(t, err, "cannot initialize Git repow")
	err = repo.Shell.RunMany([][]string{
		{"git", "commit", "--allow-empty", "-m", "initial commit"},
	})
	assert.Nil(t, err, "cannot create initial commit: %s")
	return repo
}

// CreateTestGitTownRepo creates a GitRepo for use in tests, with a main branch and
// initial git town configuration
func CreateTestGitTownRepo(t *testing.T) GitRepo {
	repo := CreateTestRepo(t)
	err := repo.CreateBranch("main", "master")
	assert.Nil(t, err)
	err = repo.Shell.RunMany([][]string{
		{"git", "config", "git-town.main-branch-name", "main"},
		{"git", "config", "git-town.perennial-branch-names", ""},
	})
	assert.Nil(t, err)
	return repo
}

// InitGitRepository initializes a fully functioning Git repository in the given path,
// including necessary Git configuration.
// Creates missing folders as needed.
func InitGitRepository(workingDir, homeDir, binDir string) (GitRepo, error) {
	// create the folder
	err := os.MkdirAll(workingDir, 0744)
	if err != nil {
		return GitRepo{}, fmt.Errorf("cannot create directory %q: %w", workingDir, err)
	}
	// initialize the repo in the folder
	result := NewGitRepository(workingDir, homeDir, NewMockingShell(workingDir, homeDir, binDir))
	outcome, err := result.Shell.Run("git", "init")
	if err != nil {
		return result, fmt.Errorf(`error running "git init" in %q: %w\n%v`, workingDir, err, outcome)
	}
	err = result.Shell.RunMany([][]string{
		{"git", "config", "--global", "user.name", "user"},
		{"git", "config", "--global", "user.email", "email@example.com"},
		{"git", "config", "--global", "core.editor", "vim"},
	})
	return result, err
}

// NewGitRepository provides a new GitRepository instance working in the given directory.
// The directory must contain an existing Git repo.
func NewGitRepository(workingDir string, homeDir string, shell command.Shell) GitRepo {
	return GitRepo{Dir: workingDir, Shell: shell}
}

// AddRemote adds the given Git remote to this repository.
func (repo *GitRepo) AddRemote(name, value string) error {
	res, err := repo.Shell.Run("git", "remote", "add", name, value)
	if err != nil {
		return fmt.Errorf("cannot add remote %q --> %q: %w\n%s", name, value, err, res.Output())
	}
	return nil
}

// Branches provides the names of the local branches in this Git repository,
// sorted alphabetically, with the "main" branch first.
func (repo *GitRepo) Branches() (result []string, err error) {
	outcome, err := repo.Shell.Run("git", "branch")
	if err != nil {
		return result, fmt.Errorf("cannot run 'git branch' in repo %q: %w", repo.Dir, err)
	}
	for _, line := range strings.Split(outcome.OutputSanitized(), "\n") {
		line = strings.Replace(line, "* ", "", 1)
		result = append(result, strings.TrimSpace(line))
	}
	return helpers.MainFirst(sort.StringSlice(result)), nil
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (repo *GitRepo) CheckoutBranch(name string) error {
	outcome, err := repo.Shell.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf("cannot check out branch %q in repo %q: %w\n%v", name, repo.Dir, err, outcome)
	}
	return nil
}

// Commits provides a tabular list of the commits in this Git repository with the given fields.
func (repo *GitRepo) Commits(fields []string) (result []Commit, err error) {
	branches, err := repo.Branches()
	if err != nil {
		return result, fmt.Errorf("cannot determine the Git branches: %w", err)
	}
	for _, branch := range branches {
		commits, err := repo.commitsInBranch(branch, fields)
		if err != nil {
			return result, err
		}
		result = append(result, commits...)
	}
	return result, nil
}

// CommitsInBranch provides all commits in the given Git branch.
func (repo *GitRepo) commitsInBranch(branch string, fields []string) (result []Commit, err error) {
	outcome, err := repo.Shell.Run("git", "log", branch, "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	if err != nil {
		return result, fmt.Errorf("cannot get commits in branch %q: %w", branch, err)
	}
	for _, line := range strings.Split(outcome.OutputSanitized(), "\n") {
		parts := strings.Split(line, "|")
		commit := Commit{Branch: branch, SHA: parts[0], Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if util.DoesStringArrayContain(fields, "FILE NAME") {
			filenames, err := repo.FilesInCommit(commit.SHA)
			if err != nil {
				return result, fmt.Errorf("cannot determine file name for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileName = strings.Join(filenames, ", ")
		}
		if util.DoesStringArrayContain(fields, "FILE CONTENT") {
			filecontent, err := repo.FileContentInCommit(commit.SHA, commit.FileName)
			if err != nil {
				return result, fmt.Errorf("cannot determine file content for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result, nil
}

// CommitStagedChanges commits the currently staged changes.
func (repo *GitRepo) CommitStagedChanges(message bool) error {
	var out *command.Result
	var err error
	if message {
		out, err = repo.Shell.Run("git", "commit", "-m", "committing")
	} else {
		out, err = repo.Shell.Run("git", "commit", "--no-edit")
	}
	if err != nil {
		return fmt.Errorf("cannot commit staged changes: %w\n%s", err, out)
	}
	return nil
}

// Configuration returns a cached Configuration instance for this repo.
func (repo *GitRepo) Configuration(refresh bool) *git.Configuration {
	if repo.configCache == nil || refresh {
		repo.configCache = git.NewConfiguration(repo.Shell)
	}
	return repo.configCache
}

// ConnectTrackingBranch connects the branch with the given name to its remote tracking branch.
// The branch must exist.
func (repo *GitRepo) ConnectTrackingBranch(name string) error {
	out, err := repo.Shell.Run("git", "branch", "--set-upstream-to=origin/"+name, name)
	if err != nil {
		return fmt.Errorf("cannot connect tracking branch for %q: %w\n%s", name, err, out)
	}
	return nil
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (repo *GitRepo) CreateBranch(name, parent string) error {
	outcome, err := repo.Shell.Run("git", "branch", name, parent)
	if err != nil {
		return fmt.Errorf("cannot create branch %q: %w\n%v", name, err, outcome)
	}
	return nil
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (repo *GitRepo) CreateChildFeatureBranch(name string, parent string) error {
	outcome, err := repo.Shell.Run("git", "branch", name, parent)
	if err != nil {
		return fmt.Errorf("cannot create child branch %q: %w\n%v", name, err, outcome)
	}
	outcome, err = repo.Shell.Run("git", "config", fmt.Sprintf("git-town-branch.%s.parent", name), parent)
	if err != nil {
		return fmt.Errorf("cannot create child branch %q: %w\n%v", name, err, outcome)
	}
	return nil
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (repo *GitRepo) CreateCommit(commit Commit) error {
	err := repo.CheckoutBranch(commit.Branch)
	if err != nil {
		return fmt.Errorf("cannot checkout branch %q: %w", commit.Branch, err)
	}
	err = repo.CreateFile(commit.FileName, commit.FileContent)
	if err != nil {
		return fmt.Errorf("cannot create file %q needed for commit: %w", commit.FileName, err)
	}
	outcome, err := repo.Shell.Run("git", "add", commit.FileName)
	if err != nil {
		return fmt.Errorf("cannot add file to commit: %w\n%v", err, outcome)
	}
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	outcome, err = repo.Shell.Run("git", commands...)
	if err != nil {
		return fmt.Errorf("cannot commit: %w\n%v", err, outcome)
	}
	return nil
}

// CreateFeatureBranch creates a feature branch with the given name in this repository.
func (repo *GitRepo) CreateFeatureBranch(name string) error {
	err := repo.Shell.RunMany([][]string{
		{"git", "branch", name, "main"},
		{"git", "config", "git-town-branch." + name + ".parent", "main"},
	})
	if err != nil {
		return fmt.Errorf("cannot create feature branch %q: %w", name, err)
	}
	return nil
}

// CreateFeatureBranchNoParent creates a feature branch with no defined parent in this repository.
func (repo *GitRepo) CreateFeatureBranchNoParent(name string) error {
	res, err := repo.Shell.Run("git", "branch", name, "main")
	if err != nil {
		return fmt.Errorf("cannot create feature branch %q: %w\n%s", name, err, res.Output())
	}
	return nil
}

// CreateFile creates a file with the given name and content in this repository.
func (repo *GitRepo) CreateFile(name, content string) error {
	filePath := filepath.Join(repo.Dir, name)
	folderPath := filepath.Dir(filePath)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create folder %q: %v", folderPath, err)
	}
	err = ioutil.WriteFile(filePath, []byte(content), 0744)
	if err != nil {
		return fmt.Errorf("cannot create file %q: %w", name, err)
	}
	return nil
}

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func (repo *GitRepo) CreatePerennialBranches(names ...string) error {
	for _, name := range names {
		err := repo.CreateBranch(name, "main")
		if err != nil {
			return fmt.Errorf("cannot create perennial branch %q in repo %q: %w", name, repo.Dir, err)
		}
	}
	repo.Configuration(false).AddToPerennialBranches(names...)
	return nil
}

// CreateTag creates a tag with the given name
func (repo *GitRepo) CreateTag(name string) error {
	_, err := repo.Shell.Run("git", "tag", "-a", name, "-m", name)
	return err
}

// CreateStandaloneTag creates a tag not on a branch
func (repo *GitRepo) CreateStandaloneTag(name string) error {
	return repo.Shell.RunMany([][]string{
		{"git", "checkout", "-b", "temp"},
		{"touch", "a.txt"},
		{"git", "add", "-A"},
		{"git", "commit", "-m", "temp"},
		{"git", "tag", "-a", name, "-m", name},
		{"git", "checkout", "-"},
		{"git", "branch", "-D", "temp"},
	})
}

// CurrentBranch provides the currently checked out branch for this repo.
func (repo *GitRepo) CurrentBranch() (result string, err error) {
	outcome, err := repo.Shell.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return result, fmt.Errorf("cannot determine the current branch: %w\n%s", err, outcome.Output())
	}
	return outcome.OutputSanitized(), nil
}

// FileContent provides the current content of a file.
func (repo *GitRepo) FileContent(filename string) (result string, err error) {
	outcome, err := repo.Shell.Run("cat", filename)
	if err != nil {
		return result, err
	}
	return outcome.Output(), nil
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (repo *GitRepo) DeleteMainBranchConfiguration() error {
	res, err := repo.Shell.Run("git", "config", "--unset", "git-town.main-branch-name")
	if err != nil {
		return fmt.Errorf("cannot delete main branch configuration: %w\n%s", err, res.Output())
	}
	return nil
}

// Fetch retrieves the updates from the remote repo.
func (repo *GitRepo) Fetch() error {
	_, err := repo.Shell.Run("git", "fetch")
	if err != nil {
		return fmt.Errorf("cannot fetch: %w", err)
	}
	return nil
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (repo *GitRepo) FileContentInCommit(sha string, filename string) (result string, err error) {
	outcome, err := repo.Shell.Run("git", "show", sha+":"+filename)
	if err != nil {
		return result, fmt.Errorf("cannot determine the content for file %q in commit %q: %w", filename, sha, err)
	}
	return outcome.OutputSanitized(), nil
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (repo *GitRepo) FilesInCommit(sha string) (result []string, err error) {
	outcome, err := repo.Shell.Run("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return result, fmt.Errorf("cannot get files for commit %q: %w", sha, err)
	}
	return strings.Split(outcome.OutputSanitized(), "\n"), nil
}

// FilesInBranch provides the list of the files present in the given branch.
func (repo *GitRepo) FilesInBranch(branch string) (result []string, err error) {
	outcome, err := repo.Shell.Run("git", "ls-tree", "-r", "--name-only", branch)
	if err != nil {
		return result, fmt.Errorf("cannot determine files in branch %q in repo %q: %w", branch, repo.Dir, err)
	}
	for _, line := range strings.Split(outcome.OutputSanitized(), "\n") {
		file := strings.TrimSpace(line)
		if file != "" {
			result = append(result, file)
		}
	}
	return result, err
}

// FilesInBranches provides a data table of files and their content in all branches.
func (repo *GitRepo) FilesInBranches() (result DataTable, err error) {
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := repo.Branches()
	if err != nil {
		return result, err
	}
	for _, branch := range branches {
		files, err := repo.FilesInBranch(branch)
		if err != nil {
			return result, err
		}
		for _, file := range files {
			content, err := repo.FileContentInCommit(branch, file)
			if err != nil {
				return result, err
			}
			result.AddRow(branch, file, content)
		}
	}
	return result, err
}

// HasBranchesOutOfSync indicates whether one or more local branches are out of sync with their remote
func (repo *GitRepo) HasBranchesOutOfSync() (bool, error) {
	res, err := repo.Shell.Run("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	if err != nil {
		return false, fmt.Errorf("cannot determine if branches are out of sync in %q: %w %q", repo.Dir, err, res.Output())
	}
	return strings.Contains(res.Output(), "["), nil
}

// HasFile indicates whether this repository contains a file with the given name and content.
func (repo *GitRepo) HasFile(name, content string) (result bool, err error) {
	rawContent, err := ioutil.ReadFile(filepath.Join(repo.Dir, name))
	if err != nil {
		return result, fmt.Errorf("repo doesn't have file %q: %w", name, err)
	}
	actualContent := string(rawContent)
	if actualContent != content {
		return result, fmt.Errorf("file %q should have content %q but has %q", name, content, actualContent)
	}
	return true, nil
}

// HasGitTownConfigNow indicates whether this repository contain Git Town specific configuration.
func (repo *GitRepo) HasGitTownConfigNow() (result bool, err error) {
	outcome, err := repo.Shell.Run("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		exitError := err.(*exec.ExitError)
		if exitError.ExitCode() == 1 {
			return false, nil
		}
	}
	return outcome.OutputSanitized() != "", err
}

// HasMergeInProgress indicates whether this Git repository currently has a merge in progress.
func (repo *GitRepo) HasMergeInProgress() (result bool, err error) {
	res, err := repo.Shell.Run("git", "status")
	if err != nil {
		return result, fmt.Errorf("cannot determine merge in %q progress: %w", repo.Dir, err)
	}
	return strings.Contains(res.OutputSanitized(), "You have unmerged paths"), nil
}

// HasRebaseInProgress indicates whether this Git repository currently has a rebase in progress.
func (repo *GitRepo) HasRebaseInProgress() (result bool, err error) {
	res, err := repo.Shell.Run("git", "status")
	if err != nil {
		return result, fmt.Errorf("cannot determine rebase in %q progress: %w", repo.Dir, err)
	}
	return strings.Contains(res.OutputSanitized(), "You are currently rebasing"), nil
}

// IsOffline indicates whether Git Town is offline.
func (repo *GitRepo) IsOffline() (result bool, err error) {
	res, err := repo.Shell.Run("git", "config", "--get", "git-town.offline")
	if err != nil {
		return false, fmt.Errorf("cannot determine offline status: %w\n%s", err, res.Output())
	}
	if res.OutputSanitized() == "true" {
		return true, nil
	}
	return false, nil
}

// LastActiveDir provides the directory that was last used in this repo.
func (repo *GitRepo) LastActiveDir() (string, error) {
	res, err := repo.Shell.Run("git", "rev-parse", "--show-toplevel")
	return res.OutputSanitized(), err
}

// PushBranch pushes the branch with the given name to the remote.
func (repo *GitRepo) PushBranch(name string) error {
	outcome, err := repo.Shell.Run("git", "push", "-u", "origin", name)
	if err != nil {
		return fmt.Errorf("cannot push branch %q in repo %q to origin: %w\n%v", name, repo.Dir, err, outcome)
	}
	return nil
}

// Remotes provides the names of all Git remotes in this repository.
func (repo *GitRepo) Remotes() (names []string, err error) {
	out, err := repo.Shell.Run("git", "remote")
	if err != nil {
		return names, err
	}
	if out.OutputSanitized() == "" {
		return []string{}, nil
	}
	return out.OutputLines(), nil
}

// RemoveBranch deletes the branch with the given name from this repo.
func (repo *GitRepo) RemoveBranch(name string) error {
	res, err := repo.Shell.Run("git", "branch", "-D", name)
	if err != nil {
		return fmt.Errorf("cannot delete branch %q: %w\n%s", name, err, res.Output())
	}
	return nil
}

// RemoveRemote deletes the Git remote with the given name.
func (repo *GitRepo) RemoveRemote(name string) error {
	_, err := repo.Shell.Run("git", "remote", "rm", name)
	return err
}

// RemoveUnnecessaryFiles trims all files that aren't necessary in this repo.
func (repo *GitRepo) RemoveUnnecessaryFiles() error {
	fullPath := filepath.Join(repo.Dir, ".git", "hooks")
	err := os.RemoveAll(fullPath)
	if err != nil {
		return fmt.Errorf("cannot remove unnecessary files in %q: %w", fullPath, err)
	}
	_ = os.Remove(filepath.Join(repo.Dir, ".git", "COMMIT_EDITMSG"))
	_ = os.Remove(filepath.Join(repo.Dir, ".git", "description"))
	return nil
}

// SetOffline enables or disables offline mode for this GitRepository.
func (repo *GitRepo) SetOffline(enabled bool) error {
	outcome, err := repo.Shell.Run("git", "config", "--global", "git-town.offline", strconv.FormatBool(enabled))
	if err != nil {
		return fmt.Errorf("cannot set offline mode in repo %q: %w\n%v", repo.Dir, err, outcome)
	}
	return nil
}

// Stash adds the current files to the Git stash.
func (repo *GitRepo) Stash() error {
	err := repo.Shell.RunMany([][]string{
		{"git", "add", "."},
		{"git", "stash"},
	})
	if err != nil {
		return fmt.Errorf("cannot stash: %w", err)
	}
	return nil
}

// StashSize provides the number of stashes in this repository.
func (repo *GitRepo) StashSize() (result int, err error) {
	res, err := repo.Shell.Run("git", "stash", "list")
	if err != nil {
		return result, fmt.Errorf("command %q failed: %w", res.FullCmd(), err)
	}
	if res.OutputSanitized() == "" {
		return 0, nil
	}
	return len(res.OutputLines()), nil
}

// Tags provides a list of the tags in this repository
func (repo *GitRepo) Tags() (result []string, err error) {
	res, err := repo.Shell.Run("git", "tag")
	if err != nil {
		return result, fmt.Errorf("cannot determine tags in repo %q: %w", repo.Dir, err)
	}
	for _, line := range strings.Split(res.OutputSanitized(), "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result, err
}

// UncommittedFiles provides the names of the files not committed into Git.
func (repo *GitRepo) UncommittedFiles() (result []string, err error) {
	res, err := repo.Shell.Run("git", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return result, fmt.Errorf("cannot determine uncommitted files in %q: %w", repo.Dir, err)
	}
	lines := res.OutputLines()
	for l := range lines {
		if lines[l] == "" {
			continue
		}
		parts := strings.Split(lines[l], " ")
		result = append(result, parts[1])
	}
	return result, nil
}

// ShaForCommit provides the SHA for the commit with the given name.
func (repo *GitRepo) ShaForCommit(name string) (result string, err error) {
	var args []string
	if name == "Initial commit" {
		args = []string{"reflog", "--grep=" + name, "--format=%H", "--max-count=1"}
	} else {
		args = []string{"reflog", "--grep-reflog=commit: " + name, "--format=%H"}
	}
	res, err := repo.Shell.Run("git", args...)
	if err != nil {
		return result, fmt.Errorf("cannot determine SHA of commit %q: %w\n%s", name, err, res.Output())
	}
	if res.OutputSanitized() == "" {
		return result, fmt.Errorf("cannot find the SHA of commit %q", name)
	}
	return res.OutputSanitized(), nil
}

// StageFiles adds the file with the given name to the Git index.
func (repo *GitRepo) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	_, err := repo.Shell.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot stage files %s: %w", strings.Join(names, ", "), err)
	}
	return nil
}
