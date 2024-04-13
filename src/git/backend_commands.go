package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/cache"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/messages"
)

type BackendRunner interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
	RunMany(commands [][]string) error
}

// BackendCommands are Git commands that Git Town executes to determine which frontend commands to run.
// They don't change the user's repo, execute instantaneously, and Git Town needs to know their output.
// They are invisible to the end user unless the "verbose" option is set.
type BackendCommands struct {
	Config             *config.Config                 // the known state of the Git repository
	CurrentBranchCache *cache.LocalBranchWithPrevious // caches the currently checked out Git branch
	DryRun             bool
	RemotesCache       *cache.Remotes // caches Git remotes
	Runner             BackendRunner  // executes shell commands in the directory of the Git repo
}

// Author provides the locally Git configured user.
func (self *BackendCommands) Author() (string, error) {
	email := self.Config.FullConfig.GitUserEmail
	if email == "" {
		return "", errors.New(messages.GitUserEmailMissing)
	}
	name := self.Config.FullConfig.GitUserName
	if name == "" {
		return "", errors.New(messages.GitUserEmailMissing)
	}
	return name + " <" + email + ">", nil
}

// BranchAuthors provides the user accounts that contributed to the given branch.
// Returns lines of "name <email>".
func (self *BackendCommands) BranchAuthors(branch, parent gitdomain.LocalBranchName) ([]string, error) {
	output, err := self.Runner.QueryTrim("git", "shortlog", "-s", "-n", "-e", parent.String()+".."+branch.String())
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, line := range stringslice.Lines(output) {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result = append(result, parts[1])
	}
	return result, nil
}

func (self *BackendCommands) BranchExists(branch gitdomain.LocalBranchName) bool {
	err := self.Runner.Run("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch.String())
	return err == nil
}

// BranchHasUnmergedChanges indicates whether the branch with the given name
// contains changes that were not merged into the main branch.
func (self *BackendCommands) BranchHasUnmergedChanges(branch, parent gitdomain.LocalBranchName) (bool, error) {
	out, err := self.Runner.QueryTrim("git", "diff", parent.String()+".."+branch.String())
	if err != nil {
		return false, fmt.Errorf(messages.BranchDiffProblem, branch, err)
	}
	return out != "", nil
}

// BranchesSnapshot provides detailed information about the sync status of all branches.
func (self *BackendCommands) BranchesSnapshot() (gitdomain.BranchesSnapshot, error) { //nolint:nonamedreturns
	output, err := self.Runner.Query("git", "branch", "-vva", "--sort=refname")
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), err
	}
	branches, currentBranch := ParseVerboseBranchesOutput(output)
	if !currentBranch.IsEmpty() {
		self.CurrentBranchCache.Set(currentBranch)
	}
	return gitdomain.BranchesSnapshot{
		Branches: branches,
		Active:   currentBranch,
	}, nil
}

// CheckoutBranch checks out the Git branch with the given name.
func (self *BackendCommands) CheckoutBranch(name gitdomain.LocalBranchName) error {
	if !self.DryRun {
		err := self.CheckoutBranchUncached(name)
		if err != nil {
			return err
		}
	}
	if name.String() != "-" {
		self.CurrentBranchCache.Set(name)
	} else {
		self.CurrentBranchCache.Invalidate()
	}
	return nil
}

