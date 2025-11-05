package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/git-town/git-town/v22/internal/cli/format"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	prodgit "github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/mapstools"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/internal/test/datatable"
	"github.com/git-town/git-town/v22/internal/test/subshell"
	"github.com/git-town/git-town/v22/internal/test/testgit"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	ConfigFileCommitMessage = "persisted config file"
	FileCommitMessage       = "persisted file"
	MainBranchName          = "main"
	deletedText             = " (deleted)"
)

// TestCommands defines Git commands used only in test code.
type TestCommands struct {
	*subshell.TestRunner
	Config    config.UnvalidatedConfig
	Git       *prodgit.Commands
	SnapShots configdomain.BeginConfigSnapshot // copy of the low-level Git config data, for verifying it in end-to-end tests
}

// AddRemote adds a Git remote with the given name and URL to this repository.
func (self *TestCommands) AddRemote(name gitdomain.Remote, url string) {
	self.MustRun("git", "remote", "add", name.String(), url)
	self.Git.RemotesCache.Invalidate()
}

// AddSubmodule adds a Git submodule with the given URL to this repository.
func (self *TestCommands) AddSubmodule(url string) {
	self.MustRun("git", "submodule", "add", url)
	self.MustRun("git", "commit", "-m", "added submodule")
}

func (self *TestCommands) AddWorktree(path string, branch gitdomain.LocalBranchName) {
	self.MustRun("git", "worktree", "add", path, branch.String())
}

