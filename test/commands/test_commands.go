package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	prodgit "github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/slice"
	"github.com/git-town/git-town/v9/src/stringslice"
	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/git-town/git-town/v9/test/datatable"
	"github.com/git-town/git-town/v9/test/git"
	"github.com/git-town/git-town/v9/test/subshell"
)

// TestCommands defines Git commands used only in test code.
type TestCommands struct {
	*subshell.TestRunner
	*prodgit.BackendCommands // TODO: remove this dependency on BackendCommands
}

// AddRemote adds a Git remote with the given name and URL to this repository.
func (tc *TestCommands) AddRemote(name domain.Remote, url string) {
	tc.MustRun("git", "remote", "add", name.String(), url)
	tc.RemotesCache.Invalidate()
}

// AddSubmodule adds a Git submodule with the given URL to this repository.
func (tc *TestCommands) AddSubmodule(url string) {
	tc.MustRun("git", "submodule", "add", url)
	tc.MustRun("git", "commit", "-m", "added submodule")
}

// BranchHierarchyTable provides the currently configured branch hierarchy information as a DataTable.
func (tc *TestCommands) BranchHierarchyTable() datatable.DataTable {
	result := datatable.DataTable{}
	tc.Config.Reload()
	lineage := tc.Config.Lineage()
	result.AddRow("BRANCH", "PARENT")
	for _, branchName := range lineage.BranchNames() {
		result.AddRow(branchName.String(), lineage[branchName].String())
	}
	return result
}

// .CheckoutBranch checks out the Git branch with the given name in this repo.
func (tc *TestCommands) CheckoutBranch(branch domain.LocalBranchName) {
	asserts.NoError(tc.BackendCommands.CheckoutBranch(branch))
}