func IsAhead(branchName, remoteText string) (bool, gitdomain.RemoteBranchName) {
	reText := fmt.Sprintf(`\[(\w+\/%s): ahead \d+\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, gitdomain.NewRemoteBranchName(matches[1])
	}
	return false, gitdomain.EmptyRemoteBranchName()
}

func IsAheadAndBehind(branchName, remoteText string) (bool, gitdomain.RemoteBranchName) {
	reText := fmt.Sprintf(`\[(\w+\/%s): ahead \d+, behind \d+\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, gitdomain.NewRemoteBranchName(matches[1])
	}
	return false, gitdomain.EmptyRemoteBranchName()
}

func IsBehind(branchName, remoteText string) (bool, gitdomain.RemoteBranchName) {
	reText := fmt.Sprintf(`\[(\w+\/%s): behind \d+\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, gitdomain.NewRemoteBranchName(matches[1])
	}
	return false, gitdomain.EmptyRemoteBranchName()
}

func IsInSync(branchName, remoteText string) (bool, gitdomain.RemoteBranchName) {
	reText := fmt.Sprintf(`\[(\w+\/%s)\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, gitdomain.NewRemoteBranchName(matches[1])
	}
	return false, gitdomain.EmptyRemoteBranchName()
}

// IsRemoteGone indicates whether the given remoteText indicates a deleted tracking branch.
func IsRemoteGone(branchName, remoteText string) (bool, gitdomain.RemoteBranchName) {
	reText := fmt.Sprintf(`^\[(\w+\/%s): gone\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, gitdomain.NewRemoteBranchName(matches[1])
	}
	return false, gitdomain.EmptyRemoteBranchName()
}

// CheckoutBranch checks out the Git branch with the given name.
func (self *BackendCommands) CheckoutBranchUncached(name gitdomain.LocalBranchName) error {
	err := self.Runner.Run("git", "checkout", name.String())
	if err != nil {
		return fmt.Errorf(messages.BranchCheckoutProblem, name, err)
	}
	return nil
}

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// Adds the given prefix with the newline if provided.
func (self *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {
	squashMessageFile := ".git/SQUASH_MSG"
	contentBytes, err := os.ReadFile(squashMessageFile)
	if err != nil {
		return fmt.Errorf(messages.SquashCannotReadFile, squashMessageFile, err)
	}
	content := string(contentBytes)
	if prefix != "" {
		content = prefix + "\n" + content
	}
	content = regexp.MustCompile("(?m)^").ReplaceAllString(content, "# ")
	return os.WriteFile(squashMessageFile, []byte(content), 0o600)
}

func (self *BackendCommands) CommitsInBranch(branch, parent gitdomain.LocalBranchName) (gitdomain.Commits, error) {
	if parent.IsEmpty() {
		return self.CommitsInPerennialBranch()
	}
	return self.CommitsInFeatureBranch(branch, parent)
}

func (self *BackendCommands) CommitsInFeatureBranch(branch, parent gitdomain.LocalBranchName) (gitdomain.Commits, error) {
	output, err := self.Runner.QueryTrim("git", "cherry", "-v", parent.String(), branch.String())
	if err != nil {
		return gitdomain.Commits{}, err
	}
	lines := strings.Split(output, "\n")
	result := make(gitdomain.Commits, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		sha, message, ok := strings.Cut(line[2:], " ")
		if !ok {
			continue
		}
		result = append(result, gitdomain.Commit{
			Message: gitdomain.CommitMessage(message),
			SHA:     gitdomain.NewSHA(sha),
		})
	}
	return result, nil
}

func (self *BackendCommands) CommitsInPerennialBranch() (gitdomain.Commits, error) {
	output, err := self.Runner.QueryTrim("git", "log", "--pretty=format:%h %s", "-10")
	if err != nil {
		return gitdomain.Commits{}, err
	}
	lines := stringslice.Lines(output)
	result := make(gitdomain.Commits, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		sha, message, ok := strings.Cut(line, " ")
		if !ok {
			continue
		}
		result = append(result, gitdomain.Commit{
			Message: gitdomain.CommitMessage(message),
			SHA:     gitdomain.NewSHA(sha),
		})
	}
	return result, nil
}

// CurrentBranch provides the name of the currently checked out branch.
func (self *BackendCommands) CurrentBranch() (gitdomain.LocalBranchName, error) {
	if !self.CurrentBranchCache.Initialized() {
		currentBranch, err := self.CurrentBranchUncached()
		if err != nil {
			return currentBranch, err
		}
		self.CurrentBranchCache.Set(currentBranch)
	}
	return self.CurrentBranchCache.Value(), nil
}

// CurrentBranch provides the currently checked out branch.
func (self *BackendCommands) CurrentBranchUncached() (gitdomain.LocalBranchName, error) {
	repoStatus, err := self.RepoStatus()
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), fmt.Errorf(messages.BranchCurrentProblem, err)
	}
	if repoStatus.RebaseInProgress {
		currentBranch, err := self.currentBranchDuringRebase()
		if err != nil {
			return gitdomain.EmptyLocalBranchName(), err
		}
		return currentBranch, nil
	}
	output, err := self.Runner.QueryTrim("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), fmt.Errorf(messages.BranchCurrentProblem, err)
	}
	return gitdomain.NewLocalBranchName(output), nil
}

// CurrentSHA provides the SHA of the currently checked out branch/commit.
func (self *BackendCommands) CurrentSHA() (gitdomain.SHA, error) {
	return self.SHAForBranch(gitdomain.NewBranchName("HEAD"))
}

func (self *BackendCommands) DefaultBranch() gitdomain.LocalBranchName {
	name, _ := self.Runner.QueryTrim("git", "config", "--get", "init.defaultbranch")
	return gitdomain.LocalBranchName(name)
}

func (self *BackendCommands) FirstExistingBranch(branches gitdomain.LocalBranchNames, mainBranch gitdomain.LocalBranchName) gitdomain.LocalBranchName {
	for _, branch := range branches {
		if self.BranchExists(branch) {
			return branch
		}
	}
	return mainBranch
}

// HasLocalBranch indicates whether this repo has a local branch with the given name.
func (self *BackendCommands) HasLocalBranch(name gitdomain.LocalBranchName) bool {
	return self.Runner.Run("git", "show-ref", "--quiet", "refs/heads/"+name.String()) == nil
}

// HasMergeInProgress indicates whether this Git repository currently has a merge in progress.
func (self *BackendCommands) HasMergeInProgress() bool {
	err := self.Runner.Run("git", "rev-parse", "-q", "--verify", "MERGE_HEAD")
	return err == nil
}

// HasShippableChanges indicates whether the given branch has changes
// not currently in the main branch.
func (self *BackendCommands) HasShippableChanges(branch, mainBranch gitdomain.LocalBranchName) (bool, error) {
	out, err := self.Runner.QueryTrim("git", "diff", mainBranch.String()+".."+branch.String())
	if err != nil {
		return false, fmt.Errorf(messages.ShippableChangesProblem, branch, err)
	}
	return out != "", nil
}

// LastCommitMessage provides the commit message for the last commit.
func (self *BackendCommands) LastCommitMessage() (gitdomain.CommitMessage, error) {
	out, err := self.Runner.QueryTrim("git", "log", "-1", "--format=%B")
	if err != nil {
		return "", fmt.Errorf(messages.CommitMessageProblem, err)
	}
	return gitdomain.CommitMessage(out), nil
}

// PreviouslyCheckedOutBranch provides the name of the branch that was previously checked out in this repo.
func (self *BackendCommands) PreviouslyCheckedOutBranch() gitdomain.LocalBranchName {
	output, err := self.Runner.QueryTrim("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return gitdomain.EmptyLocalBranchName()
	}
	if output == "" {
		return gitdomain.EmptyLocalBranchName()
	}
	return gitdomain.NewLocalBranchName(output)
}

// Remotes provides the names of all Git remotes in this repository.
func (self *BackendCommands) Remotes() (gitdomain.Remotes, error) {
	if !self.RemotesCache.Initialized() {
		remotes, err := self.RemotesUncached()
		if err != nil {
			return remotes, err
		}
		self.RemotesCache.Set(remotes)
	}
	return self.RemotesCache.Value(), nil
}

// Remotes provides the names of all Git remotes in this repository.
func (self *BackendCommands) RemotesUncached() (gitdomain.Remotes, error) {
	out, err := self.Runner.QueryTrim("git", "remote")
	if err != nil {
		return gitdomain.Remotes{}, fmt.Errorf(messages.RemotesProblem, err)
	}
	if out == "" {
		return gitdomain.Remotes{}, nil
	}
	return gitdomain.NewRemotes(stringslice.Lines(out)...), nil
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (self *BackendCommands) RemoveOutdatedConfiguration(localBranches gitdomain.LocalBranchNames) error {
	for child, parent := range self.Config.FullConfig.Lineage {
		hasChildBranch := localBranches.Contains(child)
		hasParentBranch := localBranches.Contains(parent)
		if !hasChildBranch || !hasParentBranch {
			self.Config.RemoveParent(child)
		}
	}
	return nil
}

// RepoStatus provides a summary of the state the current workspace is in right now: rebasing, has conflicts, has open changes, etc.
func (self *BackendCommands) RepoStatus() (gitdomain.RepoStatus, error) {
	output, err := self.Runner.QueryTrim("git", "status", "--long", "--ignore-submodules")
	if err != nil {
		return gitdomain.RepoStatus{}, fmt.Errorf(messages.ConflictDetectionProblem, err)
	}
	hasConflicts := strings.Contains(output, "Unmerged paths")
	hasOpenChanges := outputIndicatesOpenChanges(output)
	hasUntrackedChanges := outputIndicatesUntrackedChanges(output)
	rebaseInProgress := outputIndicatesRebaseInProgress(output)
	return gitdomain.RepoStatus{
		Conflicts:        hasConflicts,
		OpenChanges:      hasOpenChanges,
		RebaseInProgress: rebaseInProgress,
		UntrackedChanges: hasUntrackedChanges,
	}, nil
}

// RootDirectory provides the path of the root directory of the current repository,
// i.e. the directory that contains the ".git" folder.
func (self *BackendCommands) RootDirectory() gitdomain.RepoRootDir {
	output, err := self.Runner.QueryTrim("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return gitdomain.EmptyRepoRootDir()
	}
	return gitdomain.NewRepoRootDir(filepath.FromSlash(output))
}

// SHAForBranch provides the SHA for the local branch with the given name.
func (self *BackendCommands) SHAForBranch(name gitdomain.BranchName) (gitdomain.SHA, error) {
	output, err := self.Runner.QueryTrim("git", "rev-parse", "--short", name.String())
	if err != nil {
		return gitdomain.EmptySHA(), fmt.Errorf(messages.BranchLocalSHAProblem, name, err)
	}
	return gitdomain.NewSHA(output), nil
}

// ShouldPushBranch returns whether the local branch with the given name
// contains commits that have not been pushed to its tracking branch.
func (self *BackendCommands) ShouldPushBranch(branch gitdomain.LocalBranchName, trackingBranch gitdomain.RemoteBranchName) (bool, error) {
	out, err := self.Runner.QueryTrim("git", "rev-list", "--left-right", branch.String()+"..."+trackingBranch.String())
	if err != nil {
		return false, fmt.Errorf(messages.DiffProblem, branch, branch, err)
	}
	return out != "", nil
}

// StashSize provides the number of stashes in this repository.
func (self *BackendCommands) StashSize() (gitdomain.StashSize, error) {
	output, err := self.Runner.QueryTrim("git", "stash", "list")
	return gitdomain.StashSize(len(stringslice.Lines(output))), err
}

// Version indicates whether the needed Git version is installed.
func (self *BackendCommands) Version() (major int, minor int, err error) {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	output, err := self.Runner.QueryTrim("git", "version")
	if err != nil {
		return 0, 0, fmt.Errorf(messages.GitVersionProblem, err)
	}
	matches := versionRegexp.FindStringSubmatch(output)
	if matches == nil {
		return 0, 0, fmt.Errorf(messages.GitVersionUnexpectedOutput, output)
	}
	majorVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf(messages.GitVersionMajorNotNumber, matches[1], err)
	}
	minorVersion, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf(messages.GitVersionMinorNotNumber, matches[2], err)
	}
	return majorVersion, minorVersion, nil
}

func (self *BackendCommands) currentBranchDuringRebase() (gitdomain.LocalBranchName, error) {
	output, err := self.Runner.QueryTrim("git", "branch", "--list")
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), err
	}
	lines := stringslice.Lines(output)
	linesWithStar := stringslice.LinesWithPrefix(lines, "* ")
	if len(linesWithStar) == 0 {
		return gitdomain.EmptyLocalBranchName(), err
	}
	if len(linesWithStar) > 1 {
		panic("multiple lines with star found:\n " + output)
	}
	lineWithStar := linesWithStar[0]
	return ParseActiveBranchDuringRebase(lineWithStar), nil
}

func ParseActiveBranchDuringRebase(lineWithStar string) gitdomain.LocalBranchName {
	parts := strings.Split(lineWithStar, " ")
	partsWithBranchName := parts[4:]
	branchNameWithClosingParen := strings.Join(partsWithBranchName, " ")
	return gitdomain.NewLocalBranchName(branchNameWithClosingParen[:len(branchNameWithClosingParen)-1])
}

// ParseVerboseBranchesOutput provides the branches in the given Git output as well as the name of the currently checked out branch.
func ParseVerboseBranchesOutput(output string) (gitdomain.BranchInfos, gitdomain.LocalBranchName) {
	result := gitdomain.BranchInfos{}
	spaceRE := regexp.MustCompile(" +")
	lines := stringslice.Lines(output)
	checkedoutBranch := gitdomain.EmptyLocalBranchName()
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := spaceRE.Split(line[2:], 3)
		if strings.HasSuffix(parts[0], "/HEAD") {
			continue
		}
		if len(parts) < 2 {
			// This shouldn't happen, but did happen in https://github.com/git-town/git-town/issues/2562.
			fmt.Printf(messages.GitOutputIrregular, line, output)
			os.Exit(1)
		}
		branchName := parts[0]
		var sha gitdomain.SHA
		if parts[1] == "branch," {
			// we are rebasing and don't need the SHA
			sha = gitdomain.EmptySHA()
		} else {
			sha = gitdomain.NewSHA(parts[1])
		}
		remoteText := parts[2]
		if line[0] == '*' && branchName != "(no" { // "(no" as in "(no branch, rebasing main)" is what we get when a rebase is active, in which case no branch is checked out
			checkedoutBranch = gitdomain.NewLocalBranchName(branchName)
		}
		syncStatus, trackingBranchName := determineSyncStatus(branchName, remoteText)
		switch {
		case line[0] == '+':
			result = append(result, gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName(branchName),
				LocalSHA:   sha,
				SyncStatus: gitdomain.SyncStatusOtherWorktree,
				RemoteName: trackingBranchName,
				RemoteSHA:  gitdomain.EmptySHA(),
			})
		case isLocalBranchName(branchName):
			result = append(result, gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName(branchName),
				LocalSHA:   sha,
				SyncStatus: syncStatus,
				RemoteName: trackingBranchName,
				RemoteSHA:  gitdomain.EmptySHA(), // will be added later
			})
		default:
			remoteBranchName := gitdomain.NewRemoteBranchName(strings.TrimPrefix(branchName, "remotes/"))
			existingBranchWithTracking := result.FindByRemoteName(remoteBranchName)
			if existingBranchWithTracking != nil {
				existingBranchWithTracking.RemoteSHA = sha
			} else {
				result = append(result, gitdomain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: remoteBranchName,
					RemoteSHA:  sha,
				})
			}
		}
	}
	return result, checkedoutBranch
}

func determineSyncStatus(branchName, remoteText string) (syncStatus gitdomain.SyncStatus, trackingBranchName gitdomain.RemoteBranchName) {
	isInSync, trackingBranchName := IsInSync(branchName, remoteText)
	if isInSync {
		return gitdomain.SyncStatusUpToDate, trackingBranchName
	}
	isGone, trackingBranchName := IsRemoteGone(branchName, remoteText)
	if isGone {
		return gitdomain.SyncStatusDeletedAtRemote, trackingBranchName
	}
	IsAhead, trackingBranchName := IsAhead(branchName, remoteText)
	if IsAhead {
		return gitdomain.SyncStatusNotInSync, trackingBranchName
	}
	IsBehind, trackingBranchName := IsBehind(branchName, remoteText)
	if IsBehind {
		return gitdomain.SyncStatusNotInSync, trackingBranchName
	}
	IsAheadAndBehind, trackingBranchName := IsAheadAndBehind(branchName, remoteText)
	if IsAheadAndBehind {
		return gitdomain.SyncStatusNotInSync, trackingBranchName
	}
	if strings.HasPrefix(branchName, "remotes/") {
		return gitdomain.SyncStatusRemoteOnly, gitdomain.EmptyRemoteBranchName()
	}
	return gitdomain.SyncStatusLocalOnly, gitdomain.EmptyRemoteBranchName()
}

// isLocalBranchName indicates whether the branch with the given Git ref is local or remote.
func isLocalBranchName(branch string) bool {
	return !strings.HasPrefix(branch, "remotes/")
}

func outputIndicatesMergeInProgress(output string) bool {
	if strings.Contains(output, "You have unmerged paths") {
		return true
	}
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "AA ") {
			return true
		}
	}
	return false
}

// HasOpenChanges indicates whether this repo has open changes.
func outputIndicatesOpenChanges(output string) bool {
	if strings.Contains(output, "working tree clean") || strings.Contains(output, "nothing to commit") {
		return false
	}
	if outputIndicatesRebaseInProgress(output) || outputIndicatesMergeInProgress(output) {
		return false
	}
	return true
}

func outputIndicatesRebaseInProgress(output string) bool {
	return strings.Contains(output, "rebase in progress") || strings.Contains(output, "You are currently rebasing")
}

func outputIndicatesUntrackedChanges(output string) bool {
	return strings.Contains(output, "Untracked files:")
}
