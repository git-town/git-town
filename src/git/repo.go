package git

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
	"github.com/git-town/git-town/src/util"
	"github.com/git-town/git-town/test/helpers"
)

var repoInCurrentDir *Repo

// RepoInCurrentDir provides a Repo instance in the current working directory.
func RepoInCurrentDir(dryRun bool) *Repo {
	if repoInCurrentDir == nil {
		repoInCurrentDir = &Repo{
			Dir:    ".",
			dryRun: dryRun,
			Shell:  &command.ShellInCurrentDir{},
		}
	}
	return repoInCurrentDir
}

// Repo is a Git repository that exists inside a Git environment.
type Repo struct {

	// Dir contains the path of the directory that this repository is in.
	Dir string

	// Shell runs console commands in this repo.
	command.Shell

	// configCache contains the Git Town configuration to use.
	// This value is lazy loaded. Please use Configuration() to access it.
	configCache *Configuration

	// currentBranch contains the current Git branch we are on in this repo
	currentBranch string

	// dryRun indicates whether dryRun is enabled
	dryRun bool
}

// AddRemote adds the given Git remote to this repository.
func (repo *Repo) AddRemote(name, value string) error {
	res, err := repo.Run("git", "remote", "add", name, value)
	if err != nil {
		return fmt.Errorf("cannot add remote %q --> %q: %w\n%s", name, value, err, res.Output())
	}
	return nil
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (repo *Repo) CheckoutBranch(name string) error {
	outcome, err := repo.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf("cannot check out branch %q in repo %q: %w\n%v", name, repo.Dir, err, outcome)
	}
	repo.currentBranch = name
	return nil
}

// CommitStagedChanges commits the currently staged changes.
func (repo *Repo) CommitStagedChanges(message bool) error {
	var out *command.Result
	var err error
	if message {
		out, err = repo.Run("git", "commit", "-m", "committing")
	} else {
		out, err = repo.Run("git", "commit", "--no-edit")
	}
	if err != nil {
		return fmt.Errorf("cannot commit staged changes: %w\n%s", err, out)
	}
	return nil
}

// Config returns a cached Config instance for this repo.
func (repo *Repo) Config(refresh bool) *Configuration {
	if repo.configCache == nil || refresh {
		repo.configCache = NewConfiguration(repo.Shell)
	}
	return repo.configCache
}

