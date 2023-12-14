package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v11/src/config"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/gohacks/cache"
	"github.com/git-town/git-town/v11/src/gohacks/stringslice"
	"github.com/git-town/git-town/v11/src/messages"
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
	BackendRunner                                     // executes shell commands in the directory of the Git repo
	*config.GitTown                                   // the known state of the Git repository
	CurrentBranchCache *cache.LocalBranchWithPrevious // caches the currently checked out Git branch
	RemotesCache       *cache.Remotes                 // caches Git remotes
}

// Author provides the locally Git configured user.
func (self *BackendCommands) Author() (string, error) {
	// TODO: read this from the config cache?
	// If not possible, comment here why.
	name, err := self.QueryTrim("git", "config", "user.name")
	if err != nil {
		return "", err
	}
	email, err := self.QueryTrim("git", "config", "user.email")
	if err != nil {
		return "", err
	}
	return name + " <" + email + ">", nil
}

// BranchAuthors provides the user accounts that contributed to the given branch.
// Returns lines of "name <email>".
func (self *BackendCommands) BranchAuthors(branch, parent domain.LocalBranchName) ([]string, error) {
	output, err := self.QueryTrim("git", "shortlog", "-s", "-n", "-e", parent.String()+".."+branch.String())
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

func (self *BackendCommands) BranchExists(branch domain.LocalBranchName) bool {
	err := self.Run("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch.String())
	return err == nil
}

// BranchHasUnmergedChanges indicates whether the branch with the given name
// contains changes that were not merged into the main branch.
func (self *BackendCommands) BranchHasUnmergedChanges(branch, parent domain.LocalBranchName) (bool, error) {
	out, err := self.QueryTrim("git", "diff", parent.String()+".."+branch.String())
	if err != nil {
		return false, fmt.Errorf(messages.BranchDiffProblem, branch, err)
	}
	return out != "", nil
}

// BranchHasUnmergedCommits indicates whether the branch with the given name
// contains commits that are not merged into the main branch.
func (self *BackendCommands) BranchHasUnmergedCommits(branch domain.LocalBranchName, parent domain.Location) (bool, error) {
	out, err := self.QueryTrim("git", "log", parent.String()+".."+branch.String())
	if err != nil {
		return false, fmt.Errorf(messages.BranchDiffProblem, branch, err)
	}
	return out != "", nil
}

// BranchesSnapshot provides detailed information about the sync status of all branches.
func (self *BackendCommands) BranchesSnapshot() (domain.BranchesSnapshot, error) { //nolint:nonamedreturns
	output, err := self.Query("git", "branch", "-vva")
	if err != nil {
		return domain.EmptyBranchesSnapshot(), err
	}
	branches, currentBranch := ParseVerboseBranchesOutput(output)
	if !currentBranch.IsEmpty() {
		self.CurrentBranchCache.Set(currentBranch)
	}
	return domain.BranchesSnapshot{
		Branches: branches,
		Active:   currentBranch,
	}, nil
}

// CheckoutBranch checks out the Git branch with the given name.
func (self *BackendCommands) CheckoutBranch(name domain.LocalBranchName) error {
	if !self.GitTown.DryRun {
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

func IsAhead(branchName, remoteText string) (bool, domain.RemoteBranchName) {
	reText := fmt.Sprintf(`\[(\w+\/%s): ahead \d+\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, domain.NewRemoteBranchName(matches[1])
	}
	return false, domain.EmptyRemoteBranchName()
}

func IsAheadAndBehind(branchName, remoteText string) (bool, domain.RemoteBranchName) {
	reText := fmt.Sprintf(`\[(\w+\/%s): ahead \d+, behind \d+\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, domain.NewRemoteBranchName(matches[1])
	}
	return false, domain.EmptyRemoteBranchName()
}

func IsBehind(branchName, remoteText string) (bool, domain.RemoteBranchName) {
	reText := fmt.Sprintf(`\[(\w+\/%s): behind \d+\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, domain.NewRemoteBranchName(matches[1])
	}
	return false, domain.EmptyRemoteBranchName()
}

func IsInSync(branchName, remoteText string) (bool, domain.RemoteBranchName) {
	reText := fmt.Sprintf(`\[(\w+\/%s)\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, domain.NewRemoteBranchName(matches[1])
	}
	return false, domain.EmptyRemoteBranchName()
}

// IsRemoteGone indicates whether the given remoteText indicates a deleted tracking branch.
func IsRemoteGone(branchName, remoteText string) (bool, domain.RemoteBranchName) {
	reText := fmt.Sprintf(`^\[(\w+\/%s): gone\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, domain.NewRemoteBranchName(matches[1])
	}
	return false, domain.EmptyRemoteBranchName()
}

// ParseVerboseBranchesOutput provides the branches in the given Git output as well as the name of the currently checked out branch.
func ParseVerboseBranchesOutput(output string) (domain.BranchInfos, domain.LocalBranchName) {
	result := domain.BranchInfos{}
	spaceRE := regexp.MustCompile(" +")
	lines := stringslice.Lines(output)
	checkedoutBranch := domain.EmptyLocalBranchName()
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := spaceRE.Split(line[2:], 3)
		if parts[0] == "remotes/origin/HEAD" {
			continue
		}
		if len(parts) < 2 {
			// This shouldn't happen, but did happen in https://github.com/git-town/git-town/issues/2562.
			fmt.Println("ERROR: Encountered irregular Git output")
			fmt.Println()
			fmt.Println("PLEASE REPORT THE OUTPUT BELOW AT https://github.com/git-town/git-town/issues/new")
			fmt.Println()
			fmt.Printf("Problematic line: %q\n", line)
			fmt.Println()
			fmt.Println("BEGIN OUTPUT FROM 'git branch -vva'")
			fmt.Println(output)
			fmt.Println("END OUTPUT FROM 'git branch -vva'")
			os.Exit(1)
		}
		branchName := parts[0]
		var sha domain.SHA
		if parts[1] == "branch," {
			// we are rebasing and don't need the SHA
			sha = domain.EmptySHA()
		} else {
			sha = domain.NewSHA(parts[1])
		}
		remoteText := parts[2]
		if line[0] == '*' && branchName != "(no" { // "(no" as in "(no branch, rebasing main)" is what we get when a rebase is active, in which case no branch is checked out
			checkedoutBranch = domain.NewLocalBranchName(branchName)
		}
		syncStatus, trackingBranchName := determineSyncStatus(branchName, remoteText)
		switch {
		case line[0] == '+':
			result = append(result, domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName(branchName),
				LocalSHA:   sha,
				SyncStatus: domain.SyncStatusOtherWorktree,
				RemoteName: trackingBranchName,
				RemoteSHA:  domain.EmptySHA(),
			})
		case isLocalBranchName(branchName):
			result = append(result, domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName(branchName),
				LocalSHA:   sha,
				SyncStatus: syncStatus,
				RemoteName: trackingBranchName,
				RemoteSHA:  domain.EmptySHA(), // will be added later
			})
		default:
			remoteBranchName := domain.NewRemoteBranchName(strings.TrimPrefix(branchName, "remotes/"))
			existingBranchWithTracking := result.FindByRemoteName(remoteBranchName)
			if existingBranchWithTracking != nil {
				existingBranchWithTracking.RemoteSHA = sha
			} else {
				result = append(result, domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: remoteBranchName,
					RemoteSHA:  sha,
				})
			}
		}
	}
	return result, checkedoutBranch
}

func determineSyncStatus(branchName, remoteText string) (syncStatus domain.SyncStatus, trackingBranchName domain.RemoteBranchName) {
	isInSync, trackingBranchName := IsInSync(branchName, remoteText)
	if isInSync {
		return domain.SyncStatusUpToDate, trackingBranchName
	}
	isGone, trackingBranchName := IsRemoteGone(branchName, remoteText)
	if isGone {
		return domain.SyncStatusDeletedAtRemote, trackingBranchName
	}
	IsAhead, trackingBranchName := IsAhead(branchName, remoteText)
	if IsAhead {
		return domain.SyncStatusNotInSync, trackingBranchName
	}
	IsBehind, trackingBranchName := IsBehind(branchName, remoteText)
	if IsBehind {
		return domain.SyncStatusNotInSync, trackingBranchName
	}
	IsAheadAndBehind, trackingBranchName := IsAheadAndBehind(branchName, remoteText)
	if IsAheadAndBehind {
		return domain.SyncStatusNotInSync, trackingBranchName
	}
	if strings.HasPrefix(branchName, "remotes/") {
		return domain.SyncStatusRemoteOnly, domain.EmptyRemoteBranchName()
	}
	return domain.SyncStatusLocalOnly, domain.EmptyRemoteBranchName()
}

// isLocalBranchName indicates whether the branch with the given Git ref is local or remote.
func isLocalBranchName(branch string) bool {
	return !strings.HasPrefix(branch, "remotes/")
}

// CheckoutBranch checks out the Git branch with the given name.
func (self *BackendCommands) CheckoutBranchUncached(name domain.LocalBranchName) error {
	err := self.Run("git", "checkout", name.String())
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

func (self *BackendCommands) CommitsInBranch(branch, parent domain.LocalBranchName) (domain.SHAs, error) {
	if parent.IsEmpty() {
		return self.CommitsInPerennialBranch()
	}
	return self.CommitsInFeatureBranch(branch, parent)
}

func (self *BackendCommands) CommitsInFeatureBranch(branch, parent domain.LocalBranchName) (domain.SHAs, error) {
	output, err := self.QueryTrim("git", "cherry", parent.String(), branch.String())
	if err != nil {
		return domain.SHAs{}, err
	}
	lines := strings.Split(output, "\n")
	result := make([]domain.SHA, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 {
			result = append(result, domain.NewSHA(line[2:9]))
		}
	}
	return result, nil
}

func (self *BackendCommands) CommitsInPerennialBranch() (domain.SHAs, error) {
	output, err := self.QueryTrim("git", "log", "--pretty=format:%h", "-10")
	if err != nil {
		return domain.SHAs{}, err
	}
	lines := strings.Split(output, "\n")
	result := make([]domain.SHA, 0, len(lines))
	for _, line := range lines {
		result = append(result, domain.NewSHA(line))
	}
	return result, nil
}

// CreateFeatureBranch creates a feature branch with the given name in this repository.
func (self *BackendCommands) CreateFeatureBranch(name domain.LocalBranchName) error {
	err := self.RunMany([][]string{
		{"git", "branch", name.String(), "main"},
		{"git", "config", "git-town-branch." + name.String() + ".parent", "main"},
	})
	if err != nil {
		return fmt.Errorf(messages.BranchFeatureCannotCreate, name, err)
	}
	return nil
}

// CurrentBranch provides the name of the currently checked out branch.
func (self *BackendCommands) CurrentBranch() (domain.LocalBranchName, error) {
	if self.GitTown.DryRun {
		return self.CurrentBranchCache.Value(), nil
	}
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
func (self *BackendCommands) CurrentBranchUncached() (domain.LocalBranchName, error) {
	repoStatus, err := self.RepoStatus()
	if err != nil {
		return domain.EmptyLocalBranchName(), fmt.Errorf(messages.BranchCurrentProblem, err)
	}
	if repoStatus.RebaseInProgress {
		currentBranch, err := self.currentBranchDuringRebase()
		if err != nil {
			return domain.EmptyLocalBranchName(), err
		}
		return currentBranch, nil
	}
	output, err := self.QueryTrim("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return domain.EmptyLocalBranchName(), fmt.Errorf(messages.BranchCurrentProblem, err)
	}
	return domain.NewLocalBranchName(output), nil
}

// CurrentSHA provides the SHA of the currently checked out branch/commit.
func (self *BackendCommands) CurrentSHA() (domain.SHA, error) {
	return self.SHAForBranch(domain.NewBranchName("HEAD"))
}

func (self *BackendCommands) FirstExistingBranch(branches domain.LocalBranchNames, mainBranch domain.LocalBranchName) domain.LocalBranchName {
	for _, branch := range branches {
		if self.BranchExists(branch) {
			return branch
		}
	}
	return mainBranch
}

// HasLocalBranch indicates whether this repo has a local branch with the given name.
func (self *BackendCommands) HasLocalBranch(name domain.LocalBranchName) bool {
	return self.Run("git", "show-ref", "--quiet", "refs/heads/"+name.String()) == nil
}

// HasMergeInProgress indicates whether this Git repository currently has a merge in progress.
func (self *BackendCommands) HasMergeInProgress() bool {
	err := self.Run("git", "rev-parse", "-q", "--verify", "MERGE_HEAD")
	return err == nil
}

// HasShippableChanges indicates whether the given branch has changes
// not currently in the main branch.
func (self *BackendCommands) HasShippableChanges(branch, mainBranch domain.LocalBranchName) (bool, error) {
	out, err := self.QueryTrim("git", "diff", mainBranch.String()+".."+branch.String())
	if err != nil {
		return false, fmt.Errorf(messages.ShippableChangesProblem, branch, err)
	}
	return out != "", nil
}

// LastCommitMessage provides the commit message for the last commit.
func (self *BackendCommands) LastCommitMessage() (string, error) {
	out, err := self.QueryTrim("git", "log", "-1", "--format=%B")
	if err != nil {
		return "", fmt.Errorf(messages.CommitMessageProblem, err)
	}
	return out, nil
}

// PreviouslyCheckedOutBranch provides the name of the branch that was previously checked out in this repo.
func (self *BackendCommands) PreviouslyCheckedOutBranch() domain.LocalBranchName {
	output, err := self.QueryTrim("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return domain.EmptyLocalBranchName()
	}
	if output == "" {
		return domain.EmptyLocalBranchName()
	}
	return domain.NewLocalBranchName(output)
}

// Remotes provides the names of all Git remotes in this repository.
func (self *BackendCommands) Remotes() (domain.Remotes, error) {
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
func (self *BackendCommands) RemotesUncached() (domain.Remotes, error) {
	out, err := self.QueryTrim("git", "remote")
	if err != nil {
		return domain.Remotes{}, fmt.Errorf(messages.RemotesProblem, err)
	}
	if out == "" {
		return domain.Remotes{}, nil
	}
	return domain.NewRemotes(stringslice.Lines(out)...), nil
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (self *BackendCommands) RemoveOutdatedConfiguration(allBranches domain.BranchInfos) error {
	for child, parent := range self.GitTown.Lineage(self.GitTown.RemoveLocalConfigValue) {
		hasChildBranch := allBranches.HasLocalBranch(child)
		hasParentBranch := allBranches.HasLocalBranch(parent)
		if !hasChildBranch || !hasParentBranch {
			self.GitTown.RemoveParent(child)
		}
	}
	return nil
}

// HasConflicts returns whether the local repository currently has unresolved merge conflicts.
func (self *BackendCommands) RepoStatus() (domain.RepoStatus, error) {
	output, err := self.QueryTrim("git", "status", "--long", "--ignore-submodules")
	if err != nil {
		return domain.RepoStatus{}, fmt.Errorf(messages.ConflictDetectionProblem, err)
	}
	hasConflicts := strings.Contains(output, "Unmerged paths")
	hasOpenChanges := outputIndicatesOpenChanges(output)
	hasUntrackedChanges := outputIndicatesUntrackedChanges(output)
	rebaseInProgress := outputIndicatesRebaseInProgress(output)
	return domain.RepoStatus{
		Conflicts:        hasConflicts,
		OpenChanges:      hasOpenChanges,
		RebaseInProgress: rebaseInProgress,
		UntrackedChanges: hasUntrackedChanges,
	}, nil
}

// RootDirectory provides the path of the rood directory of the current repository,
// i.e. the directory that contains the ".git" folder.
func (self *BackendCommands) RootDirectory() domain.RepoRootDir {
	output, err := self.QueryTrim("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return domain.EmptyRepoRootDir()
	}
	return domain.NewRepoRootDir(filepath.FromSlash(output))
}

// SHAForBranch provides the SHA for the local branch with the given name.
func (self *BackendCommands) SHAForBranch(name domain.BranchName) (domain.SHA, error) {
	output, err := self.QueryTrim("git", "rev-parse", "--short", name.String())
	if err != nil {
		return domain.EmptySHA(), fmt.Errorf(messages.BranchLocalSHAProblem, name, err)
	}
	return domain.NewSHA(output), nil
}

// ShouldPushBranch returns whether the local branch with the given name
// contains commits that have not been pushed to its tracking branch.
func (self *BackendCommands) ShouldPushBranch(branch domain.LocalBranchName, trackingBranch domain.RemoteBranchName) (bool, error) {
	out, err := self.QueryTrim("git", "rev-list", "--left-right", branch.String()+"..."+trackingBranch.String())
	if err != nil {
		return false, fmt.Errorf(messages.DiffProblem, branch, branch, err)
	}
	return out != "", nil
}

// StashSnapshot provides the number of stashes in this repository.
func (self *BackendCommands) StashSnapshot() (domain.StashSnapshot, error) {
	output, err := self.QueryTrim("git", "stash", "list")
	return domain.StashSnapshot(len(stringslice.Lines(output))), err
}

// Version indicates whether the needed Git version is installed.
func (self *BackendCommands) Version() (major int, minor int, err error) {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	output, err := self.QueryTrim("git", "version")
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

func (self *BackendCommands) currentBranchDuringRebase() (domain.LocalBranchName, error) {
	rootDir := self.RootDirectory()
	rawContent, err := os.ReadFile(fmt.Sprintf("%s/.git/rebase-apply/head-name", rootDir))
	if err != nil {
		// Git 2.26 introduces a new rebase backend, see https://github.com/git/git/blob/master/Documentation/RelNotes/2.26.0.txt
		rawContent, err = os.ReadFile(fmt.Sprintf("%s/.git/rebase-merge/head-name", rootDir))
		if err != nil {
			return domain.EmptyLocalBranchName(), err
		}
	}
	content := strings.TrimSpace(string(rawContent))
	return domain.NewLocalBranchName(strings.ReplaceAll(content, "refs/heads/", "")), nil
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