func (self *TestCommands) AmendCommit() {
	self.MustRun("git", "commit", "--amend", "--no-edit")
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (self *TestCommands) CheckoutBranch(branch gitdomain.LocalBranchName) {
	asserts.NoError(self.Git.CheckoutBranchUncached(self.TestRunner, branch, false))
}

func (self *TestCommands) CommitSHA(querier subshelldomain.Querier, title gitdomain.CommitTitle, branch gitdomain.LocalBranchName, parent gitdomain.BranchName) gitdomain.SHA {
	commits := asserts.NoError1(self.Git.CommitsInFeatureBranch(querier, branch, parent))
	for _, commit := range commits {
		if commit.Message.Parts().Title == title {
			return commit.SHA
		}
	}
	panic(fmt.Errorf("no commit with title %q found", title))
}

func (self *TestCommands) CommitSHAs() gitdomain.Commits {
	output := self.MustQuery("git", "log", "--all", "--format=%H %s")
	if output == "" {
		return gitdomain.Commits{}
	}
	lines := stringslice.NonEmptyLines(output)
	result := make(gitdomain.Commits, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		sha, commitMessage, found := strings.Cut(line, " ")
		if !found {
			// No space found means the entire line is the SHA (empty commit message)
			sha = line
			commitMessage = ""
		}
		result = append(result, gitdomain.Commit{
			Message: gitdomain.CommitMessage(commitMessage),
			SHA:     gitdomain.NewSHA(sha),
		})
	}
	return result
}

// CommitStagedChanges commits the currently staged changes.
func (self *TestCommands) CommitStagedChanges(message gitdomain.CommitMessage) {
	self.MustRun("git", "commit", "-m", message.String())
}

// Commits provides a list of the commits in this Git repository with the given fields.
func (self *TestCommands) Commits(fields []string, lineage configdomain.Lineage, order configdomain.Order) []testgit.Commit {
	// NOTE: This method uses the provided lineage instead of self.Config.NormalConfig.Lineage
	//       because it might determine the commits on a remote repo, and that repo has no lineage information.
	//       We therefore always provide the lineage of the local repo.
	branches, branchesInOtherWorktree := asserts.NoError2(self.LocalBranchesOrderedHierarchically(lineage, order))
	var result []testgit.Commit
	for _, branch := range branches {
		if slices.Contains(branchesInOtherWorktree, branch) {
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
func (self *TestCommands) CommitsInBranch(branch gitdomain.LocalBranchName, parentOpt Option[gitdomain.BranchName], fields []string) []testgit.Commit {
	args := []string{"log", "--format=%H%x00%s%x00%an <%ae>", "--topo-order", "--reverse"}
	if parent, hasParent := parentOpt.Get(); hasParent {
		args = append(args, parent.RefName()+".."+branch.RefName())
	} else {
		args = append(args, branch.RefName())
	}
	output := self.MustQuery("git", args...)
	lines := stringslice.NonEmptyLines(output)
	result := make([]testgit.Commit, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, "\x00")
		commit := testgit.Commit{
			Branch:  branch,
			SHA:     gitdomain.NewSHA(parts[0]),
			Message: gitdomain.CommitMessage(parts[1]),
			Author:  gitdomain.Author(parts[2]),
		}
		if strings.EqualFold(commit.Message.String(), "initial commit") || strings.EqualFold(commit.Message.String(), ConfigFileCommitMessage) {
			continue
		}
		if slices.Contains(fields, "FILE NAME") {
			filenames := self.FilesInCommit(commit.SHA)
			commit.FileName = strings.Join(filenames, ", ")
		}
		if slices.Contains(fields, "FILE CONTENT") {
			filecontent := ""
			if commit.FileName != "" {
				var deleted bool
				filecontent, deleted = self.FileContentInCommit(commit.SHA.Location(), commit.FileName)
				if deleted {
					commit.FileName += deletedText
				}
			}
			commit.FileContent = filecontent
		}
		result = append(result, commit)
	}
	return result
}

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func (self *TestCommands) ConnectTrackingBranch(name gitdomain.LocalBranchName) {
	self.MustRun("git", "branch", fmt.Sprintf("--set-upstream-to=origin/%s", name), name.String())
}

// CreateAndCheckoutBranch creates a new branch with the given name and checks it out.
// The created branch is a normal branch.
// To create feature branches, use CreateAndCheckoutFeatureBranch.
func (self *TestCommands) CreateAndCheckoutBranch(name gitdomain.LocalBranchName, parent gitdomain.BranchName) {
	self.Git.CurrentBranchCache.Set(name)
	self.MustRun("git", "checkout", "-b", name.String(), parent.String())
}

// CreateAndCheckoutFeatureBranch creates a feature branch with the given name and checks it out.
func (self *TestCommands) CreateAndCheckoutFeatureBranch(name gitdomain.LocalBranchName, parent gitdomain.LocalBranchName) {
	asserts.NoError(self.Git.CreateAndCheckoutBranchWithParent(self, name, parent.Location()))
	asserts.NoError(gitconfig.SetParent(self, name, parent))
}

// CreateBranch creates a new branch with the given name without checking it out.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *TestCommands) CreateBranch(name gitdomain.LocalBranchName, parent gitdomain.BranchName) {
	self.MustRun("git", "branch", name.String(), parent.String())
}

func (self *TestCommands) CreateBranchOfType(name gitdomain.LocalBranchName, parentOpt Option[gitdomain.LocalBranchName], branchType configdomain.BranchType) {
	if parent, hasParent := parentOpt.Get(); hasParent {
		self.CreateFeatureBranch(name, parent.BranchName())
	} else {
		self.CreateBranch(name, MainBranchName)
	}
	asserts.NoError(gitconfig.SetBranchTypeOverride(self.TestRunner, branchType, name))
}

func (self *TestCommands) CreateChildBranch(name, parent gitdomain.LocalBranchName) {
	self.CheckoutBranch(parent)
	env := append(os.Environ(), "GIT_TOWN_NEW_BRANCH_TYPE=feature")
	asserts.NoError(self.RunWithEnv(env, "git-town", "append", "--no-sync", name.String()))
	self.Git.CurrentBranchCache.Set(name)
}

// CreateCommit creates a commit with the given properties in this Git repo.
func (self *TestCommands) CreateCommit(commit testgit.Commit) {
	self.CheckoutBranch(commit.Branch)
	self.CreateFile(commit.FileName, commit.FileContent)
	self.MustRun("git", "add", commit.FileName)
	commands := []string{"commit", "-m", commit.Message.String()}
	if len(commit.Author) > 0 {
		commands = append(commands, "--author="+commit.Author.String())
	}
	self.MustRun("git", commands...)
}

// CreateFeatureBranch creates a feature branch with the given name without checking it out.
func (self *TestCommands) CreateFeatureBranch(name gitdomain.LocalBranchName, parent gitdomain.BranchName) {
	self.CreateBranch(name, parent)
	asserts.NoError(self.Config.NormalConfig.SetParent(self.TestRunner, name, parent.LocalName()))
}

// CreateFile creates a file with the given name and content in this repository.
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

func (self *TestCommands) CreateLocalBranchUsingGitTown(branchSetup datatable.BranchSetup) {
	// step 1: create the local branch if one is needed
	if parent, hasParent := branchSetup.Parent.Get(); hasParent {
		self.CreateChildBranch(branchSetup.Name, parent)
	} else {
		self.CreateAndCheckoutBranch(branchSetup.Name, MainBranchName)
	}
	// step 2: create the tracking branch
	if branchSetup.Locations.Contains(testgit.LocationOrigin) {
		self.PushBranchToRemote(branchSetup.Name, gitdomain.RemoteOrigin)
	}
	if branchSetup.Locations.Contains(testgit.LocationUpstream) {
		self.PushBranchToRemote(branchSetup.Name, gitdomain.RemoteUpstream)
	}
	// step 3: set the branch type
	if branchType, hasBranchType := branchSetup.BranchType.Get(); hasBranchType {
		switch branchType {
		case configdomain.BranchTypeContributionBranch:
			self.MustRun("git-town", "contribute")
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeMainBranch:
			// nothing to do
		case configdomain.BranchTypePerennialBranch:
			self.MustRun("git-town", "set-parent", "--none")
		case configdomain.BranchTypeObservedBranch:
			self.MustRun("git-town", "observe")
		case configdomain.BranchTypeParkedBranch:
			self.MustRun("git-town", "park")
		case configdomain.BranchTypePrototypeBranch:
			self.MustRun("git-town", "prototype")
		}
	}
}

// CreateStandaloneTag creates a tag not on a branch.
func (self *TestCommands) CreateStandaloneTag(name string) {
	const tempBranchName = "temp"
	self.MustRun("git", "checkout", "-b", tempBranchName)
	filePath := filepath.Join(self.WorkingDir, "a.txt")
	file, err := os.Create(filePath)
	asserts.NoError(err)
	err = file.Close()
	asserts.NoError(err)
	self.MustRun("git", "add", "-A")
	self.MustRun("git", "commit", "-m", tempBranchName)
	self.MustRun("git", "tag", "-a", name, "-m", name)
	self.MustRun("git", "checkout", "-")
	self.MustRun("git", "branch", "-D", tempBranchName)
}

// CreateTag creates a tag with the given name.
func (self *TestCommands) CreateTag(name string) {
	self.MustRun("git", "tag", "-a", name, "-m", name)
}

// ExistingParent provides the first ancestor of the given branch that actually exists in the repo.
func (self *TestCommands) ExistingParent(branch gitdomain.LocalBranchName, lineage configdomain.Lineage) Option[gitdomain.BranchName] {
	for {
		parentOpt := lineage.Parent(branch)
		parent, hasParent := parentOpt.Get()
		if !hasParent {
			return None[gitdomain.BranchName]()
		}
		if self.Git.BranchExists(self, parent) {
			return Some(parent.BranchName())
		}
		if self.Git.BranchExistsRemotely(self, parent, gitdomain.RemoteOrigin) {
			return Some(parent.AtRemote(gitdomain.RemoteOrigin).BranchName())
		}
		branch = parent
	}
}

// Fetch retrieves the updates from the origin repo.
func (self *TestCommands) Fetch() {
	self.MustRun("git", "fetch", "--prune")
}

// FileContent provides the current content of a file.
func (self *TestCommands) FileContent(filename string) string {
	return asserts.NoError1(self.FileContentErr(filename))
}

// FileContentErr provides the current content of a file, returning an error if the file cannot be read.
func (self *TestCommands) FileContentErr(filename string) (string, error) {
	content, err := os.ReadFile(filepath.Join(self.WorkingDir, filename))
	return string(content), err
}

// FileContentInCommit provides the content of the file with the given name in the commit with the given SHA.
// If the file was deleted, the return content is empty and the bool return variable is set.
func (self *TestCommands) FileContentInCommit(location gitdomain.Location, filename string) (content string, deleted bool) {
	output, err := self.Query("git", "show", location.String()+":"+filename)
	if err != nil {
		return "", true
	}
	return output, false
}

// FilesInBranch provides the list of the files present in the given branch.
func (self *TestCommands) FilesInBranch(branch gitdomain.LocalBranchName) []string {
	output := self.MustQuery("git", "ls-tree", "-r", "--name-only", branch.String())
	return stringslice.NonEmptyLines(output)
}

// FilesInBranches provides a data table of files and their content in all branches.
func (self *TestCommands) FilesInBranches(mainBranch gitdomain.LocalBranchName) datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("BRANCH", "NAME", "CONTENT")
	branches, _ := asserts.NoError2(self.LocalBranchesMainFirst(mainBranch))
	var lastBranch gitdomain.LocalBranchName
	for _, branch := range branches {
		files := self.FilesInBranch(branch)
		for _, file := range files {
			content, deleted := self.FileContentInCommit(branch.Location(), file)
			if deleted {
				file += deletedText
			}
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
	return stringslice.NonEmptyLines(output)
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

func (self *TestCommands) Lineage() configdomain.Lineage {
	localSnapshot, _ := gitconfig.LoadSnapshot(self.TestRunner, Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedNo)
	localGitConfig, _ := config.NewPartialConfigFromSnapshot(localSnapshot, false, false, nil)
	return localGitConfig.Lineage
}

// LineageText formats the given lineage as a readable text representation.
func (self *TestCommands) LineageText(lineage configdomain.Lineage) string {
	return format.BranchLineage(lineage, configdomain.OrderAsc)
}

// LocalBranches provides the names of all branches in the local repository,
// ordered alphabetically.
func (self *TestCommands) LocalBranches() (allBranches, branchesInOtherWorktrees gitdomain.LocalBranchNames, err error) {
	forEachRefFormat := strings.Join(
		[]string{
			// worktree marker
			"%(if)", "%(HEAD)", "%(then)", // If the branch is checked out in the current worktree
			"H", // literal "H"
			"%(else)",
			"%(if)", "%(worktreepath)", "%(then)", // If the branch is checked out in any (other) worktree
			"W", // literal "W"
			"%(else)",
			"-", // literal "-"
			"%(end)",
			"%(end)",

			" ", // space separator

			"%(refname:lstrip=2)", // the branch name (without refs/heads/)
		},
		"")
	output, err := self.Query("git", "for-each-ref", "--format="+forEachRefFormat, "refs/heads/")
	if err != nil {
		return gitdomain.LocalBranchNames{}, gitdomain.LocalBranchNames{}, err
	}
	for _, line := range stringslice.Lines(output) {
		marker := line[0]
		branch := line[2:]
		allBranches = append(allBranches, gitdomain.NewLocalBranchName(branch))
		switch marker {
		case 'H', '-':
		case 'W':
			branchesInOtherWorktrees = append(branchesInOtherWorktrees, gitdomain.NewLocalBranchName(branch))
		default:
			panic(fmt.Sprintf("unexpected marker %q in line %q", marker, line))
		}
	}
	slice.NaturalSort(allBranches)
	slice.NaturalSort(branchesInOtherWorktrees)
	return allBranches, branchesInOtherWorktrees, nil
}

// LocalBranchesMainFirst provides the names of all local branches in this repo.
func (self *TestCommands) LocalBranchesMainFirst(mainBranch gitdomain.LocalBranchName) (allBranches, branchesInOtherWorktrees gitdomain.LocalBranchNames, err error) {
	allBranches, branchesInOtherWorktrees, err = self.LocalBranches()
	if err != nil {
		return
	}
	allBranches = slice.Hoist(allBranches, mainBranch)
	return
}

// LocalBranchesMainFirst provides the names of all local branches in this repo.
func (self *TestCommands) LocalBranchesOrderedHierarchically(lineage configdomain.Lineage, order configdomain.Order) (allBranches, branchesInOtherWorktrees gitdomain.LocalBranchNames, err error) {
	allBranches, branchesInOtherWorktrees, err = self.LocalBranches()
	if err != nil {
		return
	}
	allBranches = lineage.OrderHierarchically(allBranches, order)
	return
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

func (self *TestCommands) Reload() {
	self.SnapShots = self.Config.Reload(self.TestRunner)
}

// RemoveBranch deletes the branch with the given name from this repo.
func (self *TestCommands) RemoveBranch(name gitdomain.LocalBranchName) {
	self.MustRun("git", "branch", "-D", name.String())
}

// RemovePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (self *TestCommands) RemovePerennialBranchConfiguration() error {
	return gitconfig.RemovePerennialBranches(self.TestRunner)
}

// RemoveRemote deletes the Git remote with the given name.
func (self *TestCommands) RemoveRemote(name gitdomain.Remote) {
	self.Git.RemotesCache.Invalidate()
	self.MustRun("git", "remote", "rm", name.String())
}

// RemoveUnnecessaryFiles trims all files that aren't necessary in this repo.
func (self *TestCommands) RemoveUnnecessaryFiles() {
	fullPath := filepath.Join(self.WorkingDir, ".git", "hooks")
	asserts.NoError(os.RemoveAll(fullPath))
	_ = os.Remove(filepath.Join(self.WorkingDir, ".git", "COMMIT_EDITMSG"))
	_ = os.Remove(filepath.Join(self.WorkingDir, ".git", "description"))
}

func (self *TestCommands) RenameRemote(oldName, newName string) {
	self.Git.RemotesCache.Invalidate()
	self.MustRun("git", "remote", "rename", oldName, newName)
}

// SHAforBranch provides the SHA that the given branch points to.
func (self *TestCommands) SHAforBranch(branch gitdomain.LocalBranchName) gitdomain.SHA {
	output := self.MustQuery("git", "rev-parse", branch.String())
	return gitdomain.NewSHA(output)
}

// SHAForCommit provides the SHA for the commit with the given name.
func (self *TestCommands) SHAsForCommit(message gitdomain.CommitMessage) gitdomain.SHAs {
	output := self.MustQuery("git", "reflog", "--format=%H %s")
	if output == "" {
		panic(fmt.Sprintf("cannot find the SHA of commit %q", message))
	}
	shasWithMessage := make(gitdomain.SHAs, 0, 1)
	for text := range strings.SplitSeq(output, "\n") {
		shaText, commitMessage, found := strings.Cut(text, " ")
		if found && commitMessage == message.String() {
			sha := gitdomain.NewSHA(shaText)
			shasWithMessage = append(shasWithMessage, sha)
		}
	}
	return shasWithMessage
}

func (self *TestCommands) SetDefaultGitBranch(value gitdomain.LocalBranchName) {
	self.MustRun("git", "config", "init.defaultbranch", value.String())
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
	return stringslice.NonEmptyLines(output)
}

// UncommittedFiles provides the names of the files not committed into Git.
func (self *TestCommands) UncommittedFiles() []string {
	output := self.MustQuery("git", "status", "-z", "--untracked-files=all")
	statuses, err := prodgit.ParseGitStatusZ(output)
	asserts.NoError(err)
	result := make([]string, 0, len(statuses))
	for _, status := range statuses {
		result = append(result, status.Path)
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
	self.Config.Reload(self.TestRunner)
	for aliasName, aliasValue := range mapstools.SortedKeyValues(self.Config.NormalConfig.Aliases) {
		if strings.HasPrefix(aliasValue, "town ") {
			return fmt.Errorf("unexpected Git Town alias %q with value %q. All aliases: %#v", aliasName, aliasValue, self.Config.NormalConfig.Aliases)
		}
	}
	return nil
}
