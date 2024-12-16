package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	prodgit "github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/slice"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/git-town/git-town/v16/test/asserts"
	"github.com/git-town/git-town/v16/test/datatable"
	"github.com/git-town/git-town/v16/test/git"
	"github.com/git-town/git-town/v16/test/subshell"
)

const ConfigFileCommitMessage = "persisted config file"

// TestCommands defines Git commands used only in test code.
type TestCommands struct {
	*subshell.TestRunner
	*prodgit.Commands
	Config config.UnvalidatedConfig
}

// AddRemote adds a Git remote with the given name and URL to this repository.
func (self *TestCommands) AddRemote(name gitdomain.Remote, url string) {
	self.MustRun("git", "remote", "add", name.String(), url)
	self.RemotesCache.Invalidate()
}

// AddSubmodule adds a Git submodule with the given URL to this repository.
func (self *TestCommands) AddSubmodule(url string) {
	self.MustRun("git", "submodule", "add", url)
	self.MustRun("git", "commit", "-m", "added submodule")
}

func (self *TestCommands) AddWorktree(path string, branch gitdomain.LocalBranchName) {
	self.MustRun("git", "worktree", "add", path, branch.String())
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (self *TestCommands) CheckoutBranch(branch gitdomain.LocalBranchName) {
	asserts.NoError(self.Commands.CheckoutBranch(self.TestRunner, branch, false))
}

func (self *TestCommands) CommitSHAs() map[string]gitdomain.SHA {
	result := map[string]gitdomain.SHA{}
	output := self.MustQuery("git", "log", "--all", "--pretty=format:%h %s")
	for _, line := range strings.Split(output, "\n") {
		parts := strings.SplitN(line, " ", 2)
		sha := parts[0]
		commitMessage := parts[1]
		result[commitMessage] = gitdomain.NewSHA(sha)
	}
	return result
}

// CommitStagedChanges commits the currently staged changes.
func (self *TestCommands) CommitStagedChanges(message gitdomain.CommitMessage) {
	self.MustRun("git", "commit", "-m", message.String())
}

// Commits provides a list of the commits in this Git repository with the given fields.
func (self *TestCommands) Commits(fields []string, mainBranch gitdomain.BranchName, lineage configdomain.Lineage) []git.Commit {
	// NOTE: This method uses the provided lineage instead of self.Config.NormalConfig.Lineage
	//       because it might determine the commits on a remote repo, and that repo has no lineage information.
	//       We therefore always provide the lineage of the local repo.
	branches := asserts.NoError1(self.LocalBranchesMainFirst(mainBranch.LocalName()))
	var result []git.Commit
	for _, branch := range branches {
		if strings.HasPrefix(branch.String(), "+ ") {
			// branch is checked out in another workspace --> skip here
			continue
		}
		parent := self.ExistingParent(branch, lineage)
		commits := self.CommitsInBranch(branch, parent, fields)
		result = append(result, commits...)
	}
	return result
}

// CommitsInBranch provides all commits in the given Git branch.
func (self *TestCommands) CommitsInBranch(branch gitdomain.LocalBranchName, parentOpt Option[gitdomain.BranchName], fields []string) []git.Commit {
	args := []string{"log", "--format=%h|%s|%an <%ae>", "--topo-order", "--reverse"}
	if parent, hasParent := parentOpt.Get(); hasParent {
		args = append(args, fmt.Sprintf("%s..%s", parent, branch))
	} else {
		args = append(args, branch.String())
	}
	output := self.MustQuery("git", args...)
	lines := strings.Split(output, "\n")
	result := make([]git.Commit, 0, len(lines))
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		parts := strings.Split(line, "|")
		commit := git.Commit{Branch: branch, SHA: gitdomain.NewSHA(parts[0]), Message: gitdomain.CommitMessage(parts[1]), Author: gitdomain.Author(parts[2])}
		if strings.EqualFold(commit.Message.String(), "initial commit") || strings.EqualFold(commit.Message.String(), ConfigFileCommitMessage) {
			continue
		}
		if slices.Contains(fields, "FILE NAME") {
			filenames := self.FilesInCommit(commit.SHA)
			commit.FileName = strings.Join(filenames, ", ")
		}
		if slices.Contains(fields, "FILE CONTENT") {
			filecontent := self.FileContentInCommit(commit.SHA.Location(), commit.FileName)
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result
}

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func (self *TestCommands) ConnectTrackingBranch(name gitdomain.LocalBranchName) {
	self.MustRun("git", "branch", "--set-upstream-to=origin/"+name.String(), name.String())
}

// creates a feature branch with the given name in this repository
func (self *TestCommands) CreateAndCheckoutFeatureBranch(name gitdomain.LocalBranchName, parent gitdomain.Location) {
	asserts.NoError(self.CreateAndCheckoutBranchWithParent(self, name, parent))
	self.MustRun("git", "config", "git-town-branch."+name.String()+".parent", parent.String())
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *TestCommands) CreateBranch(name gitdomain.LocalBranchName, parent gitdomain.BranchName) {
	self.MustRun("git", "branch", name.String(), parent.String())
}

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func (self *TestCommands) CreateChildFeatureBranch(branch gitdomain.LocalBranchName, parent gitdomain.LocalBranchName) {
	self.CreateBranch(branch, parent.BranchName())
	asserts.NoError(self.Config.NormalConfig.SetParent(branch, parent))
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (self *TestCommands) CreateCommit(commit git.Commit) {
	self.CheckoutBranch(commit.Branch)
	self.CreateFile(commit.FileName, commit.FileContent)
	self.MustRun("git", "add", commit.FileName)
	commands := []string{"commit", "-m", commit.Message.String()}
	if len(commit.Author) > 0 {
		commands = append(commands, "--author="+commit.Author.String())
	}
	self.MustRun("git", commands...)
}

// creates a contribution branches with the given name in this repository
func (self *TestCommands) CreateContributionBranch(name gitdomain.LocalBranchName) {
	self.CreateBranch(name, "main")
	asserts.NoError(self.Config.NormalConfig.AddToContributionBranches(name))
}

// creates a feature branch with the given name in this repository
func (self *TestCommands) CreateFeatureBranch(name gitdomain.LocalBranchName, parent gitdomain.BranchName) {
	self.CreateBranch(name, parent)
	self.MustRun("git", "config", "git-town-branch."+name.String()+".parent", parent.String())
}

// creates a file with the given name and content in this repository
func (self *TestCommands) CreateFile(name, content string) {
	filePath := filepath.Join(self.WorkingDir, name)
	folderPath := filepath.Dir(filePath)
	asserts.NoError(os.MkdirAll(folderPath, os.ModePerm))
	//nolint:gosec // need permission 700 here in order for tests to work
	asserts.NoError(os.WriteFile(filePath, []byte(content), 0o700))
}

// CreateFolder creates a folder with the given name in this repository.
func (self *TestCommands) CreateFolder(name string) {
	folderPath := filepath.Join(self.WorkingDir, name)
	asserts.NoError(os.MkdirAll(folderPath, os.ModePerm))
}

// creates an observed branch with the given name in this repository
func (self *TestCommands) CreateObservedBranch(name gitdomain.LocalBranchName) {
	self.CreateBranch(name, "main")
	asserts.NoError(self.Config.NormalConfig.AddToObservedBranches(name))
}

// creates a parked branch with the given name and parent in this repository
func (self *TestCommands) CreateParkedBranch(name, parent gitdomain.LocalBranchName) {
	self.CreateFeatureBranch(name, parent.BranchName())
	asserts.NoError(self.Config.NormalConfig.AddToParkedBranches(name))
}

// creates a perennial branch with the given name in this repository
func (self *TestCommands) CreatePerennialBranch(name gitdomain.LocalBranchName) {
	self.CreateBranch(name, "main")
	asserts.NoError(self.Config.NormalConfig.AddToPerennialBranches(name))
}

// creates a prototype branch with the given name and parent in this repository
func (self *TestCommands) CreatePrototypeBranch(name, parent gitdomain.LocalBranchName) {
	self.CreateFeatureBranch(name, parent.BranchName())
	asserts.NoError(self.Config.NormalConfig.AddToPrototypeBranches(name))
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

func (self *TestCommands) CurrentCommitMessage() string {
	return self.MustQuery("git", "log", "-1", "--pretty=%B")
}

// provides the first ancestor of the given branch that actually exists in the repo
func (self *TestCommands) ExistingParent(branch gitdomain.LocalBranchName, lineage configdomain.Lineage) Option[gitdomain.BranchName] {
	for {
		parentOpt := lineage.Parent(branch)
		parent, hasParent := parentOpt.Get()
		if !hasParent {
			return None[gitdomain.BranchName]()
		}
		if self.BranchExists(self, parent) {
			return Some(parent.BranchName())
		}
		if self.BranchExistsRemotely(self, parent) {
			return Some(parent.AtRemote(git.RemoteOrigin).BranchName())
		}
		branch = parent
	}
}

// Fetch retrieves the updates from the origin repo.
func (self *TestCommands) Fetch() {
	self.MustRun("git", "fetch")
}

// FileContent provides the current content of a file.
func (self *TestCommands) FileContent(filename string) string {
	return asserts.NoError1(self.FileContentErr(filename))
}

// FileContent provides the current content of a file.
func (self *TestCommands) FileContentErr(filename string) (string, error) {
	content, err := os.ReadFile(filepath.Join(self.WorkingDir, filename))
	return string(content), err
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
func (self *TestCommands) FileContentInCommit(location gitdomain.Location, filename string) string {
	output := self.MustQuery("git", "show", location.String()+":"+filename)
	if strings.HasPrefix(output, "tree ") {
		// merge commits get an empty file content instead of "tree <SHA>"
		return ""
	}
	return output
}

// FilesInBranch provides the list of the files present in the given branch.
func (self *TestCommands) FilesInBranch(branch gitdomain.LocalBranchName) []string {
	output := self.MustQuery("git", "ls-tree", "-r", "--name-only", branch.String())
	var result []string
	for _, line := range strings.Split(output, "\n") {
		file := strings.TrimSpace(line)
		if len(file) > 0 {
			result = append(result, file)
		}
	}
	return result
}

// FilesInBranches provides a data table of files and their content in all branches.
func (self *TestCommands) FilesInBranches(mainBranch gitdomain.LocalBranchName) datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches := asserts.NoError1(self.LocalBranchesMainFirst(mainBranch))
	var lastBranch gitdomain.LocalBranchName
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
func (self *TestCommands) FilesInCommit(sha gitdomain.SHA) []string {
	output := self.MustQuery("git", "show", "--name-only", "--pretty=format:", sha.String())
	return strings.Split(output, "\n")
}

func (self *TestCommands) FilesInWorkspace() []string {
	files := asserts.NoError1(os.ReadDir(self.WorkingDir))
	result := make([]string, 0, len(files))
	for _, file := range files {
		fileName := file.Name()
		if fileName == ".git" {
			continue
		}
		result = append(result, fileName)
	}
	return result
}

func (self *TestCommands) GitConfig(scope configdomain.ConfigScope, name configdomain.Key) Option[string] {
	args := []string{"config"}
	switch scope {
	case configdomain.ConfigScopeGlobal:
		args = append(args, "--global")
	case configdomain.ConfigScopeLocal:
		args = append(args, "--local")
	}
	args = append(args, "--get", name.String())
	output, err := self.Query("git", args...)
	if err != nil {
		return None[string]()
	}
	return Some(output)
}

// HasBranchesOutOfSync indicates whether one or more local branches are out of sync with their tracking branch.
func (self *TestCommands) HasBranchesOutOfSync() (bool, string) {
	output := self.MustQuery("git", "for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	return strings.Contains(output, "["), output
}

// HasFile indicates whether this repository contains a file with the given name and content.
// An empty error message means a file with the given name and content exists.
func (self *TestCommands) HasFile(name, content string) string {
	rawContent, err := os.ReadFile(filepath.Join(self.WorkingDir, name))
	if err != nil {
		return fmt.Sprintf("repo doesn't have file %q", name)
	}
	actualContent := strings.TrimSpace(string(rawContent))
	if actualContent != content {
		return fmt.Sprintf("file %q should have content %q but has %q", name, content, actualContent)
	}
	return ""
}

// LineageTable provides the currently configured lineage information as a DataTable.
func (self *TestCommands) LineageTable() datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("BRANCH", "PARENT")
	_, localGitConfig, _ := self.Config.NormalConfig.GitConfigAccess.LoadLocal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	lineage := localGitConfig.Lineage
	for _, entry := range lineage.Entries() {
		result.AddRow(entry.Child.String(), entry.Parent.String())
	}
	result.Sort()
	return result
}

// SetGitAlias sets the Git alias with the given name to the given value.
func (self *TestCommands) LoadGitAlias(name configdomain.AliasableCommand) (string, error) {
	return self.Query("git", "config", "--get", "--global", configdomain.AliasKeyPrefix+name.String())
}

// LocalBranches provides the names of all branches in the local repository,
// ordered alphabetically.
func (self *TestCommands) LocalBranches() (gitdomain.LocalBranchNames, error) {
	output, err := self.QueryTrim("git", "branch")
	if err != nil {
		return gitdomain.LocalBranchNames{}, err
	}
	result := gitdomain.LocalBranchNames{}
	for _, line := range stringslice.Lines(output) {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, gitdomain.NewLocalBranchName(line))
	}
	return result, nil
}

// LocalBranchesMainFirst provides the names of all local branches in this repo.
func (self *TestCommands) LocalBranchesMainFirst(mainBranch gitdomain.LocalBranchName) (gitdomain.LocalBranchNames, error) {
	branches, err := self.LocalBranches()
	if err != nil {
		return gitdomain.LocalBranchNames{}, err
	}
	branches = slice.Hoist(branches, mainBranch)
	return branches, nil
}

func (self *TestCommands) MergeBranch(branch gitdomain.LocalBranchName) error {
	return self.Run("git", "merge", branch.String())
}

func (self *TestCommands) PushBranch() {
	self.MustRun("git", "push")
}

func (self *TestCommands) PushBranchToRemote(branch gitdomain.LocalBranchName, remote gitdomain.Remote) {
	self.MustRun("git", "push", "-u", remote.String(), branch.String())
}

func (self *TestCommands) RebaseAgainstBranch(branch gitdomain.LocalBranchName) error {
	return self.Run("git", "rebase", branch.String())
}

// RemoveBranch deletes the branch with the given name from this repo.
func (self *TestCommands) RemoveBranch(name gitdomain.LocalBranchName) {
	self.MustRun("git", "branch", "-D", name.String())
}

// DeleteMainBranchConfiguration removes the configuration for which branch is the main branch.
func (self *TestCommands) RemoveMainBranchConfiguration() {
	self.MustRun("git", "config", "--unset", configdomain.KeyMainBranch.String())
}

// RemovePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (self *TestCommands) RemovePerennialBranchConfiguration() error {
	return self.Config.NormalConfig.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyPerennialBranches)
}

// RemoveRemote deletes the Git remote with the given name.
func (self *TestCommands) RemoveRemote(name gitdomain.Remote) {
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

func (self *TestCommands) RenameRemote(oldName, newName string) {
	self.RemotesCache.Invalidate()
	self.MustRun("git", "remote", "rename", oldName, newName)
}

// SHAForCommit provides the SHA for the commit with the given name.
func (self *TestCommands) SHAsForCommit(name string) gitdomain.SHAs {
	output := self.MustQuery("git", "reflog", "--format=%h %s")
	if output == "" {
		panic(fmt.Sprintf("cannot find the SHA of commit %q", name))
	}
	shasWithMessage := make(gitdomain.SHAs, 0, 1)
	for _, text := range strings.Split(output, "\n") {
		shaText, commitMessage, found := strings.Cut(text, " ")
		if found && commitMessage == name {
			sha := gitdomain.NewSHA(shaText)
			shasWithMessage = append(shasWithMessage, sha)
		}
	}
	return shasWithMessage
}

// SetColorUI configures whether Git output contains color codes.
func (self *TestCommands) SetColorUI(value string) error {
	return self.Run("git", "config", "color.ui", value)
}

func (self *TestCommands) SetDefaultGitBranch(value gitdomain.LocalBranchName) {
	self.MustRun("git", "config", "init.defaultbranch", value.String())
}

// SetGitAlias sets the Git alias with the given name to the given value.
func (self *TestCommands) SetGitAlias(name configdomain.AliasableCommand, value string) error {
	return self.Run("git", "config", "--global", configdomain.AliasKeyPrefix+name.String(), value)
}

// StageFiles adds the file with the given name to the Git index.
func (self *TestCommands) StageFiles(names ...string) {
	args := append([]string{"add"}, names...)
	self.MustRun("git", args...)
}

// StashOpenFiles stashes the open files away.
func (self *TestCommands) StashOpenFiles() {
	self.MustRun("git", "add", "-A")
	self.MustRun("git", "stash")
}

// Tags provides a list of the tags in this repository.
func (self *TestCommands) Tags() []string {
	output := self.MustQuery("git", "tag")
	lines := strings.Split(output, "\n")
	result := make([]string, len(lines))
	for l, line := range lines {
		result[l] = strings.TrimSpace(line)
	}
	return result
}

// UncommittedFiles provides the names of the files not committed into Git.
func (self *TestCommands) UncommittedFiles() []string {
	output := self.MustQuery("git", "status", "--porcelain", "--untracked-files=all")
	lines := stringslice.Lines(output)
	result := make([]string, 0, len(lines))
	for _, line := range lines {
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

// HasGitTownConfigNow indicates whether this repository contain Git Town specific configuration.
func (self *TestCommands) VerifyNoGitTownConfiguration() error {
	output, _ := self.Query("git", "config", "--get-regex", "git-town")
	if len(output) > 0 {
		return fmt.Errorf("unexpected Git Town configuration:\n%s", output)
	}
	output, _ = self.Query("git", "config", "--get-regex", "git-town-branch")
	if len(output) > 0 {
		return fmt.Errorf("unexpected Git Town configuration:\n%s", output)
	}
	self.Config.Reload()
	for aliasName, aliasValue := range self.Config.NormalConfig.Aliases {
		if strings.HasPrefix(aliasValue, "town ") {
			return fmt.Errorf("unexpected Git Town alias %q with value %q. All aliases: %#v", aliasName, aliasValue, self.Config.NormalConfig.Aliases)
		}
	}
	return nil
}