// ConnectTrackingBranch connects the branch with the given name to its remote tracking branch.
// The branch must exist.
func (repo *Repo) ConnectTrackingBranch(name string) error {
	out, err := repo.Run("git", "branch", "--set-upstream-to=origin/"+name, name)
	if err != nil {
		return fmt.Errorf("cannot connect tracking branch for %q: %w\n%s", name, err, out)
	}
	return nil
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (repo *Repo) CreateBranch(name, parent string) error {
	outcome, err := repo.Run("git", "branch", name, parent)
	if err != nil {
		return fmt.Errorf("cannot create branch %q: %w\n%v", name, err, outcome)
	}
	return nil
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (repo *Repo) CreateChildFeatureBranch(name string, parent string) error {
	outcome, err := repo.Run("git", "branch", name, parent)
	if err != nil {
		return fmt.Errorf("cannot create child branch %q: %w\n%v", name, err, outcome)
	}
	outcome, err = repo.Run("git", "config", fmt.Sprintf("git-town-branch.%s.parent", name), parent)
	if err != nil {
		return fmt.Errorf("cannot create child branch %q: %w\n%v", name, err, outcome)
	}
	return nil
}

// CreateFeatureBranch creates a feature branch with the given name in this repository.
func (repo *Repo) CreateFeatureBranch(name string) error {
	err := repo.RunMany([][]string{
		{"git", "branch", name, "main"},
		{"git", "config", "git-town-branch." + name + ".parent", "main"},
	})
	if err != nil {
		return fmt.Errorf("cannot create feature branch %q: %w", name, err)
	}
	return nil
}

// CreateFeatureBranchNoParent creates a feature branch with no defined parent in this repository.
func (repo *Repo) CreateFeatureBranchNoParent(name string) error {
	res, err := repo.Run("git", "branch", name, "main")
	if err != nil {
		return fmt.Errorf("cannot create feature branch %q: %w\n%s", name, err, res.Output())
	}
	return nil
}

// CreateFile creates a file with the given name and content in this repository.
func (repo *Repo) CreateFile(name, content string) error {
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
func (repo *Repo) CreatePerennialBranches(names ...string) error {
	for _, name := range names {
		err := repo.CreateBranch(name, "main")
		if err != nil {
			return fmt.Errorf("cannot create perennial branch %q in repo %q: %w", name, repo.Dir, err)
		}
	}
	repo.Config(false).AddToPerennialBranches(names...)
	return nil
}

// CreateTag creates a tag with the given name
func (repo *Repo) CreateTag(name string) error {
	_, err := repo.Run("git", "tag", "-a", name, "-m", name)
	return err
}

// CreateStandaloneTag creates a tag not on a branch
func (repo *Repo) CreateStandaloneTag(name string) error {
	return repo.RunMany([][]string{
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
func (repo *Repo) CurrentBranch() (result string, err error) {
	outcome, err := repo.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return result, fmt.Errorf("cannot determine the current branch: %w\n%s", err, outcome.Output())
	}
	return outcome.OutputSanitized(), nil
}

// FileContent provides the current content of a file.
func (repo *Repo) FileContent(filename string) (result string, err error) {
	outcome, err := repo.Run("cat", filename)
	if err != nil {
		return result, err
	}
	return outcome.Output(), nil
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (repo *Repo) DeleteMainBranchConfiguration() error {
	res, err := repo.Run("git", "config", "--unset", "git-town.main-branch-name")
	if err != nil {
		return fmt.Errorf("cannot delete main branch configuration: %w\n%s", err, res.Output())
	}
	return nil
}

// Fetch retrieves the updates from the remote repo.
func (repo *Repo) Fetch() error {
	_, err := repo.Run("git", "fetch")
	if err != nil {
		return fmt.Errorf("cannot fetch: %w", err)
	}
	return nil
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (repo *Repo) FileContentInCommit(sha string, filename string) (result string, err error) {
	outcome, err := repo.Run("git", "show", sha+":"+filename)
	if err != nil {
		return result, fmt.Errorf("cannot determine the content for file %q in commit %q: %w", filename, sha, err)
	}
	return outcome.OutputSanitized(), nil
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (repo *Repo) FilesInCommit(sha string) (result []string, err error) {
	outcome, err := repo.Run("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return result, fmt.Errorf("cannot get files for commit %q: %w", sha, err)
	}
	return strings.Split(outcome.OutputSanitized(), "\n"), nil
}

// FilesInBranch provides the list of the files present in the given branch.
func (repo *Repo) FilesInBranch(branch string) (result []string, err error) {
	outcome, err := repo.Run("git", "ls-tree", "-r", "--name-only", branch)
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

// HasBranchesOutOfSync indicates whether one or more local branches are out of sync with their remote
func (repo *Repo) HasBranchesOutOfSync() (bool, error) {
	res, err := repo.Run("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	if err != nil {
		return false, fmt.Errorf("cannot determine if branches are out of sync in %q: %w %q", repo.Dir, err, res.Output())
	}
	return strings.Contains(res.Output(), "["), nil
}

// HasFile indicates whether this repository contains a file with the given name and content.
func (repo *Repo) HasFile(name, content string) (result bool, err error) {
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
func (repo *Repo) HasGitTownConfigNow() (result bool, err error) {
	outcome, err := repo.Run("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		exitError := err.(*exec.ExitError)
		if exitError.ExitCode() == 1 {
			return false, nil
		}
	}
	return outcome.OutputSanitized() != "", err
}

// HasLocalBranch indicates whether this repo has a local branch with the given name.
func (repo *Repo) HasLocalBranch(name string) (bool, error) {
	branches, err := repo.LocalBranches()
	if err != nil {
		return false, fmt.Errorf("cannot determine whether the local branch %q exists: %w", name, err)
	}
	return util.DoesStringArrayContain(branches, name), nil
}

// HasLocalOrRemoteBranch indicates whether this repo has a local or remote branch with the given name.
func (repo *Repo) HasLocalOrRemoteBranch(name string) (bool, error) {
	branches, err := repo.LocalAndRemoteBranches()
	if err != nil {
		return false, fmt.Errorf("cannot determine whether the local or remote branch %q exists: %w", name, err)
	}
	return util.DoesStringArrayContain(branches, name), nil
}

// HasMergeInProgress indicates whether this Git repository currently has a merge in progress.
func (repo *Repo) HasMergeInProgress() (result bool, err error) {
	res, err := repo.Run("git", "status")
	if err != nil {
		return result, fmt.Errorf("cannot determine merge in %q progress: %w", repo.Dir, err)
	}
	return strings.Contains(res.OutputSanitized(), "You have unmerged paths"), nil
}

// HasRebaseInProgress indicates whether this Git repository currently has a rebase in progress.
func (repo *Repo) HasRebaseInProgress() (result bool, err error) {
	res, err := repo.Run("git", "status")
	if err != nil {
		return result, fmt.Errorf("cannot determine rebase in %q progress: %w", repo.Dir, err)
	}
	return strings.Contains(res.OutputSanitized(), "You are currently rebasing"), nil
}

// HasRemote indicates whether this repo has a remote with the given name.
func (repo *Repo) HasRemote(name string) (result bool, err error) {
	remotes, err := repo.Remotes()
	if err != nil {
		return false, fmt.Errorf("cannot determine if remote %q exists: %w", name, err)
	}
	return util.DoesStringArrayContain(remotes, name), nil
}

// IsOffline indicates whether Git Town is offline.
func (repo *Repo) IsOffline() (result bool, err error) {
	res, err := repo.Run("git", "config", "--get", "git-town.offline")
	if err != nil {
		return false, fmt.Errorf("cannot determine offline status: %w\n%s", err, res.Output())
	}
	if res.OutputSanitized() == "true" {
		return true, nil
	}
	return false, nil
}

// LastActiveDir provides the directory that was last used in this repo.
func (repo *Repo) LastActiveDir() (string, error) {
	res, err := repo.Run("git", "rev-parse", "--show-toplevel")
	return res.OutputSanitized(), err
}

// LocalBranches provides the names of all local branches in this repo.
func (repo *Repo) LocalBranches() (result []string, err error) {
	outcome, err := repo.Run("git", "branch")
	if err != nil {
		return result, fmt.Errorf("cannot determine the local branches")
	}
	lines := outcome.OutputLines()
	for l := range lines {
		result = append(result, strings.TrimSpace(strings.Trim(lines[l], "* ")))
	}
	return helpers.MainFirst(sort.StringSlice(result)), nil
}

// LocalAndRemoteBranches provides the names of all local branches in this repo.
func (repo *Repo) LocalAndRemoteBranches() ([]string, error) {
	outcome, err := repo.Run("git", "branch", "-a")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine the local branches")
	}
	lines := outcome.OutputLines()
	branchNames := make(map[string]struct{})
	for l := range lines {
		if !strings.Contains(lines[l], " -> ") {
			branchNames[strings.TrimSpace(strings.Replace(strings.Replace(lines[l], "* ", "", 1), "remotes/origin/", "", 1))] = struct{}{}
		}
	}
	result := make([]string, len(branchNames))
	i := 0
	for branchName := range branchNames {
		result[i] = branchName
		i++
	}
	sort.Strings(result)
	return helpers.MainFirst(result), nil
}

// PushBranch pushes the branch with the given name to the remote.
func (repo *Repo) PushBranch(name string) error {
	outcome, err := repo.Run("git", "push", "-u", "origin", name)
	if err != nil {
		return fmt.Errorf("cannot push branch %q in repo %q to origin: %w\n%v", name, repo.Dir, err, outcome)
	}
	return nil
}

// Remotes provides the names of all Git remotes in this repository.
func (repo *Repo) Remotes() (names []string, err error) {
	out, err := repo.Run("git", "remote")
	if err != nil {
		return names, err
	}
	if out.OutputSanitized() == "" {
		return []string{}, nil
	}
	return out.OutputLines(), nil
}

// RemoveBranch deletes the branch with the given name from this repo.
func (repo *Repo) RemoveBranch(name string) error {
	res, err := repo.Run("git", "branch", "-D", name)
	if err != nil {
		return fmt.Errorf("cannot delete branch %q: %w\n%s", name, err, res.Output())
	}
	return nil
}

// RemoveRemote deletes the Git remote with the given name.
func (repo *Repo) RemoveRemote(name string) error {
	_, err := repo.Run("git", "remote", "rm", name)
	return err
}

// RemoveUnnecessaryFiles trims all files that aren't necessary in this repo.
func (repo *Repo) RemoveUnnecessaryFiles() error {
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
func (repo *Repo) SetOffline(enabled bool) error {
	outcome, err := repo.Run("git", "config", "--global", "git-town.offline", strconv.FormatBool(enabled))
	if err != nil {
		return fmt.Errorf("cannot set offline mode in repo %q: %w\n%v", repo.Dir, err, outcome)
	}
	return nil
}

// Stash adds the current files to the Git stash.
func (repo *Repo) Stash() error {
	err := repo.RunMany([][]string{
		{"git", "add", "."},
		{"git", "stash"},
	})
	if err != nil {
		return fmt.Errorf("cannot stash: %w", err)
	}
	return nil
}

// StashSize provides the number of stashes in this repository.
func (repo *Repo) StashSize() (result int, err error) {
	res, err := repo.Run("git", "stash", "list")
	if err != nil {
		return result, fmt.Errorf("command %q failed: %w", res.FullCmd(), err)
	}
	if res.OutputSanitized() == "" {
		return 0, nil
	}
	return len(res.OutputLines()), nil
}

// Tags provides a list of the tags in this repository
func (repo *Repo) Tags() (result []string, err error) {
	res, err := repo.Run("git", "tag")
	if err != nil {
		return result, fmt.Errorf("cannot determine tags in repo %q: %w", repo.Dir, err)
	}
	for _, line := range strings.Split(res.OutputSanitized(), "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result, err
}

// UncommittedFiles provides the names of the files not committed into Git.
func (repo *Repo) UncommittedFiles() (result []string, err error) {
	res, err := repo.Run("git", "status", "--porcelain", "--untracked-files=all")
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
func (repo *Repo) ShaForCommit(name string) (result string, err error) {
	var args []string
	if name == "Initial commit" {
		args = []string{"reflog", "--grep=" + name, "--format=%H", "--max-count=1"}
	} else {
		args = []string{"reflog", "--grep-reflog=commit: " + name, "--format=%H"}
	}
	res, err := repo.Run("git", args...)
	if err != nil {
		return result, fmt.Errorf("cannot determine SHA of commit %q: %w\n%s", name, err, res.Output())
	}
	if res.OutputSanitized() == "" {
		return result, fmt.Errorf("cannot find the SHA of commit %q", name)
	}
	return res.OutputSanitized(), nil
}

// StageFiles adds the file with the given name to the Git index.
func (repo *Repo) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	_, err := repo.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot stage files %s: %w", strings.Join(names, ", "), err)
	}
	return nil
}
