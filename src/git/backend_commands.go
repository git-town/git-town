package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v9/src/cache"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/stringslice"
)

type BackendRunner interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
	RunMany([][]string) error
}

// BackendCommands are Git commands that Git Town executes to determine which frontend commands to run.
// They don't change the user's repo, execute instantaneously, and Git Town needs to know their output.
// They are invisible to the end user unless the "debug" option is set.
type BackendCommands struct {
	BackendRunner                     // executes shell commands in the directory of the Git repo
	Config             *RepoConfig    // the known state of the Git repository
	CurrentBranchCache *cache.String  // caches the currently checked out Git branch
	RemoteBranchCache  *cache.Strings // caches the remote branches of this Git repo
	RemotesCache       *cache.Strings // caches Git remotes
}

// Author provides the locally Git configured user.
func (bc *BackendCommands) Author() (string, error) {
	name, err := bc.QueryTrim("git", "config", "user.name")
	if err != nil {
		return "", err
	}
	email, err := bc.QueryTrim("git", "config", "user.email")
	if err != nil {
		return "", err
	}
	return name + " <" + email + ">", nil
}

// BranchAuthors provides the user accounts that contributed to the given branch.
// Returns lines of "name <email>".
func (bc *BackendCommands) BranchAuthors(branch, parent string) ([]string, error) {
	output, err := bc.QueryTrim("git", "shortlog", "-s", "-n", "-e", parent+".."+branch)
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

// BranchHasUnmergedCommits indicates whether the branch with the given name
// contains commits that are not merged into the main branch.
func (bc *BackendCommands) BranchHasUnmergedCommits(branch, parent string) (bool, error) {
	out, err := bc.QueryTrim("git", "log", parent+".."+branch)
	if err != nil {
		return false, fmt.Errorf(messages.BranchDiffProblem, branch, err)
	}
	return out != "", nil
}

// BranchesSyncStatus provides detailed information about the sync status of all branches.
func (bc *BackendCommands) BranchesSyncStatus() (branches BranchesSyncStatus, currentBranch string, err error) { //nolint:nonamedreturns
	output, err := bc.Query("git", "branch", "-vva")
	if err != nil {
		return
	}
	branches, currentBranch = ParseVerboseBranchesOutput(output)
	if currentBranch != "" {
		bc.CurrentBranchCache.Set(currentBranch)
	}
	return branches, currentBranch, nil
}

// ParseVerboseBranchesOutput provides the branches in the given Git output as well as the name of the currently checked out branch.
func ParseVerboseBranchesOutput(output string) (BranchesSyncStatus, string) {
	result := BranchesSyncStatus{}
	spaceRE := regexp.MustCompile("[ ]+")
	lines := stringslice.Lines(output)
	checkedoutBranch := ""
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := spaceRE.Split(line[2:], 3)
		branchName := parts[0]
		sha := parts[1]
		remoteText := parts[2]
		if line[0] == '*' && branchName != "(no" { // "(no" is what we get when a rebase is active, in which case no branch is checked out
			checkedoutBranch = branchName
		}
		syncStatus := determineSyncStatus(branchName, remoteText)
		branchName = strings.TrimPrefix(branchName, "remotes/origin/")
		if !result.Contains(branchName) {
			result = append(result, BranchSyncStatus{
				LocalName:  branchName,
				SyncStatus: syncStatus,
			})
		}
	}
	return result, checkedoutBranch
}

func determineSyncStatus(branchName, remoteText string) SyncStatus {
	if remoteText[0] == '[' {
		closingBracketPos := strings.IndexRune(remoteText, ']')
		textInBrackets := remoteText[1:closingBracketPos]
		_, remoteStatus, _ := strings.Cut(textInBrackets, ": ")
		if remoteStatus == "" {
			return SyncStatusUpToDate
		}
		if remoteStatus == "gone" {
			return SyncStatusDeletedAtRemote
		}
		if strings.Contains(remoteStatus, ", behind ") {
			return SyncStatusAheadAndBehind
		}
		if strings.HasPrefix(remoteStatus, "ahead ") {
			return SyncStatusAhead
		}
		if strings.HasPrefix(remoteStatus, "behind ") {
			return SyncStatusBehind
		}
		panic(fmt.Sprintf("cannot determine the sync status for Git remote %q and branch name %q", remoteText, branchName))
	} else {
		if strings.HasPrefix(branchName, "remotes/origin/") {
			return SyncStatusRemoteOnly
		}
		return SyncStatusLocalOnly
	}
}

// CheckoutBranch checks out the Git branch with the given name.
func (bc *BackendCommands) CheckoutBranchUncached(name string) error {
	err := bc.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf(messages.BranchCheckoutProblem, name, err)
	}
	return nil
}

