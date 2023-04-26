package test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/git-town/git-town/v8/src/config"
	prodgit "github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/stringslice"
	"github.com/git-town/git-town/v8/test/datatable"
	"github.com/git-town/git-town/v8/test/git"
	"github.com/git-town/git-town/v8/test/runner"
)

// testCommands defines Git commands used only in test code.
type testCommands struct {
	runner.Mocking
	config prodgit.RepoConfig
	*prodgit.BackendCommands
}

// AddRemote adds a Git remote with the given name and URL to this repository.
func (r *testCommands) AddRemote(name, url string) error {
	_, err := r.Run("git", "remote", "add", name, url)
	if err != nil {
		return fmt.Errorf("cannot add remote %q --> %q: %w", name, url, err)
	}
	r.config.RemotesCache.Invalidate()
	return nil
}

// AddSubmodule adds a Git submodule with the given URL to this repository.
func (r *testCommands) AddSubmodule(url string) error {
	_, err := r.Run("git", "submodule", "add", url)
	if err != nil {
		return err
	}
	_, err = r.Run("git", "commit", "-m", "added submodule")
	return err
}

// BranchHierarchyTable provides the currently configured branch hierarchy information as a DataTable.
func (r *testCommands) BranchHierarchyTable() datatable.DataTable {
	result := datatable.DataTable{}
	r.config.Reload()
	parentBranchMap := r.config.ParentBranchMap()
	result.AddRow("BRANCH", "PARENT")
	childBranches := make([]string, 0, len(parentBranchMap))
	for child := range parentBranchMap {
		childBranches = append(childBranches, child)
	}
	sort.Strings(childBranches)
	for _, child := range childBranches {
		result.AddRow(child, parentBranchMap[child])
	}
	return result
}

