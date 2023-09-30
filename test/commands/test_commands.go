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
func (t *TestCommands) AddRemote(name domain.Remote, url string) {
	t.MustRun("git", "remote", "add", name.String(), url)
	t.RemotesCache.Invalidate()
}

// AddSubmodule adds a Git submodule with the given URL to this repository.
func (t *TestCommands) AddSubmodule(url string) {
	t.MustRun("git", "submodule", "add", url)
	t.MustRun("git", "commit", "-m", "added submodule")
}

// BranchHierarchyTable provides the currently configured branch hierarchy information as a DataTable.
func (t *TestCommands) BranchHierarchyTable() datatable.DataTable {
	result := datatable.DataTable{}
	t.Config.Reload()
	lineage := t.Config.Lineage()
	result.AddRow("BRANCH", "PARENT")
	for _, branchName := range lineage.BranchNames() {
		result.AddRow(branchName.String(), lineage[branchName].String())
	}
	return result
}

// .CheckoutBranch checks out the Git branch with the given name in this repo.
func (t *TestCommands) CheckoutBranch(branch domain.LocalBranchName) {
	asserts.NoError(t.BackendCommands.CheckoutBranch(branch))
}

func (t *TestCommands) CommitSHAs() map[string]domain.SHA {
	result := map[string]domain.SHA{}
	output := t.MustQuery("git", "log", "--all", "--pretty=format:%h %s")
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
func (t *TestCommands) CreateBranch(name, parent domain.LocalBranchName) {
	t.MustRun("git", "branch", name.String(), parent.String())
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (t *TestCommands) CreateChildFeatureBranch(branch domain.LocalBranchName, parent domain.LocalBranchName) {
	t.CreateBranch(branch, parent)
	asserts.NoError(t.Config.SetParent(branch, parent))
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (t *TestCommands) CreateCommit(commit git.Commit) {
	t.CheckoutBranch(commit.Branch)
	t.CreateFile(commit.FileName, commit.FileContent)
	t.MustRun("git", "add", commit.FileName)
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	t.MustRun("git", commands...)
}

// CreateFile creates a file with the given name and content in this repository.
func (t *TestCommands) CreateFile(name, content string) {
	filePath := filepath.Join(t.WorkingDir, name)
	folderPath := filepath.Dir(filePath)
	asserts.NoError(os.MkdirAll(folderPath, os.ModePerm))
	//nolint:gosec // need permission 700 here in order for tests to work
	asserts.NoError(os.WriteFile(filePath, []byte(content), 0x700))
}

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func (t *TestCommands) CreatePerennialBranches(names ...domain.LocalBranchName) {
	for _, name := range names {
		t.CreateBranch(name, domain.NewLocalBranchName("main"))
	}
	asserts.NoError(t.Config.AddToPerennialBranches(names...))
}

// CreateStandaloneTag creates a tag not on a branch.
func (t *TestCommands) CreateStandaloneTag(name string) {
	t.MustRun("git", "checkout", "-b", "temp")
	t.MustRun("touch", "a.txt")
	t.MustRun("git", "add", "-A")
	t.MustRun("git", "commit", "-m", "temp")
	t.MustRun("git", "tag", "-a", name, "-m", name)
	t.MustRun("git", "checkout", "-")
	t.MustRun("git", "branch", "-D", "temp")
}

// CreateTag creates a tag with the given name.
func (t *TestCommands) CreateTag(name string) {
	t.MustRun("git", "tag", "-a", name, "-m", name)
}

// Commits provides a list of the commits in this Git repository with the given fields.
func (t *TestCommands) Commits(fields []string, mainBranch domain.LocalBranchName) []git.Commit {
	branches, err := t.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	result := []git.Commit{}
	for _, branch := range branches {
		commits := t.CommitsInBranch(branch, fields)
		result = append(result, commits...)
	}
	return result
}

// CommitsInBranch provides all commits in the given Git branch.
func (t *TestCommands) CommitsInBranch(branch domain.LocalBranchName, fields []string) []git.Commit {
	output := t.MustQuery("git", "log", branch.String(), "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	result := []git.Commit{}
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, "|")
		commit := git.Commit{Branch: branch, SHA: domain.NewSHA(parts[0]), Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if slice.Contains(fields, "FILE NAME") {
			filenames := t.FilesInCommit(commit.SHA)
			commit.FileName = strings.Join(filenames, ", ")
		}
		if slice.Contains(fields, "FILE CONTENT") {
			filecontent := t.FileContentInCommit(commit.SHA.Location(), commit.FileName)
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result
}

// CommitStagedChanges commits the currently staged changes.
func (t *TestCommands) CommitStagedChanges(message string) {
	t.MustRun("git", "commit", "-m", message)
}

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func (t *TestCommands) ConnectTrackingBranch(name domain.LocalBranchName) {
	t.MustRun("git", "branch", "--set-upstream-to=origin/"+name.String(), name.String())
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (t *TestCommands) DeleteMainBranchConfiguration() {
	t.MustRun("git", "config", "--unset", config.KeyMainBranch.String())
}

// Fetch retrieves the updates from the origin repo.
func (t *TestCommands) Fetch() {
	t.MustRun("git", "fetch")
}

// FileContent provides the current content of a file.
func (t *TestCommands) FileContent(filename string) string {
	content, err := os.ReadFile(filepath.Join(t.WorkingDir, filename))
	asserts.NoError(err)
	return string(content)
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (t *TestCommands) FileContentInCommit(location domain.Location, filename string) string {
	output := t.MustQuery("git", "show", location.String()+":"+filename)
	if strings.HasPrefix(output, "tree ") {
		// merge commits get an empty file content instead of "tree <SHA>"
		return ""
	}
	return output
}

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (t *TestCommands) FilesInCommit(sha domain.SHA) []string {
	output := t.MustQuery("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha.String())
	return strings.Split(output, "\n")
}

// FilesInBranch provides the list of the files present in the given branch.
func (t *TestCommands) FilesInBranch(branch domain.LocalBranchName) []string {
	output := t.MustQuery("git", "ls-tree", "-r", "--name-only", branch.String())
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
func (t *TestCommands) FilesInBranches(mainBranch domain.LocalBranchName) datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := t.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	lastBranch := domain.LocalBranchName{}
	for _, branch := range branches {
		files := t.FilesInBranch(branch)
		for _, file := range files {
			content := t.FileContentInCommit(branch.Location(), file)
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
func (t *TestCommands) HasBranchesOutOfSync() bool {
	output := t.MustQuery("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	return strings.Contains(output, "[")
}

// HasFile indicates whether this repository contains a file with the given name and content.
// An empty error message means a file with the given name and content exists.
func (t *TestCommands) HasFile(name, content string) string {
	rawContent, err := os.ReadFile(filepath.Join(t.WorkingDir, name))
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
func (t *TestCommands) HasGitTownConfigNow() bool {
	output, err := t.Query("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		return false
	}
	return output != ""
}

// LocalBranches provides the names of all branches in the local repository,
// ordered alphabetically.
func (t *TestCommands) LocalBranches() (domain.LocalBranchNames, error) {
	output, err := t.QueryTrim("git", "branch")
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
func (t *TestCommands) LocalBranchesMainFirst(mainBranch domain.LocalBranchName) (domain.LocalBranchNames, error) {
	branches, err := t.LocalBranches()
	if err != nil {
		return domain.LocalBranchNames{}, err
	}
	return slice.Hoist(branches, mainBranch), nil
}

func (t *TestCommands) MergeBranch(branch domain.LocalBranchName) error {
	return t.Run("git", "merge", branch.String())
}

func (t *TestCommands) PushBranch() {
	t.MustRun("git", "push")
}

func (t *TestCommands) PushBranchToRemote(branch domain.LocalBranchName, remote domain.Remote) {
	t.MustRun("git", "push", "-u", remote.String(), branch.String())
}

func (t *TestCommands) RebaseAgainstBranch(branch domain.LocalBranchName) error {
	return t.Run("git", "rebase", branch.String())
}

// RemoveBranch deletes the branch with the given name from this repo.
func (t *TestCommands) RemoveBranch(name domain.LocalBranchName) {
	t.MustRun("git", "branch", "-D", name.String())
}

// RemoveRemote deletes the Git remote with the given name.
func (t *TestCommands) RemoveRemote(name domain.Remote) {
	t.RemotesCache.Invalidate()
	t.MustRun("git", "remote", "rm", name.String())
}

// RemoveUnnecessaryFiles trims all files that aren't necessary in this repo.
func (t *TestCommands) RemoveUnnecessaryFiles() {
	fullPath := filepath.Join(t.WorkingDir, ".git", "hooks")
	asserts.NoError(os.RemoveAll(fullPath))
	_ = os.Remove(filepath.Join(t.WorkingDir, ".git", "COMMIT_EDITMSG"))
	_ = os.Remove(filepath.Join(t.WorkingDir, ".git", "description"))
}

// SHAForCommit provides the SHA for the commit with the given name.
func (t *TestCommands) SHAForCommit(name string) string {
	output := t.MustQuery("git", "log", "--reflog", "--format=%h", "--grep=^"+name+"$")
	if output == "" {
		log.Fatalf("cannot find the SHA of commit %q", name)
	}
	return strings.Split(output, "\n")[0]
}

// StageFiles adds the file with the given name to the Git index.
func (t *TestCommands) StageFiles(names ...string) {
	args := append([]string{"add"}, names...)
	t.MustRun("git", args...)
}

// StashOpenFiles stashes the open files away.
func (t *TestCommands) StashOpenFiles() {
	t.MustRunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
}

// Tags provides a list of the tags in this repository.
func (t *TestCommands) Tags() []string {
	output := t.MustQuery("git", "tag")
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result
}

// UncommittedFiles provides the names of the files not committed into Git.
func (t *TestCommands) UncommittedFiles() []string {
	output := t.MustQuery("git", "status", "--porcelain", "--untracked-files=all")
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

func (t *TestCommands) UnstashOpenFiles() error {
	return t.Run("git", "stash", "pop")
}