// CheckoutBranch checks out the Git branch with the given name.
func (bc *BackendCommands) CheckoutBranch(name string) error {
	if !bc.Config.DryRun {
		err := bc.CheckoutBranchUncached(name)
		if err != nil {
			return err
		}
	}
	if name != "-" {
		bc.CurrentBranchCache.Set(name)
	} else {
		bc.CurrentBranchCache.Invalidate()
	}
	return nil
}

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// Adds the given prefix with the newline if provided.
func (bc *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {
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

// CreateFeatureBranch creates a feature branch with the given name in this repository.
func (bc *BackendCommands) CreateFeatureBranch(name string) error {
	err := bc.RunMany([][]string{
		{"git", "branch", name, "main"},
		{"git", "config", "git-town-branch." + name + ".parent", "main"},
	})
	if err != nil {
		return fmt.Errorf(messages.BranchFeatureCannotCreate, name, err)
	}
	return nil
}

// CurrentBranch provides the currently checked out branch.
func (bc *BackendCommands) CurrentBranchUncached() (string, error) {
	rebasing, err := bc.HasRebaseInProgress()
	if err != nil {
		return "", fmt.Errorf(messages.BranchCurrentProblem, err)
	}
	if rebasing {
		currentBranch, err := bc.currentBranchDuringRebase()
		if err != nil {
			return "", err
		}
		return currentBranch, nil
	}
	output, err := bc.QueryTrim("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf(messages.BranchCurrentProblem, err)
	}
	return output, nil
}

// CurrentBranch provides the currently checked out branch.
func (bc *BackendCommands) CurrentBranch() (string, error) {
	if bc.Config.DryRun {
		return bc.CurrentBranchCache.Value(), nil
	}
	if !bc.CurrentBranchCache.Initialized() {
		currentBranch, err := bc.CurrentBranchUncached()
		if err != nil {
			return currentBranch, err
		}
		bc.CurrentBranchCache.Set(currentBranch)
	}
	return bc.CurrentBranchCache.Value(), nil
}

func (bc *BackendCommands) currentBranchDuringRebase() (string, error) {
	rootDir := bc.RootDirectory()
	rawContent, err := os.ReadFile(fmt.Sprintf("%s/.git/rebase-apply/head-name", rootDir))
	if err != nil {
		// Git 2.26 introduces a new rebase backend, see https://github.com/git/git/blob/master/Documentation/RelNotes/2.26.0.txt
		rawContent, err = os.ReadFile(fmt.Sprintf("%s/.git/rebase-merge/head-name", rootDir))
		if err != nil {
			return "", err
		}
	}
	content := strings.TrimSpace(string(rawContent))
	return strings.ReplaceAll(content, "refs/heads/", ""), nil
}

// CurrentSha provides the SHA of the currently checked out branch/commit.
func (bc *BackendCommands) CurrentSha() (string, error) {
	return bc.ShaForBranch("HEAD")
}

// ExpectedPreviouslyCheckedOutBranch returns what is the expected previously checked out branch
// given the inputs.
func (bc *BackendCommands) ExpectedPreviouslyCheckedOutBranch(initialPreviouslyCheckedOutBranch, initialBranch, mainBranch string) (string, error) {
	hasInitialPreviouslyCheckedOutBranch, err := bc.HasLocalBranch(initialPreviouslyCheckedOutBranch)
	if err != nil {
		return "", err
	}
	if hasInitialPreviouslyCheckedOutBranch {
		currentBranch, err := bc.CurrentBranch()
		if err != nil {
			return "", err
		}
		hasInitialBranch, err := bc.HasLocalBranch(initialBranch)
		if err != nil {
			return "", err
		}
		if currentBranch == initialBranch || !hasInitialBranch {
			return initialPreviouslyCheckedOutBranch, nil
		}
		return initialBranch, nil
	}
	return mainBranch, nil
}

// HasConflicts returns whether the local repository currently has unresolved merge conflicts.
func (bc *BackendCommands) HasConflicts() (bool, error) {
	output, err := bc.QueryTrim("git", "status")
	if err != nil {
		return false, fmt.Errorf(messages.ConflictDetectionProblem, err)
	}
	return strings.Contains(output, "Unmerged paths"), nil
}

// HasLocalBranch indicates whether this repo has a local branch with the given name.
func (bc *BackendCommands) HasLocalBranch(name string) (bool, error) {
	branches, err := bc.LocalBranches()
	if err != nil {
		return false, fmt.Errorf(messages.BranchLocalProblem, name, err)
	}
	return stringslice.Contains(branches, name), nil
}

// HasMergeInProgress indicates whether this Git repository currently has a merge in progress.
func (bc *BackendCommands) HasMergeInProgress() bool {
	err := bc.Run("git", "rev-parse", "-q", "--verify", "MERGE_HEAD")
	return err == nil
}

// HasOpenChanges indicates whether this repo has open changes.
func (bc *BackendCommands) HasOpenChanges() (bool, error) {
	output, err := bc.QueryTrim("git", "status", "--porcelain", "--ignore-submodules")
	if err != nil {
		return false, fmt.Errorf(messages.OpenChangesProblem, err)
	}
	return output != "", nil
}

// HasRebaseInProgress indicates whether this Git repository currently has a rebase in progress.
func (bc *BackendCommands) HasRebaseInProgress() (bool, error) {
	output, err := bc.QueryTrim("git", "status")
	if err != nil {
		return false, fmt.Errorf(messages.RebaseProblem, err)
	}
	if strings.Contains(output, "You are currently rebasing") {
		return true, nil
	}
	if strings.Contains(output, "rebase in progress") {
		return true, nil
	}
	return false, nil
}

// HasShippableChanges indicates whether the given branch has changes
// not currently in the main branch.
func (bc *BackendCommands) HasShippableChanges(branch, mainBranch string) (bool, error) {
	out, err := bc.QueryTrim("git", "diff", mainBranch+".."+branch)
	if err != nil {
		return false, fmt.Errorf(messages.ShippableChangesProblem, branch, err)
	}
	return out != "", nil
}

// LastCommitMessage provides the commit message for the last commit.
func (bc *BackendCommands) LastCommitMessage() (string, error) {
	out, err := bc.QueryTrim("git", "log", "-1", "--format=%B")
	if err != nil {
		return "", fmt.Errorf(messages.CommitMessageProblem, err)
	}
	return out, nil
}

// LocalBranches provides the names of all branches in the local repository,
// ordered alphabetically.
// TODO: can we derive this info from allBranchesSyncStatus?
func (bc *BackendCommands) LocalBranches() ([]string, error) {
	output, err := bc.QueryTrim("git", "branch")
	if err != nil {
		return []string{}, err
	}
	result := []string{}
	for _, line := range stringslice.Lines(output) {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, line)
	}
	return result, nil
}

// LocalBranchesMainFirst provides the names of all local branches in this repo.
func (bc *BackendCommands) LocalBranchesMainFirst(mainBranch string) ([]string, error) {
	branches, err := bc.LocalBranches()
	if err != nil {
		return []string{}, err
	}
	return stringslice.Hoist(sort.StringSlice(branches), mainBranch), nil
}

// PreviouslyCheckedOutBranch provides the name of the branch that was previously checked out in this repo.
func (bc *BackendCommands) PreviouslyCheckedOutBranch() string {
	output, err := bc.QueryTrim("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return ""
	}
	return output
}

// Remotes provides the names of all Git remotes in this repository.
func (bc *BackendCommands) RemotesUncached() ([]string, error) {
	out, err := bc.QueryTrim("git", "remote")
	if err != nil {
		return []string{}, fmt.Errorf(messages.RemotesProblem, err)
	}
	if out == "" {
		return []string{}, nil
	}
	return stringslice.Lines(out), nil
}

// Remotes provides the names of all Git remotes in this repository.
func (bc *BackendCommands) Remotes() (config.Remotes, error) {
	if !bc.RemotesCache.Initialized() {
		remotes, err := bc.RemotesUncached()
		if err != nil {
			return remotes, err
		}
		bc.RemotesCache.Set(remotes)
	}
	return bc.RemotesCache.Value(), nil
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (bc *BackendCommands) RemoveOutdatedConfiguration(allBranches BranchesSyncStatus) error {
	for child, parent := range bc.Config.Lineage() {
		hasChildBranch := allBranches.Contains(child)
		hasParentBranch := allBranches.Contains(parent)
		if !hasChildBranch || !hasParentBranch {
			// TODO
			err := bc.Config.RemoveParent(child)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RootDirectory provides the path of the rood directory of the current repository,
// i.e. the directory that contains the ".git" folder.
func (bc *BackendCommands) RootDirectory() string {
	output, err := bc.QueryTrim("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return ""
	}
	return filepath.FromSlash(output)
}

// ShaForBranch provides the SHA for the local branch with the given name.
func (bc *BackendCommands) ShaForBranch(name string) (string, error) {
	output, err := bc.QueryTrim("git", "rev-parse", name)
	if err != nil {
		return "", fmt.Errorf(messages.BranchLocalShaProblem, name, err)
	}
	return output, nil
}

// ShouldPushBranch returns whether the local branch with the given name
// contains commits that have not been pushed to its tracking branch.
func (bc *BackendCommands) ShouldPushBranch(branch string) (bool, error) {
	trackingBranch := TrackingBranchName(branch)
	out, err := bc.QueryTrim("git", "rev-list", "--left-right", branch+"..."+trackingBranch)
	if err != nil {
		return false, fmt.Errorf(messages.DiffProblem, branch, trackingBranch, err)
	}
	return out != "", nil
}

// Version indicates whether the needed Git version is installed.
func (bc *BackendCommands) Version() (major int, minor int, err error) {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	output, err := bc.QueryTrim("git", "version")
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
