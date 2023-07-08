package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	prodgit "github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/stringslice"
	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/git-town/git-town/v9/test/datatable"
	"github.com/git-town/git-town/v9/test/git"
	"github.com/git-town/git-town/v9/test/subshell"
)

// TestCommands defines Git commands used only in test code.
type TestCommands struct {
	subshell.TestRunner
	*prodgit.BackendCommands // TODO: remove this dependency on BackendCommands
}

// AddRemote adds a Git remote with the given name and URL to this repository.
func (r *TestCommands) AddRemote(name, url string) {
	r.MustRun("git", "remote", "add", name, url)
	r.Config.RemotesCache.Invalidate()
}

// AddSubmodule adds a Git submodule with the given URL to this repository.
func (r *TestCommands) AddSubmodule(url string) {
	r.MustRun("git", "submodule", "add", url)
	r.MustRun("git", "commit", "-m", "added submodule")
}

// BranchHierarchyTable provides the currently configured branch hierarchy information as a DataTable.
func (r *TestCommands) BranchHierarchyTable() datatable.DataTable {
	result := datatable.DataTable{}
	r.Config.Reload()
	parentBranchMap := r.Config.Lineage().Entries
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

// .CheckoutBranch checks out the Git branch with the given name in this repo.
func (r *TestCommands) CheckoutBranch(name string) {
	r.MustRun("git", "checkout", name)
	if name != "-" {
		r.Config.CurrentBranchCache.Set(name)
	} else {
		r.Config.CurrentBranchCache.Invalidate()
	}
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (r *TestCommands) CreateBranch(name, parent string) {
	r.MustRun("git", "branch", name, parent)
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (r *TestCommands) CreateChildFeatureBranch(name string, parent string) {
	r.CreateBranch(name, parent)
	asserts.NoError(r.Config.SetParent(name, parent))
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (r *TestCommands) CreateCommit(commit git.Commit) {
	r.CheckoutBranch(commit.Branch)
	r.CreateFile(commit.FileName, commit.FileContent)
	r.MustRun("git", "add", commit.FileName)
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	r.MustRun("git", commands...)
}

// CreateFile creates a file with the given name and content in this repository.
func (r *TestCommands) CreateFile(name, content string) {
	filePath := filepath.Join(r.WorkingDir, name)
	folderPath := filepath.Dir(filePath)
	asserts.NoError(os.MkdirAll(folderPath, os.ModePerm))
	//nolint:gosec // need permission 700 here in order for tests to work
	asserts.NoError(os.WriteFile(filePath, []byte(content), 0x700))
}

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func (r *TestCommands) CreatePerennialBranches(names ...string) {
	for _, name := range names {
		r.CreateBranch(name, "main")
	}
	asserts.NoError(r.Config.AddToPerennialBranches(names...))
}

// CreateStandaloneTag creates a tag not on a branch.
func (r *TestCommands) CreateStandaloneTag(name string) {
	r.MustRun("git", "checkout", "-b", "temp")
	r.MustRun("touch", "a.txt")
	r.MustRun("git", "add", "-A")
	r.MustRun("git", "commit", "-m", "temp")
	r.MustRun("git", "tag", "-a", name, "-m", name)
	r.MustRun("git", "checkout", "-")
	r.MustRun("git", "branch", "-D", "temp")
}

// CreateTag creates a tag with the given name.
func (r *TestCommands) CreateTag(name string) {
	r.MustRun("git", "tag", "-a", name, "-m", name)
}

// Commits provides a list of the commits in this Git repository with the given fields.
func (r *TestCommands) Commits(fields []string, mainBranch string) []git.Commit {
	branches, err := r.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	result := []git.Commit{}
	for _, branch := range branches {
		commits := r.CommitsInBranch(branch, fields)
		result = append(result, commits...)
	}
	return result
}

// CommitsInBranch provides all commits in the given Git branch.
func (r *TestCommands) CommitsInBranch(branch string, fields []string) []git.Commit {
	output := r.MustQuery("git", "log", branch, "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	result := []git.Commit{}
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, "|")
		commit := git.Commit{Branch: branch, SHA: parts[0], Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if stringslice.Contains(fields, "FILE NAME") {
			filenames := r.FilesInCommit(commit.SHA)
			commit.FileName = strings.Join(filenames, ", ")
		}
		if stringslice.Contains(fields, "FILE CONTENT") {
			filecontent := r.FileContentInCommit(commit.SHA, commit.FileName)
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result
}

// CommitStagedChanges commits the currently staged changes.
func (r *TestCommands) CommitStagedChanges(message string) {
	r.MustRun("git", "commit", "-m", message)
}

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func (r *TestCommands) ConnectTrackingBranch(name string) {
	r.MustRun("git", "branch", "--set-upstream-to=origin/"+name, name)
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (r *TestCommands) DeleteMainBranchConfiguration() {
	r.MustRun("git", "config", "--unset", config.MainBranchKey)
}

// Fetch retrieves the updates from the origin repo.
func (r *TestCommands) Fetch() {
	r.MustRun("git", "fetch")
}

// FileContent provides the current content of a file.
func (r *TestCommands) FileContent(filename string) string {
	content, err := os.ReadFile(filepath.Join(r.WorkingDir, filename))
	asserts.NoError(err)
	return string(content)
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (r *TestCommands) FileContentInCommit(sha string, filename string) string {
	output := r.MustQuery("git", "show", sha+":"+filename)
	if strings.HasPrefix(output, "tree ") {
		// merge commits get an empty file content instead of "tree <SHA>"
		return ""
	}
	return output
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (r *TestCommands) FilesInCommit(sha string) []string {
	output := r.MustQuery("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	return strings.Split(output, "\n")
}

// FilesInBranch provides the list of the files present in the given branch.
func (r *TestCommands) FilesInBranch(branch string) []string {
	output := r.MustQuery("git", "ls-tree", "-r", "--name-only", branch)
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		file := strings.TrimSpace(line)
		if file != "" {
			result = append(result, file)
		}
	}
	return result
}

// FilesInBranches provides a data table of files and their content in all branches.
func (r *TestCommands) FilesInBranches(mainBranch string) datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := r.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	lastBranch := ""
	for _, branch := range branches {
		files := r.FilesInBranch(branch)
		for _, file := range files {
			content := r.FileContentInCommit(branch, file)
			if branch == lastBranch {
				result.AddRow("", file, content)
			} else {
				result.AddRow(branch, file, content)
			}
			lastBranch = branch
		}
	}
	return result
}

// HasBranchesOutOfSync indicates whether one or more local branches are out of sync with their tracking branch.
func (r *TestCommands) HasBranchesOutOfSync() bool {
	output := r.MustQuery("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	return strings.Contains(output, "[")
}

// HasFile indicates whether this repository contains a file with the given name and content.
// An empty error message means a file with the given name and content exists.
func (r *TestCommands) HasFile(name, content string) string {
	rawContent, err := os.ReadFile(filepath.Join(r.WorkingDir, name))
	if err != nil {
		return fmt.Sprintf("repo doesn't have file %q", name)
	}
	actualContent := string(rawContent)
	if actualContent != content {
		return fmt.Sprintf("file %q should have content %q but has %q", name, content, actualContent)
	}
	return ""
}

// HasGitTownConfigNow indicates whether this repository contain Git Town specific configuration.
func (r *TestCommands) HasGitTownConfigNow() bool {
	output, err := r.Query("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		return false
	}
	return output != ""
}

func (r *TestCommands) PushBranch() {
	r.MustRun("git", "push")
}

func (r *TestCommands) PushBranchToRemote(branch, remote string) {
	r.MustRun("git", "push", "-u", remote, branch)
}

// RemoveBranch deletes the branch with the given name from this repo.
func (r *TestCommands) RemoveBranch(name string) {
	r.MustRun("git", "branch", "-D", name)
}

// RemoveRemote deletes the Git remote with the given name.
func (r *TestCommands) RemoveRemote(name string) {
	r.Config.RemotesCache.Invalidate()
	r.MustRun("git", "remote", "rm", name)
}

// RemoveUnnecessaryFiles trims all files that aren't necessary in this repo.
func (r *TestCommands) RemoveUnnecessaryFiles() {
	fullPath := filepath.Join(r.WorkingDir, ".git", "hooks")
	asserts.NoError(os.RemoveAll(fullPath))
	_ = os.Remove(filepath.Join(r.WorkingDir, ".git", "COMMIT_EDITMSG"))
	_ = os.Remove(filepath.Join(r.WorkingDir, ".git", "description"))
}

// ShaForCommit provides the SHA for the commit with the given name.
func (r *TestCommands) ShaForCommit(name string) string {
	output := r.MustQuery("git", "log", "--reflog", "--format=%H", "--grep=^"+name+"$")
	if output == "" {
		log.Fatalf("cannot find the SHA of commit %q", name)
	}
	return strings.Split(output, "\n")[0]
}

// StageFiles adds the file with the given name to the Git index.
func (r *TestCommands) StageFiles(names ...string) {
	args := append([]string{"add"}, names...)
	r.MustRun("git", args...)
}

// StashSize provides the number of stashes in this repository.
func (r *TestCommands) StashSize() int {
	output := r.MustQuery("git", "stash", "list")
	if output == "" {
		return 0
	}
	return len(stringslice.Lines(output))
}

// Tags provides a list of the tags in this repository.
func (r *TestCommands) Tags() []string {
	output := r.MustQuery("git", "tag")
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result
}

// UncommittedFiles provides the names of the files not committed into Git.
func (r *TestCommands) UncommittedFiles() []string {
	output := r.MustQuery("git", "status", "--porcelain", "--untracked-files=all")
	result := []string{}
	for _, line := range stringslice.Lines(output) {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		result = append(result, parts[1])
	}
	return result
}
