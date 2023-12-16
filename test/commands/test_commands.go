package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	prodgit "github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/git-town/git-town/v11/src/gohacks/stringslice"
	"github.com/git-town/git-town/v11/test/asserts"
	"github.com/git-town/git-town/v11/test/datatable"
	"github.com/git-town/git-town/v11/test/git"
	"github.com/git-town/git-town/v11/test/subshell"
)

// TestCommands defines Git commands used only in test code.
type TestCommands struct {
	*subshell.TestRunner
	*prodgit.BackendCommands // TODO: remove this dependency on BackendCommands
}

// AddRemote adds a Git remote with the given name and URL to this repository.
func (self *TestCommands) AddRemote(name domain.Remote, url string) {
	self.MustRun("git", "remote", "add", name.String(), url)
	self.RemotesCache.Invalidate()
}

// AddSubmodule adds a Git submodule with the given URL to this repository.
func (self *TestCommands) AddSubmodule(url string) {
	self.MustRun("git", "submodule", "add", url)
	self.MustRun("git", "commit", "-m", "added submodule")
}

func (self *TestCommands) AddWorktree(path string, branch domain.LocalBranchName) {
	self.MustRun("git", "worktree", "add", path, branch.String())
}

// .CheckoutBranch checks out the Git branch with the given name in this repo.
func (self *TestCommands) CheckoutBranch(branch domain.LocalBranchName) {
	asserts.NoError(self.BackendCommands.CheckoutBranch(branch))
}

func (self *TestCommands) CommitSHAs() map[string]domain.SHA {
	result := map[string]domain.SHA{}
	output := self.MustQuery("git", "log", "--all", "--pretty=format:%h %s")
	for _, line := range strings.Split(output, "\n") {
		parts := strings.SplitN(line, " ", 2)
		sha := parts[0]
		commitMessage := parts[1]
		result[commitMessage] = domain.NewSHA(sha)
	}
	return result
}

// CommitStagedChanges commits the currently staged changes.
func (self *TestCommands) CommitStagedChanges(message string) {
	self.MustRun("git", "commit", "-m", message)
}

// Commits provides a list of the commits in this Git repository with the given fields.
func (self *TestCommands) Commits(fields []string, mainBranch domain.LocalBranchName) []git.Commit {
	branches, err := self.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	result := []git.Commit{}
	for _, branch := range branches {
		if strings.HasPrefix(branch.String(), "+ ") {
			continue
		}
		commits := self.CommitsInBranch(branch, fields)
		result = append(result, commits...)
	}
	return result
}