// Clone creates a clone of the repository managed by this test.Runner into the given directory.
// The cloned repo uses the same homeDir and binDir as its origin.
func (r *testCommands) Clone(targetDir string) (Runner, error) {
	_, err := r.Run("git", "clone", r.WorkingDir, targetDir)
	if err != nil {
		return Runner{}, fmt.Errorf("cannot clone repo %q: %w", r.WorkingDir, err)
	}
	return newRunner(targetDir, r.HomeDir, r.BinDir), nil
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (r *testCommands) CheckoutBranch(name string) error {
	_, err := r.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf("cannot check out branch %q: %w", name, err)
	}
	if name != "-" {
		r.config.CurrentBranchCache.Set(name)
	} else {
		r.config.CurrentBranchCache.Invalidate()
	}
	return nil
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (r *testCommands) CreateBranch(name, parent string) error {
	_, err := r.Run("git", "branch", name, parent)
	return err
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (r *testCommands) CreateChildFeatureBranch(name string, parent string) error {
	err := r.CreateBranch(name, parent)
	if err != nil {
		return err
	}
	return r.config.SetParent(name, parent)
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (r *testCommands) CreateCommit(commit git.Commit) error {
	err := r.CheckoutBranch(commit.Branch)
	if err != nil {
		return fmt.Errorf("cannot checkout branch %q: %w", commit.Branch, err)
	}
	err = r.CreateFile(commit.FileName, commit.FileContent)
	if err != nil {
		return fmt.Errorf("cannot create file %q needed for commit: %w", commit.FileName, err)
	}
	_, err = r.Run("git", "add", commit.FileName)
	if err != nil {
		return fmt.Errorf("cannot add file to commit: %w", err)
	}
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	_, err = r.Run("git", commands...)
	if err != nil {
		return fmt.Errorf("cannot commit: %w", err)
	}
	return nil
}

// CreateFile creates a file with the given name and content in this repository.
func (r *testCommands) CreateFile(name, content string) error {
	filePath := filepath.Join(r.WorkingDir, name)
	folderPath := filepath.Dir(filePath)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create folder %q: %w", folderPath, err)
	}
	err = os.WriteFile(filePath, []byte(content), 0o500)
	if err != nil {
		return fmt.Errorf("cannot create file %q: %w", name, err)
	}
	return nil
}

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func (r *testCommands) CreatePerennialBranches(names ...string) error {
	for _, name := range names {
		err := r.CreateBranch(name, "main")
		if err != nil {
			return fmt.Errorf("cannot create perennial branch %q in repo %q: %w", name, r.WorkingDir, err)
		}
	}
	return r.config.AddToPerennialBranches(names...)
}

// CreateStandaloneTag creates a tag not on a branch.
func (r *testCommands) CreateStandaloneTag(name string) error {
	return r.RunMany([][]string{
		{"git", "checkout", "-b", "temp"},
		{"touch", "a.txt"},
		{"git", "add", "-A"},
		{"git", "commit", "-m", "temp"},
		{"git", "tag", "-a", name, "-m", name},
		{"git", "checkout", "-"},
		{"git", "branch", "-D", "temp"},
	})
}

// CreateTag creates a tag with the given name.
func (r *testCommands) CreateTag(name string) error {
	_, err := r.Run("git", "tag", "-a", name, "-m", name)
	return err
}

// Commits provides a list of the commits in this Git repository with the given fields.
func (r *testCommands) Commits(fields []string, mainBranch string) ([]git.Commit, error) {
	branches, err := r.LocalBranchesMainFirst(mainBranch)
	if err != nil {
		return []git.Commit{}, fmt.Errorf("cannot determine the Git branches: %w", err)
	}
	result := []git.Commit{}
	for _, branch := range branches {
		commits, err := r.CommitsInBranch(branch, fields)
		if err != nil {
			return []git.Commit{}, err
		}
		result = append(result, commits...)
	}
	return result, nil
}

// CommitsInBranch provides all commits in the given Git branch.
func (r *testCommands) CommitsInBranch(branch string, fields []string) ([]git.Commit, error) {
	output, err := r.Run("git", "log", branch, "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	if err != nil {
		return []git.Commit{}, fmt.Errorf("cannot get commits in branch %q: %w", branch, err)
	}
	result := []git.Commit{}
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, "|")
		commit := git.Commit{Branch: branch, SHA: parts[0], Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if stringslice.Contains(fields, "FILE NAME") {
			filenames, err := r.FilesInCommit(commit.SHA)
			if err != nil {
				return []git.Commit{}, fmt.Errorf("cannot determine file name for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileName = strings.Join(filenames, ", ")
		}
		if stringslice.Contains(fields, "FILE CONTENT") {
			filecontent, err := r.FileContentInCommit(commit.SHA, commit.FileName)
			if err != nil {
				return []git.Commit{}, fmt.Errorf("cannot determine file content for commit %q in branch %q: %w", commit.SHA, branch, err)
			}
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result, nil
}

// CommitStagedChanges commits the currently staged changes.
func (r *testCommands) CommitStagedChanges(message string) error {
	_, err := r.Run("git", "commit", "-m", message)
	return err
}

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func (r *testCommands) ConnectTrackingBranch(name string) error {
	_, err := r.Run("git", "branch", "--set-upstream-to=origin/"+name, name)
	return err
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (r *testCommands) DeleteMainBranchConfiguration() error {
	_, err := r.Run("git", "config", "--unset", config.MainBranchKey)
	if err != nil {
		return fmt.Errorf("cannot delete main branch configuration: %w", err)
	}
	return nil
}

// Fetch retrieves the updates from the origin repo.
func (r *testCommands) Fetch() error {
	_, err := r.Run("git", "fetch")
	return err
}

// FileContent provides the current content of a file.
func (r *testCommands) FileContent(filename string) (string, error) {
	content, err := os.ReadFile(filepath.Join(r.WorkingDir, filename))
	return string(content), err
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (r *testCommands) FileContentInCommit(sha string, filename string) (string, error) {
	output, err := r.Run("git", "show", sha+":"+filename)
	if err != nil {
		return "", fmt.Errorf("cannot determine the content for file %q in commit %q: %w", filename, sha, err)
	}
	result := output
	if strings.HasPrefix(result, "tree ") {
		// merge commits get an empty file content instead of "tree <SHA>"
		result = ""
	}
	return result, nil
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (r *testCommands) FilesInCommit(sha string) ([]string, error) {
	output, err := r.Run("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return []string{}, fmt.Errorf("cannot get files for commit %q: %w", sha, err)
	}
	return strings.Split(output, "\n"), nil
}

// FilesInBranch provides the list of the files present in the given branch.
func (r *testCommands) FilesInBranch(branch string) ([]string, error) {
	output, err := r.Run("git", "ls-tree", "-r", "--name-only", branch)
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine files in branch %q in repo %q: %w", branch, r.WorkingDir, err)
	}
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		file := strings.TrimSpace(line)
		if file != "" {
			result = append(result, file)
		}
	}
	return result, err
}

// FilesInBranches provides a data table of files and their content in all branches.
func (r *testCommands) FilesInBranches(mainBranch string) (datatable.DataTable, error) {
	result := datatable.DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := r.LocalBranchesMainFirst(mainBranch)
	if err != nil {
		return datatable.DataTable{}, err
	}
	lastBranch := ""
	for _, branch := range branches {
		files, err := r.FilesInBranch(branch)
		if err != nil {
			return datatable.DataTable{}, err
		}
		for _, file := range files {
			content, err := r.FileContentInCommit(branch, file)
			if err != nil {
				return datatable.DataTable{}, err
			}
			if branch == lastBranch {
				result.AddRow("", file, content)
			} else {
				result.AddRow(branch, file, content)
			}
			lastBranch = branch
		}
	}
	return result, err
}

// HasBranchesOutOfSync indicates whether one or more local branches are out of sync with their tracking branch.
func (r *testCommands) HasBranchesOutOfSync() (bool, error) {
	output, err := r.Run("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	if err != nil {
		return false, fmt.Errorf("cannot determine if branches are out of sync in %q: %w %q", r.WorkingDir, err, output)
	}
	return strings.Contains(output, "["), nil
}

// HasFile indicates whether this repository contains a file with the given name and content.
func (r *testCommands) HasFile(name, content string) (bool, error) {
	rawContent, err := os.ReadFile(filepath.Join(r.WorkingDir, name))
	if err != nil {
		return false, fmt.Errorf("repo doesn't have file %q: %w", name, err)
	}
	actualContent := string(rawContent)
	if actualContent != content {
		return false, fmt.Errorf("file %q should have content %q but has %q", name, content, actualContent)
	}
	return true, nil
}

// HasGitTownConfigNow indicates whether this repository contain Git Town specific configuration.
func (r *testCommands) HasGitTownConfigNow() bool {
	output, err := r.Run("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		return false
	}
	return output != ""
}

func (r *testCommands) PushBranch() error {
	_, err := r.Run("git", "push")
	return err
}

func (r *testCommands) PushBranchToRemote(branch, remote string) error {
	_, err := r.Run("git", "push", "-u", remote, branch)
	return err
}

// RemoveBranch deletes the branch with the given name from this repo.
func (r *testCommands) RemoveBranch(name string) error {
	_, err := r.Run("git", "branch", "-D", name)
	return err
}

// RemoveRemote deletes the Git remote with the given name.
func (r *testCommands) RemoveRemote(name string) error {
	r.config.RemotesCache.Invalidate()
	_, err := r.Run("git", "remote", "rm", name)
	return err
}

// RemoveUnnecessaryFiles trims all files that aren't necessary in this repo.
func (r *testCommands) RemoveUnnecessaryFiles() error {
	fullPath := filepath.Join(r.WorkingDir, ".git", "hooks")
	err := os.RemoveAll(fullPath)
	if err != nil {
		return fmt.Errorf("cannot remove unnecessary files in %q: %w", fullPath, err)
	}
	_ = os.Remove(filepath.Join(r.WorkingDir, ".git", "COMMIT_EDITMSG"))
	_ = os.Remove(filepath.Join(r.WorkingDir, ".git", "description"))
	return nil
}

// ShaForCommit provides the SHA for the commit with the given name.
func (r *testCommands) ShaForCommit(name string) (string, error) {
	output, err := r.Run("git", "log", "--reflog", "--format=%H", "--grep=^"+name+"$")
	if err != nil {
		return "", fmt.Errorf("cannot determine the SHA of commit %q: %w", name, err)
	}
	result := output
	if result == "" {
		return "", fmt.Errorf("cannot find the SHA of commit %q", name)
	}
	result = strings.Split(result, "\n")[0]
	return result, nil
}

// StageFiles adds the file with the given name to the Git index.
func (r *testCommands) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	_, err := r.Run("git", args...)
	return err
}

// StashSize provides the number of stashes in this repository.
func (r *testCommands) StashSize() (int, error) {
	output, err := r.Run("git", "stash", "list")
	if err != nil {
		return 0, fmt.Errorf("cannot determine Git stash: %w", err)
	}
	if output == "" {
		return 0, nil
	}
	return len(stringslice.Lines(output)), nil
}

// Tags provides a list of the tags in this repository.
func (r *testCommands) Tags() ([]string, error) {
	output, err := r.Run("git", "tag")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine tags in repo %q: %w", r.WorkingDir, err)
	}
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result, err
}

// UncommittedFiles provides the names of the files not committed into Git.
func (r *testCommands) UncommittedFiles() ([]string, error) {
	output, err := r.Run("git", "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine uncommitted files in %q: %w", r.WorkingDir, err)
	}
	result := []string{}
	for _, line := range stringslice.Lines(output) {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		result = append(result, parts[1])
	}
	return result, nil
}
