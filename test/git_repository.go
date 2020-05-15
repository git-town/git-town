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

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/util"
	"github.com/git-town/git-town/test/helpers"
)

// GitRepository is a Git repository that exists inside a Git environment.
type GitRepository struct {

	// Dir contains the path of the directory that this repository is in.
	Dir string

	// originalCommits contains the commits in this repository before the test ran.
	originalCommits []Commit

	// Shell runs console commands in this repo.
	Shell command.Shell

	// configCache contains the Git Town configuration to use.
	// This value is lazy loaded. Please use Configuration() to access it.
	configCache *git.Configuration
}

// CloneGitRepository clones a Git repo in originDir into a new GitRepository in workingDir.
// The cloning operation is using the given homeDir as the $HOME.
func CloneGitRepository(originDir, targetDir, homeDir, binDir string) (GitRepository, error) {
	res, err := command.Run("git", "clone", originDir, targetDir)
	if err != nil {
		return GitRepository{}, fmt.Errorf("cannot clone repo %q: %w\n%s", originDir, err, res.Output())
	}
	return NewGitRepository(targetDir, homeDir, NewMockingShell(targetDir, homeDir, binDir)), nil
}

// InitGitRepository initializes a fully functioning Git repository in the given path,
// including necessary Git configuration.
// Creates missing folders as needed.
func InitGitRepository(workingDir, homeDir, binDir string) (GitRepository, error) {
	// create the folder
	err := os.MkdirAll(workingDir, 0744)
	if err != nil {
		return GitRepository{}, fmt.Errorf("cannot create directory %q: %w", workingDir, err)
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
func NewGitRepository(workingDir string, homeDir string, shell command.Shell) GitRepository {
	return GitRepository{Dir: workingDir, Shell: shell}
}

// HasBranchesOutOfSync returns if true if and only if one or more local branches are out of sync with their remote
func (repo *GitRepository) HasBranchesOutOfSync() (bool, error) {
	res, err := repo.Shell.Run("bash", "-c", "git branch -vv | grep -o \"\\[.*\\]\" | tr -d \"[]\" | awk \"{ print \\$2 }\" | grep . | wc -l")
	if err != nil {
		return false, fmt.Errorf("cannot determine if any branches are out of sync in %q: %w %q", repo.Dir, err, res.Output())
	}
	count, err := strconv.ParseInt(res.OutputSanitized(), 0, 64)
	if err != nil {
		return false, fmt.Errorf("cannot parse output as int in %q: %w %q", repo.Dir, err, res.Output())
	}
	return count != 0, nil
}

// AddRemote adds the given Git remote to this repository.
func (repo *GitRepository) AddRemote(name, value string) error {
	res, err := repo.Shell.Run("git", "remote", "add", name, value)
	if err != nil {
		return fmt.Errorf("cannot add remote %q --> %q: %w\n%s", name, value, err, res.Output())
	}
	return nil
}

// Branches provides the names of the local branches in this Git repository,
// sorted alphabetically, with the "main" branch first.
func (repo *GitRepository) Branches() (result []string, err error) {
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
func (repo *GitRepository) CheckoutBranch(name string) error {
	outcome, err := repo.Shell.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf("cannot check out branch %q in repo %q: %w\n%v", name, repo.Dir, err, outcome)
	}
	return nil
}

// Commits provides a tabular list of the commits in this Git repository with the given fields.
func (repo *GitRepository) Commits(fields []string) (result []Commit, err error) {
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
func (repo *GitRepository) commitsInBranch(branch string, fields []string) (result []Commit, err error) {
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
func (repo *GitRepository) CommitStagedChanges(message bool) error {
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
func (repo *GitRepository) Configuration(refresh bool) *git.Configuration {
	if repo.configCache == nil || refresh {
		repo.configCache = git.NewConfiguration(repo.Shell)
	}
	return repo.configCache
}

// ConnectTrackingBranch connects the branch with the given name to its remote tracking branch.
// The branch must exist.
func (repo *GitRepository) ConnectTrackingBranch(name string) error {
	out, err := repo.Shell.Run("git", "branch", "--set-upstream-to=origin/"+name, name)
	if err != nil {
		return fmt.Errorf("cannot connect tracking branch for %q: %w\n%s", name, err, out)
	}
	return nil
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (repo *GitRepository) CreateBranch(name, parent string) error {
	outcome, err := repo.Shell.Run("git", "branch", name, parent)
	if err != nil {
		return fmt.Errorf("cannot create branch %q: %w\n%v", name, err, outcome)
	}
	return nil
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (repo *GitRepository) CreateChildFeatureBranch(name string, parent string) error {
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
func (repo *GitRepository) CreateCommit(commit Commit) error {
	repo.originalCommits = append(repo.originalCommits, commit)
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
func (repo *GitRepository) CreateFeatureBranch(name string) error {
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
func (repo *GitRepository) CreateFeatureBranchNoParent(name string) error {
	res, err := repo.Shell.Run("git", "branch", name, "main")
	if err != nil {
		return fmt.Errorf("cannot create feature branch %q: %w\n%s", name, err, res.Output())
	}
	return nil
}

// CreateFile creates a file with the given name and content in this repository.
func (repo *GitRepository) CreateFile(name, content string) error {
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
func (repo *GitRepository) CreatePerennialBranches(names ...string) error {
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
func (repo *GitRepository) CreateTag(name string) error {
	_, err := repo.Shell.Run("git", "tag", "-a", name, "-m", name)
	return err
}

// CreateStandaloneTag creates a tag not on a branch
func (repo *GitRepository) CreateStandaloneTag(name string) error {
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
func (repo *GitRepository) CurrentBranch() (result string, err error) {
	outcome, err := repo.Shell.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return result, fmt.Errorf("cannot determine the current branch: %w\n%s", err, outcome.Output())
	}
	return outcome.OutputSanitized(), nil
}

// FileContent provides the current content of a file.
func (repo *GitRepository) FileContent(filename string) (result string, err error) {
	outcome, err := repo.Shell.Run("cat", filename)
	if err != nil {
		return result, err
	}
	return outcome.Output(), nil
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (repo *GitRepository) DeleteMainBranchConfiguration() error {
	res, err := repo.Shell.Run("git", "config", "--unset", "git-town.main-branch-name")
	if err != nil {
		return fmt.Errorf("cannot delete main branch configuration: %w\n%s", err, res.Output())
	}
	return nil
}

// Fetch retrieves the updates from the remote repo.
func (repo *GitRepository) Fetch() error {
	_, err := repo.Shell.Run("git", "fetch")
	if err != nil {
		return fmt.Errorf("cannot fetch: %w", err)
	}
	return nil
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (repo *GitRepository) FileContentInCommit(sha string, filename string) (result string, err error) {
	outcome, err := repo.Shell.Run("git", "show", sha+":"+filename)
	if err != nil {
		return result, fmt.Errorf("cannot determine the content for file %q in commit %q: %w", filename, sha, err)
	}
	return outcome.OutputSanitized(), nil
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (repo *GitRepository) FilesInCommit(sha string) (result []string, err error) {
	outcome, err := repo.Shell.Run("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return result, fmt.Errorf("cannot get files for commit %q: %w", sha, err)
	}
	return strings.Split(outcome.OutputSanitized(), "\n"), nil
}

// FilesInBranch provides the list of the files present in the given branch.
func (repo *GitRepository) FilesInBranch(branch string) (result []string, err error) {
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
func (repo *GitRepository) FilesInBranches() (result DataTable, err error) {
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

// HasFile indicates whether this repository contains a file with the given name and content.
func (repo *GitRepository) HasFile(name, content string) (result bool, err error) {
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
func (repo *GitRepository) HasGitTownConfigNow() (result bool, err error) {
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
func (repo *GitRepository) HasMergeInProgress() (result bool, err error) {
	res, err := repo.Shell.Run("git", "status")
	if err != nil {
		return result, fmt.Errorf("cannot determine merge in %q progress: %w", repo.Dir, err)
	}
	return strings.Contains(res.OutputSanitized(), "You have unmerged paths"), nil
}

// HasRebaseInProgress indicates whether this Git repository currently has a rebase in progress.
func (repo *GitRepository) HasRebaseInProgress() (result bool, err error) {
	res, err := repo.Shell.Run("git", "status")
	if err != nil {
		return result, fmt.Errorf("cannot determine rebase in %q progress: %w", repo.Dir, err)
	}
	return strings.Contains(res.OutputSanitized(), "You are currently rebasing"), nil
}

// IsOffline indicates whether Git Town is offline.
func (repo *GitRepository) IsOffline() (result bool, err error) {
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
func (repo *GitRepository) LastActiveDir() (string, error) {
	res, err := repo.Shell.Run("git", "rev-parse", "--show-toplevel")
	return res.OutputSanitized(), err
}

// PushBranch pushes the branch with the given name to the remote.
func (repo *GitRepository) PushBranch(name string) error {
	outcome, err := repo.Shell.Run("git", "push", "-u", "origin", name)
	if err != nil {
		return fmt.Errorf("cannot push branch %q in repo %q to origin: %w\n%v", name, repo.Dir, err, outcome)
	}
	return nil
}

// RegisterOriginalCommit tracks the given commit as existing in this repo before the system under test executed.
func (repo *GitRepository) RegisterOriginalCommit(commit Commit) {
	repo.originalCommits = append(repo.originalCommits, commit)
}

// Remotes provides the names of all Git remotes in this repository.
func (repo *GitRepository) Remotes() (names []string, err error) {
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
func (repo *GitRepository) RemoveBranch(name string) error {
	res, err := repo.Shell.Run("git", "branch", "-D", name)
	if err != nil {
		return fmt.Errorf("cannot delete branch %q: %w\n%s", name, err, res.Output())
	}
	return nil
}

// RemoveRemote deletes the Git remote with the given name.
func (repo *GitRepository) RemoveRemote(name string) error {
	_, err := repo.Shell.Run("git", "remote", "rm", name)
	return err
}

// SetOffline enables or disables offline mode for this GitRepository.
func (repo *GitRepository) SetOffline(enabled bool) error {
	outcome, err := repo.Shell.Run("git", "config", "--global", "git-town.offline", strconv.FormatBool(enabled))
	if err != nil {
		return fmt.Errorf("cannot set offline mode in repo %q: %w\n%v", repo.Dir, err, outcome)
	}
	return nil
}

// Stash adds the current files to the Git stash.
func (repo *GitRepository) Stash() error {
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
func (repo *GitRepository) StashSize() (result int, err error) {
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
func (repo *GitRepository) Tags() (result []string, err error) {
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
func (repo *GitRepository) UncommittedFiles() (result []string, err error) {
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
func (repo *GitRepository) ShaForCommit(name string) (result string, err error) {
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
func (repo *GitRepository) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	_, err := repo.Shell.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot stage files %s: %w", strings.Join(names, ", "), err)
	}
	return nil
}