// CommitsInBranch provides all commits in the given Git branch.
func (self *TestCommands) CommitsInBranch(branch domain.LocalBranchName, fields []string) []git.Commit {
	output := self.MustQuery("git", "log", branch.String(), "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse")
	result := []git.Commit{}
	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, "|")
		commit := git.Commit{Branch: branch, SHA: domain.NewSHA(parts[0]), Message: parts[1], Author: parts[2]}
		if strings.EqualFold(commit.Message, "initial commit") {
			continue
		}
		if slice.Contains(fields, "FILE NAME") {
			filenames := self.FilesInCommit(commit.SHA)
			commit.FileName = strings.Join(filenames, ", ")
		}
		if slice.Contains(fields, "FILE CONTENT") {
			filecontent := self.FileContentInCommit(commit.SHA.Location(), commit.FileName)
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result
}

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func (self *TestCommands) ConnectTrackingBranch(name domain.LocalBranchName) {
	self.MustRun("git", "branch", "--set-upstream-to=origin/"+name.String(), name.String())
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *TestCommands) CreateBranch(name, parent domain.LocalBranchName) {
	self.MustRun("git", "branch", name.String(), parent.String())
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (self *TestCommands) CreateChildFeatureBranch(branch domain.LocalBranchName, parent domain.LocalBranchName) {
	self.CreateBranch(branch, parent)
	asserts.NoError(self.GitTown.SetParent(branch, parent))
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (self *TestCommands) CreateCommit(commit git.Commit) {
	self.CheckoutBranch(commit.Branch)
	self.CreateFile(commit.FileName, commit.FileContent)
	self.MustRun("git", "add", commit.FileName)
	commands := []string{"commit", "-m", commit.Message}
	if commit.Author != "" {
		commands = append(commands, "--author="+commit.Author)
	}
	self.MustRun("git", commands...)
}

// CreateFile creates a file with the given name and content in this repository.
func (self *TestCommands) CreateFile(name, content string) {
	filePath := filepath.Join(self.WorkingDir, name)
	folderPath := filepath.Dir(filePath)
	asserts.NoError(os.MkdirAll(folderPath, os.ModePerm))
	//nolint:gosec // need permission 700 here in order for tests to work
	asserts.NoError(os.WriteFile(filePath, []byte(content), 0o700))
}

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func (self *TestCommands) CreatePerennialBranches(names ...domain.LocalBranchName) {
	for _, name := range names {
		self.CreateBranch(name, domain.NewLocalBranchName("main"))
	}
	asserts.NoError(self.GitTown.AddToPerennialBranches(names...))
}

// CreateStandaloneTag creates a tag not on a branch.
func (self *TestCommands) CreateStandaloneTag(name string) {
	self.MustRun("git", "checkout", "-b", "temp")
	self.MustRun("touch", "a.txt")
	self.MustRun("git", "add", "-A")
	self.MustRun("git", "commit", "-m", "temp")
	self.MustRun("git", "tag", "-a", name, "-m", name)
	self.MustRun("git", "checkout", "-")
	self.MustRun("git", "branch", "-D", "temp")
}

// CreateTag creates a tag with the given name.
func (self *TestCommands) CreateTag(name string) {
	self.MustRun("git", "tag", "-a", name, "-m", name)
}

// Fetch retrieves the updates from the origin repo.
func (self *TestCommands) Fetch() {
	self.MustRun("git", "fetch")
}

// FileContent provides the current content of a file.
func (self *TestCommands) FileContent(filename string) string {
	content, err := os.ReadFile(filepath.Join(self.WorkingDir, filename))
	asserts.NoError(err)
	return string(content)
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (self *TestCommands) FileContentInCommit(location domain.Location, filename string) string {
	output := self.MustQuery("git", "show", location.String()+":"+filename)
	if strings.HasPrefix(output, "tree ") {
		// merge commits get an empty file content instead of "tree <SHA>"
		return ""
	}
	return output
}

// FilesInBranch provides the list of the files present in the given branch.
func (self *TestCommands) FilesInBranch(branch domain.LocalBranchName) []string {
	output := self.MustQuery("git", "ls-tree", "-r", "--name-only", branch.String())
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
func (self *TestCommands) FilesInBranches(mainBranch domain.LocalBranchName) datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, err := self.LocalBranchesMainFirst(mainBranch)
	asserts.NoError(err)
	lastBranch := domain.EmptyLocalBranchName()
	for _, branch := range branches {
		files := self.FilesInBranch(branch)
		for _, file := range files {
			content := self.FileContentInCommit(branch.Location(), file)
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

// FilesInCommit provides the names of the files that the commit with the given SHA changes.
func (self *TestCommands) FilesInCommit(sha domain.SHA) []string {
	output := self.MustQuery("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha.String())
	return strings.Split(output, "\n")
}

// HasBranchesOutOfSync indicates whether one or more local branches are out of sync with their tracking branch.
func (self *TestCommands) HasBranchesOutOfSync() bool {
	output := self.MustQuery("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	return strings.Contains(output, "[")
}

// HasFile indicates whether this repository contains a file with the given name and content.
// An empty error message means a file with the given name and content exists.
func (self *TestCommands) HasFile(name, content string) string {
	rawContent, err := os.ReadFile(filepath.Join(self.WorkingDir, name))
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
func (self *TestCommands) HasGitTownConfigNow() bool {
	output, err := self.Query("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		return false
	}
	if output != "" {
		return true
	}
	output, err = self.Query("git", "config", "--local", "--get-regex", "git-town-branch")
	if err != nil {
		return false
	}
	if output != "" {
		return true
	}
	return false
}

// LineageTable provides the currently configured lineage information as a DataTable.
func (self *TestCommands) LineageTable() datatable.DataTable {
	result := datatable.DataTable{}
	self.GitTown.Reload()
	lineage := self.GitTown.Lineage(self.GitTown.RemoveLocalConfigValue)
	result.AddRow("BRANCH", "PARENT")
	for _, branchName := range lineage.BranchNames() {
		result.AddRow(branchName.String(), lineage[branchName].String())
	}
	return result
}

// LocalBranches provides the names of all branches in the local repository,
// ordered alphabetically.
func (self *TestCommands) LocalBranches() (domain.LocalBranchNames, error) {
	output, err := self.QueryTrim("git", "branch")
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
func (self *TestCommands) LocalBranchesMainFirst(mainBranch domain.LocalBranchName) (domain.LocalBranchNames, error) {
	branches, err := self.LocalBranches()
	if err != nil {
		return domain.LocalBranchNames{}, err
	}
	slice.Hoist(&branches, mainBranch)
	return branches, nil
}

func (self *TestCommands) MergeBranch(branch domain.LocalBranchName) error {
	return self.Run("git", "merge", branch.String())
}

func (self *TestCommands) PushBranch() {
	self.MustRun("git", "push")
}

func (self *TestCommands) PushBranchToRemote(branch domain.LocalBranchName, remote domain.Remote) {
	self.MustRun("git", "push", "-u", remote.String(), branch.String())
}

func (self *TestCommands) RebaseAgainstBranch(branch domain.LocalBranchName) error {
	return self.Run("git", "rebase", branch.String())
}

// RemoveBranch deletes the branch with the given name from this repo.
func (self *TestCommands) RemoveBranch(name domain.LocalBranchName) {
	self.MustRun("git", "branch", "-D", name.String())
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (self *TestCommands) RemoveMainBranchConfiguration() {
	self.MustRun("git", "config", "--unset", configdomain.KeyMainBranch.String())
}

// RemovePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (self *TestCommands) RemovePerennialBranchConfiguration() error {
	return self.RemoveLocalConfigValue(configdomain.KeyPerennialBranches)
}

// RemoveRemote deletes the Git remote with the given name.
func (self *TestCommands) RemoveRemote(name domain.Remote) {
	self.RemotesCache.Invalidate()
	self.MustRun("git", "remote", "rm", name.String())
}

// RemoveUnnecessaryFiles trims all files that aren'self necessary in this repo.
func (self *TestCommands) RemoveUnnecessaryFiles() {
	fullPath := filepath.Join(self.WorkingDir, ".git", "hooks")
	asserts.NoError(os.RemoveAll(fullPath))
	_ = os.Remove(filepath.Join(self.WorkingDir, ".git", "COMMIT_EDITMSG"))
	_ = os.Remove(filepath.Join(self.WorkingDir, ".git", "description"))
}

// SHAForCommit provides the SHA for the commit with the given name.
// TODO: return a domain.SHA here.
func (self *TestCommands) SHAForCommit(name string) string {
	output := self.MustQuery("git", "log", "--reflog", "--format=%h", "--grep=^"+name+"$")
	if output == "" {
		log.Fatalf("cannot find the SHA of commit %q", name)
	}
	return strings.Split(output, "\n")[0]
}

// SetColorUI configures whether Git output contains color codes.
func (self *TestCommands) SetColorUI(value string) error {
	return self.Run("git", "config", "color.ui", value)
}

// StageFiles adds the file with the given name to the Git index.
func (self *TestCommands) StageFiles(names ...string) {
	args := append([]string{"add"}, names...)
	self.MustRun("git", args...)
}

// StashOpenFiles stashes the open files away.
func (self *TestCommands) StashOpenFiles() {
	self.MustRunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
}

// Tags provides a list of the tags in this repository.
func (self *TestCommands) Tags() []string {
	output := self.MustQuery("git", "tag")
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result
}

// UncommittedFiles provides the names of the files not committed into Git.
func (self *TestCommands) UncommittedFiles() []string {
	output := self.MustQuery("git", "status", "--porcelain", "--untracked-files=all")
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

func (self *TestCommands) UnstashOpenFiles() error {
	return self.Run("git", "stash", "pop")
}