func (tc *TestCommands) CommitSHAs() map[string]domain.SHA {
	result := map[string]domain.SHA{}
	output := tc.MustQuery("git", "log", "--all", "--pretty=format:%h %s")
	for _, line := range strings.Split(output, "\n") {
		parts := strings.SplitN(line, " ", 2)
		sha := parts[0]
		commitMessage := parts[1]
		result[commitMessage] = domain.NewSHA(sha)
	}
	return result
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (tc *TestCommands) CreateBranch(name, parent domain.LocalBranchName) {
	tc.MustRun("git", "branch", name.String(), parent.String())
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (tc *TestCommands) CreateChildFeatureBranch(branch domain.LocalBranchName, parent domain.LocalBranchName) {
	tc.CreateBranch(branch, parent)
	asserts.NoError(tc.Config.SetParent(branch, parent))
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (tc *TestCommands) CreateCommit(commit git.Commit) {
	tc.CheckoutBranch(commit.Branch)
	tc.CreateFile(commit.FileName, commit.FileContent)
	tc.MustRun("git", "add", commit.FileName)
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	tc.MustRun("git", commands...)
}

// CreateFile creates a file with the given name and content in this repository.
func (tc *TestCommands) CreateFile(name, content string) {
	filePath := filepath.Join(tc.WorkingDir, name)
	folderPath := filepath.Dir(filePath)
	asserts.NoError(os.MkdirAll(folderPath, os.ModePerm))
	//nolint:gosec // need permission 700 here in order for tests to work
	asserts.NoError(os.WriteFile(filePath, []byte(content), 0x700))
}

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func (tc *TestCommands) CreatePerennialBranches(names ...domain.LocalBranchName) {
	for _, name := range names {
		tc.CreateBranch(name, domain.NewLocalBranchName("main"))
	}
	asserts.NoError(tc.Config.AddToPerennialBranches(names...))
}

// CreateStandaloneTag creates a tag not on a branch.
func (tc *TestCommands) CreateStandaloneTag(name string) {
	tc.MustRun("git", "checkout", "-b", "temp")
	tc.MustRun("touch", "a.txt")
	tc.MustRun("git", "add", "-A")
	tc.MustRun("git", "commit", "-m", "temp")
	tc.MustRun("git", "tag", "-a", name, "-m", name)
	tc.MustRun("git", "checkout", "-")
	tc.MustRun("git", "branch", "-D", "temp")
}

// CreateTag creates a tag with the given name.
func (tc *TestCommands) CreateTag(name string) {
	tc.MustRun("git", "tag", "-a", name, "-m", name)
}

// Commits provides a list of the commits in this Git repository with the given fields.
func (tc *TestCommands) Commits(fields []string, mainBranch domain.LocalBranchName) []git.Commit {
	branches, err := tc.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	result := []git.Commit{}
	for _, branch := range branches {
		commits := tc.CommitsInBranch(branch, fields)
		result = append(result, commits...)
	}
	return result
}

// CommitsInBranch provides all commits in the given Git branch.
func (tc *TestCommands) CommitsInBranch(branch domain.LocalBranchName, fields []string) []git.Commit {
	output := tc.MustQuery("git", "log", branch.String(), "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	result := []git.Commit{}
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, "|")
		commit := git.Commit{Branch: branch, SHA: domain.NewSHA(parts[0]), Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if slice.Contains(fields, "FILE NAME") {
			filenames := tc.FilesInCommit(commit.SHA)
			commit.FileName = strings.Join(filenames, ", ")
		}
		if slice.Contains(fields, "FILE CONTENT") {
			filecontent := tc.FileContentInCommit(commit.SHA.Location(), commit.FileName)
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result
}

// CommitStagedChanges commits the currently staged changes.
func (tc *TestCommands) CommitStagedChanges(message string) {
	tc.MustRun("git", "commit", "-m", message)
}

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func (tc *TestCommands) ConnectTrackingBranch(name domain.LocalBranchName) {
	tc.MustRun("git", "branch", "--set-upstream-to=origin/"+name.String(), name.String())
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (tc *TestCommands) DeleteMainBranchConfiguration() {
	tc.MustRun("git", "config", "--unset", config.KeyMainBranch.String())
}

// Fetch retrieves the updates from the origin repo.
func (tc *TestCommands) Fetch() {
	tc.MustRun("git", "fetch")
}

// FileContent provides the current content of a file.
func (tc *TestCommands) FileContent(filename string) string {
	content, err := os.ReadFile(filepath.Join(tc.WorkingDir, filename))
	asserts.NoError(err)
	return string(content)
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (tc *TestCommands) FileContentInCommit(location domain.Location, filename string) string {
	output := tc.MustQuery("git", "show", location.String()+":"+filename)
	if strings.HasPrefix(output, "tree ") {
		// merge commits get an empty file content instead of "tree <SHA>"
		return ""
	}
	return output
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (tc *TestCommands) FilesInCommit(sha domain.SHA) []string {
	output := tc.MustQuery("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha.String())
	return strings.Split(output, "\n")
}

// FilesInBranch provides the list of the files present in the given branch.
func (tc *TestCommands) FilesInBranch(branch domain.LocalBranchName) []string {
	output := tc.MustQuery("git", "ls-tree", "-r", "--name-only", branch.String())
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
func (tc *TestCommands) FilesInBranches(mainBranch domain.LocalBranchName) datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := tc.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	lastBranch := domain.LocalBranchName{}
	for _, branch := range branches {
		files := tc.FilesInBranch(branch)
		for _, file := range files {
			content := tc.FileContentInCommit(branch.Location(), file)
			if branch == lastBranch {
				result.AddRow("", file, content)
			} else {
				result.AddRow(branch.String(), file, content)
			}
			lastBranch = branch
		}
	}
	return result
}

// HasBranchesOutOfSync indicates whether one or more local branches are out of sync with their tracking branch.
func (tc *TestCommands) HasBranchesOutOfSync() bool {
	output := tc.MustQuery("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	return strings.Contains(output, "[")
}

// HasFile indicates whether this repository contains a file with the given name and content.
// An empty error message means a file with the given name and content exists.
func (tc *TestCommands) HasFile(name, content string) string {
	rawContent, err := os.ReadFile(filepath.Join(tc.WorkingDir, name))
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
func (tc *TestCommands) HasGitTownConfigNow() bool {
	output, err := tc.Query("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		return false
	}
	return output != ""
}

// LocalBranches provides the names of all branches in the local repository,
// ordered alphabetically.
func (tc *TestCommands) LocalBranches() (domain.LocalBranchNames, error) {
	output, err := tc.QueryTrim("git", "branch")
	if err != nil {
		return domain.LocalBranchNames{}, err
	}
	result := domain.LocalBranchNames{}
	for _, line := range stringslice.Lines(output) {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, domain.NewLocalBranchName(line))
	}
	return result, nil
}

// LocalBranchesMainFirst provides the names of all local branches in this repo.
func (tc *TestCommands) LocalBranchesMainFirst(mainBranch domain.LocalBranchName) (domain.LocalBranchNames, error) {
	branches, err := tc.LocalBranches()
	if err != nil {
		return domain.LocalBranchNames{}, err
	}
	return slice.Hoist(branches, mainBranch), nil
}

func (tc *TestCommands) MergeBranch(branch domain.LocalBranchName) error {
	return tc.Run("git", "merge", branch.String())
}

func (tc *TestCommands) PushBranch() {
	tc.MustRun("git", "push")
}

func (tc *TestCommands) PushBranchToRemote(branch domain.LocalBranchName, remote domain.Remote) {
	tc.MustRun("git", "push", "-u", remote.String(), branch.String())
}

func (tc *TestCommands) RebaseAgainstBranch(branch domain.LocalBranchName) error {
	return tc.Run("git", "rebase", branch.String())
}

// RemoveBranch deletes the branch with the given name from this repo.
func (tc *TestCommands) RemoveBranch(name domain.LocalBranchName) {
	tc.MustRun("git", "branch", "-D", name.String())
}

// RemoveRemote deletes the Git remote with the given name.
func (tc *TestCommands) RemoveRemote(name domain.Remote) {
	tc.RemotesCache.Invalidate()
	tc.MustRun("git", "remote", "rm", name.String())
}

// RemoveUnnecessaryFiles trims all files that aren'tc necessary in this repo.
func (tc *TestCommands) RemoveUnnecessaryFiles() {
	fullPath := filepath.Join(tc.WorkingDir, ".git", "hooks")
	asserts.NoError(os.RemoveAll(fullPath))
	_ = os.Remove(filepath.Join(tc.WorkingDir, ".git", "COMMIT_EDITMSG"))
	_ = os.Remove(filepath.Join(tc.WorkingDir, ".git", "description"))
}

// SHAForCommit provides the SHA for the commit with the given name.
func (tc *TestCommands) SHAForCommit(name string) string {
	output := tc.MustQuery("git", "log", "--reflog", "--format=%h", "--grep=^"+name+"$")
	if output == "" {
		log.Fatalf("cannot find the SHA of commit %q", name)
	}
	return strings.Split(output, "\n")[0]
}

// StageFiles adds the file with the given name to the Git index.
func (tc *TestCommands) StageFiles(names ...string) {
	args := append([]string{"add"}, names...)
	tc.MustRun("git", args...)
}

// StashOpenFiles stashes the open files away.
func (tc *TestCommands) StashOpenFiles() {
	tc.MustRunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
}

// Tags provides a list of the tags in this repository.
func (tc *TestCommands) Tags() []string {
	output := tc.MustQuery("git", "tag")
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result
}

// UncommittedFiles provides the names of the files not committed into Git.
func (tc *TestCommands) UncommittedFiles() []string {
	output := tc.MustQuery("git", "status", "--porcelain", "--untracked-files=all")
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

func (tc *TestCommands) UnstashOpenFiles() error {
	return tc.Run("git", "stash", "pop")
}
