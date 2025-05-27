package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/cache"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Commands are Git commands that Git Town executes to determine which frontend commands to run.
// They don't change the user's repo, execute instantaneously, and Git Town needs to know their output.
// They are invisible to the end user unless the "verbose" option is set.
type Commands struct {
	CurrentBranchCache *cache.WithPrevious[gitdomain.LocalBranchName] // caches the currently checked out Git branch
	RemotesCache       *cache.Cache[gitdomain.Remotes]                // caches Git remotes
}

func (self *Commands) AbortMerge(runner gitdomain.Runner) error {
	return runner.Run("git", "merge", "--abort")
}

func (self *Commands) AbortRebase(runner gitdomain.Runner) error {
	return runner.Run("git", "rebase", "--abort")
}

// BranchAuthors provides the user accounts that contributed to the given branch.
func (self *Commands) BranchAuthors(querier gitdomain.Querier, branch, parent gitdomain.LocalBranchName) ([]gitdomain.Author, error) {
	output, err := querier.QueryTrim("git", "shortlog", "-s", "-n", "-e", parent.String()+".."+branch.String())
	if err != nil {
		return []gitdomain.Author{}, err
	}
	lines := stringslice.Lines(output)
	result := make([]gitdomain.Author, len(lines))
	for l, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result[l] = gitdomain.Author(parts[1])
	}
	return result, nil
}

func (self *Commands) BranchContainsMerges(querier gitdomain.Querier, branch, parent gitdomain.LocalBranchName) (bool, error) {
	output, err := querier.QueryTrim("git", "log", "--merges", "--format=%H", fmt.Sprintf("%s..%s", parent, branch))
	return len(output) > 0, err
}

func (self *Commands) BranchExists(runner gitdomain.Runner, branch gitdomain.LocalBranchName) bool {
	err := runner.Run("git", "rev-parse", "--verify", "-q", "refs/heads/"+branch.String())
	return err == nil
}

func (self *Commands) BranchExistsRemotely(runner gitdomain.Runner, branch gitdomain.LocalBranchName, remote gitdomain.Remote) bool {
	err := runner.Run("git", "ls-remote", remote.String(), branch.String())
	return err == nil
}

// BranchHasUnmergedChanges indicates whether the branch with the given name
// contains changes that were not merged into the main branch.
func (self *Commands) BranchHasUnmergedChanges(querier gitdomain.Querier, branch, parent gitdomain.LocalBranchName) (bool, error) {
	out, err := querier.QueryTrim("git", "diff", "--shortstat", parent.String(), branch.String())
	if err != nil {
		return false, fmt.Errorf(messages.BranchDiffProblem, branch, err)
	}
	return len(out) > 0, nil
}

func (self *Commands) BranchInSyncWithParent(querier gitdomain.Querier, branch gitdomain.LocalBranchName, parent gitdomain.BranchName) (bool, error) {
	output, err := querier.QueryTrim("git", "log", "--no-merges", "--format=%H", parent.String(), "^"+branch.String())
	return len(output) == 0, err
}

// BranchInSyncWithTracking returns whether the local branch with the given name
// contains commits that have not been pushed to its tracking branch.
func (self *Commands) BranchInSyncWithTracking(querier gitdomain.Querier, branch gitdomain.LocalBranchName, devRemote gitdomain.Remote) (bool, error) {
	out, err := querier.QueryTrim("git", "rev-list", "--left-right", branch.String()+"..."+branch.TrackingBranch(devRemote).String())
	if err != nil {
		return false, fmt.Errorf(messages.DiffProblem, branch, branch, err)
	}
	return len(out) == 0, nil
}

func (self *Commands) BranchesSnapshot(querier gitdomain.Querier) (gitdomain.BranchesSnapshot, error) {
	branches, err := branchesQuery(querier)
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), err
	}
	if len(branches) == 0 {
		// We are in a brand-new repo.
		// Report the initial branch name (reported by `git branch --show-current`) as the current branch.
		currentBranch, err := self.CurrentBranchUncached(querier)
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), err
		}
		return makeBranchesSnapshotNewRepo(currentBranch), nil
	}
	result := gitdomain.BranchInfos{}
	currentBranchOpt := None[gitdomain.LocalBranchName]()
	for _, branch := range branches {
		if branch.Symref {
			// Ignore symbolic refs.
			continue
		}
		if branch.Head && branch.BranchName.IsLocal() {
			currentBranchOpt = Some(branch.BranchName.LocalName())
		}
		switch {
		case branch.Worktree && !branch.Head:
			result = append(result, gitdomain.BranchInfo{
				LocalName:  Some(branch.BranchName.LocalName()),
				LocalSHA:   Some(branch.SHA),
				RemoteName: branch.UpstreamOption,
				RemoteSHA:  None[gitdomain.SHA](), // may be added later
				SyncStatus: gitdomain.SyncStatusOtherWorktree,
			})
		case isLocalRefName(branch.RefName):
			syncStatus := determineSyncStatus(branch.Track, branch.UpstreamOption)
			result = append(result, gitdomain.BranchInfo{
				LocalName:  Some(branch.BranchName.LocalName()),
				LocalSHA:   Some(branch.SHA),
				RemoteName: branch.UpstreamOption,
				RemoteSHA:  None[gitdomain.SHA](), // may be added later
				SyncStatus: syncStatus,
			})
		default:
			// Not using `BranchName.RemoteName()` because it might not necessarily be prefixed with "origin/".
			remoteBranchName := gitdomain.NewRemoteBranchName(branch.BranchName.String())
			if existingBranchWithTracking, hasExistingBranchWithTracking := result.FindByRemoteName(remoteBranchName).Get(); hasExistingBranchWithTracking {
				existingBranchWithTracking.RemoteSHA = Some(branch.SHA)
			} else {
				result = append(result, gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(remoteBranchName),
					RemoteSHA:  Some(branch.SHA),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				})
			}
		}
	}
	if currentBranchOpt.IsNone() {
		rebaseInProgress, err := self.HasRebaseInProgress(querier)
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), err
		}
		if !rebaseInProgress {
			// We are in a detached HEAD state. Use the current HEAD location as the branch name.
			headSHA, err := self.CurrentSHA(querier)
			if err != nil {
				return gitdomain.EmptyBranchesSnapshot(), err
			}
			currentBranchOpt = gitdomain.NewLocalBranchNameOption(headSHA.String())
			// prepend to result
			result = slices.Insert(result, 0, gitdomain.BranchInfo{
				LocalName:  currentBranchOpt,
				LocalSHA:   Some(headSHA),
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
			})
		}
	}
	if currentBranch, hasCurrentBranch := currentBranchOpt.Get(); hasCurrentBranch {
		self.CurrentBranchCache.Set(currentBranch)
	}
	return gitdomain.BranchesSnapshot{
		Branches: result,
		Active:   currentBranchOpt,
	}, nil
}

func (self *Commands) ChangeDir(dir gitdomain.RepoRootDir) error {
	return os.Chdir(dir.String())
}

func (self *Commands) CheckoutBranch(runner gitdomain.Runner, name gitdomain.LocalBranchName, merge configdomain.SwitchUsingMerge) error {
	err := self.CheckoutBranchUncached(runner, name, merge)
	if err != nil {
		return err
	}
	if name.String() != "-" {
		self.CurrentBranchCache.Set(name)
	} else {
		self.CurrentBranchCache.Invalidate()
	}
	return nil
}

func (self *Commands) CheckoutBranchUncached(runner gitdomain.Runner, name gitdomain.LocalBranchName, merge configdomain.SwitchUsingMerge) error {
	args := []string{"checkout", name.String()}
	if merge {
		args = append(args, "-m")
	}
	err := runner.Run("git", args...)
	if err != nil {
		return fmt.Errorf(messages.BranchCheckoutProblem, name, err)
	}
	return nil
}

func (self *Commands) CherryPick(runner gitdomain.Runner, sha gitdomain.SHA) error {
	return runner.Run("git", "cherry-pick", sha.String())
}

func (self *Commands) CherryPickAbort(runner gitdomain.Runner) error {
	return runner.Run("git", "cherry-pick", "--abort")
}

func (self *Commands) CherryPickContinue(runner gitdomain.Runner) error {
	return runner.RunWithEnv([]string{"GIT_EDITOR=true"}, "git", "cherry-pick", "--continue")
}

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// If the given prefix has content, adds it together with a newline.
func (self *Commands) CommentOutSquashCommitMessage(prefix Option[string]) error {
	squashMessageFile := ".git/SQUASH_MSG"
	contentBytes, err := os.ReadFile(squashMessageFile)
	if err != nil {
		return fmt.Errorf(messages.SquashCannotReadFile, squashMessageFile, err)
	}
	content := string(contentBytes)
	if prefix, hasPrefix := prefix.Get(); hasPrefix {
		content = prefix + "\n" + content
	}
	content = regexp.MustCompile("(?m)^").ReplaceAllString(content, "# ")
	return os.WriteFile(squashMessageFile, []byte(content), 0o600)
}

func (self *Commands) Commit(runner gitdomain.Runner, useMessage configdomain.UseMessage, author Option[gitdomain.Author], commitHook configdomain.CommitHook) error {
	args := []string{"commit"}
	switch {
	case useMessage.IsCustomMessage():
		message := useMessage.GetCustomMessageOrPanic()
		args = append(args, "-m", message.String())
	case useMessage.IsUseDefault():
		args = append(args, "--no-edit")
	case useMessage.IsEditDefault():
		// This is the default behaviour of `git commit`.
	default:
		return fmt.Errorf("unhandled %#v case", useMessage)
	}
	if author, hasAuthor := author.Get(); hasAuthor {
		args = append(args, "--author", author.String())
	}
	switch commitHook {
	case configdomain.CommitHookDisabled:
		args = append(args, "--no-verify")
	case configdomain.CommitHookEnabled:
	}
	return runner.Run("git", args...)
}

func (self *Commands) CommitMessage(querier gitdomain.Querier, sha gitdomain.SHA) (gitdomain.CommitMessage, error) {
	output, err := querier.QueryTrim("git", "show", "--no-patch", "--format=%B", sha.String())
	return gitdomain.CommitMessage(strings.TrimSpace(output)), err
}

func (self *Commands) CommitStart(runner gitdomain.Runner) error {
	return runner.Run("git", "commit")
}

func (self *Commands) CommitsInBranch(querier gitdomain.Querier, branch gitdomain.LocalBranchName, parent Option[gitdomain.LocalBranchName]) (gitdomain.Commits, error) {
	if parent, hasParent := parent.Get(); hasParent {
		return self.CommitsInFeatureBranch(querier, branch, parent.BranchName())
	}
	return self.CommitsInPerennialBranch(querier)
}

func (self *Commands) CommitsInFeatureBranch(querier gitdomain.Querier, branch gitdomain.LocalBranchName, parent gitdomain.BranchName) (gitdomain.Commits, error) {
	output, err := querier.QueryTrim("git", "log", "--format=%H %s", fmt.Sprintf("%s..%s", parent.String(), branch.String()))
	if err != nil {
		return gitdomain.Commits{}, err
	}
	lines := strings.Split(output, "\n")
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
	slices.Reverse(result)
	return result, nil
}

func (self *Commands) CommitsInPerennialBranch(querier gitdomain.Querier) (gitdomain.Commits, error) {
	output, err := querier.QueryTrim("git", "log", "--format=%H %s", "-10")
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

func (self *Commands) ContentBlobInfo(querier gitdomain.Querier, branch gitdomain.Location, filePath string) (Option[BlobInfo], error) {
	output, err := querier.QueryTrim("git", "ls-tree", branch.String(), filePath)
	if err != nil || len(output) == 0 {
		return None[BlobInfo](), err
	}
	sha, err := ParseLsTreeOutput(output)
	return Some(sha), err
}

func (self *Commands) ContinueRebase(runner gitdomain.Runner) error {
	return runner.RunWithEnv([]string{"GIT_EDITOR=true"}, "git", "rebase", "--continue")
}

// CreateAndCheckoutBranch creates a new branch with the given name and checks it out using a single Git operation.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *Commands) CreateAndCheckoutBranch(runner gitdomain.Runner, name gitdomain.LocalBranchName) error {
	err := runner.Run("git", "checkout", "-b", name.String())
	self.CurrentBranchCache.Set(name)
	return err
}

// CreateAndCheckoutBranchWithParent creates a new branch with the given name and parent and checks it out using a single Git operation.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *Commands) CreateAndCheckoutBranchWithParent(runner gitdomain.Runner, name gitdomain.LocalBranchName, parent gitdomain.Location) error {
	args := []string{"checkout", "-b", name.String(), parent.String()}
	if parent.IsRemoteBranchName() {
		args = append(args, "--no-track")
	}
	err := runner.Run("git", args...)
	self.CurrentBranchCache.Set(name)
	return err
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *Commands) CreateBranch(runner gitdomain.Runner, name gitdomain.LocalBranchName, parent gitdomain.Location) error {
	return runner.Run("git", "branch", name.String(), parent.String())
}

func (self *Commands) CreateTrackingBranch(runner gitdomain.Runner, branch gitdomain.LocalBranchName, remote gitdomain.Remote, noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, "-u", remote.String())
	args = append(args, branch.String())
	return runner.Run("git", args...)
}

func (self *Commands) CurrentBranch(querier gitdomain.Querier) (gitdomain.LocalBranchName, error) {
	if !self.CurrentBranchCache.Initialized() {
		currentBranch, err := self.CurrentBranchUncached(querier)
		if err != nil {
			return currentBranch, err
		}
		self.CurrentBranchCache.Set(currentBranch)
	}
	return self.CurrentBranchCache.Value(), nil
}

func (self *Commands) CurrentBranchDuringRebase(querier gitdomain.Querier) (gitdomain.LocalBranchName, error) {
	gitDir, err := self.gitDirectory(querier)
	if err != nil {
		return "", err
	}
	for _, rebaseHeadFileName := range []string{"rebase-merge/head-name", "rebase-apply/head-name"} {
		rebaseHeadFilePath := filepath.Join(gitDir, rebaseHeadFileName)
		content, err := os.ReadFile(rebaseHeadFilePath)
		if err != nil {
			continue
		}
		refName := strings.TrimSpace(string(content))
		if strings.HasPrefix(refName, "refs/heads/") {
			branchName := strings.TrimPrefix(refName, "refs/heads/")
			return gitdomain.NewLocalBranchName(branchName), nil
		}
		// rebase head name is not a branch name
		break
	}
	return "", errors.New(messages.BranchCurrentProblemNoError)
}

func (self *Commands) CurrentBranchHasTrackingBranch(runner gitdomain.Runner) bool {
	err := runner.Run("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	return err == nil
}

func (self *Commands) CurrentBranchUncached(querier gitdomain.Querier) (gitdomain.LocalBranchName, error) {
	// first try to detect the current branch the normal way
	output, err := querier.QueryTrim("git", "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf(messages.BranchCurrentProblem, err)
	}
	if output != "" {
		return gitdomain.NewLocalBranchName(output), nil
	}
	// here we couldn't detect the current branch the normal way --> assume we are in a rebase and try the rebase way
	return self.CurrentBranchDuringRebase(querier)
}

// CurrentSHA provides the SHA of the currently checked out branch/commit.
func (self *Commands) CurrentSHA(querier gitdomain.Querier) (gitdomain.SHA, error) {
	return self.SHAForBranch(querier, "HEAD")
}

func (self *Commands) DefaultBranch(querier gitdomain.Querier) Option[gitdomain.LocalBranchName] {
	name, err := querier.QueryTrim("git", "config", "--get", "init.defaultbranch")
	if err != nil {
		return None[gitdomain.LocalBranchName]()
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return None[gitdomain.LocalBranchName]()
	}
	return Some(gitdomain.LocalBranchName(name))
}

func (self *Commands) DefaultRemote(querier gitdomain.Querier) gitdomain.Remote {
	output, err := querier.QueryTrim("git", "config", "--get", "clone.defaultRemoteName")
	if err != nil {
		// Git returns an error if the user has not configured a default remote name.
		// In this case use the Git default of "origin".
		return gitdomain.RemoteOrigin
	}
	return gitdomain.Remote(output)
}

// DeleteConfigEntryForgeType removes the forge type config entry.
func (self *Commands) DeleteConfigEntryForgeType(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyForgeType.String())
}

// DeleteConfigEntryOriginHostname removes the origin hostname override
func (self *Commands) DeleteConfigEntryOriginHostname(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyHostingOriginHostname.String())
}

// DeleteLastCommit resets HEAD to the previous commit.
func (self *Commands) DeleteLastCommit(runner gitdomain.Runner) error {
	return runner.Run("git", "reset", "--hard", "HEAD~1")
}

func (self *Commands) DeleteLocalBranch(runner gitdomain.Runner, name gitdomain.LocalBranchName) error {
	return runner.Run("git", "branch", "-D", name.String())
}

func (self *Commands) DeleteTrackingBranch(runner gitdomain.Runner, name gitdomain.RemoteBranchName) error {
	remote, localBranchName := name.Parts()
	return runner.Run("git", "push", remote.String(), ":"+localBranchName.String())
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (self *Commands) DiffParent(runner gitdomain.Runner, branch, parentBranch gitdomain.LocalBranchName) error {
	return runner.Run("git", "diff", parentBranch.String(), branch.String())
}

func (self *Commands) DiscardOpenChanges(runner gitdomain.Runner) error {
	return runner.Run("git", "reset", "--hard")
}

func (self *Commands) DropMostRecentStash(runner gitdomain.Runner) error {
	return runner.Run("git", "stash", "drop")
}

func (self *Commands) Fetch(runner gitdomain.Runner, syncTags configdomain.SyncTags) error {
	if syncTags.IsTrue() {
		return runner.Run("git", "fetch", "--prune", "--tags")
	}
	return runner.Run("git", "fetch", "--prune", "--no-tags")
}

func (self *Commands) FetchUpstream(runner gitdomain.Runner, branch gitdomain.LocalBranchName) error {
	return runner.Run("git", "fetch", gitdomain.RemoteUpstream.String(), branch.String())
}

func (self *Commands) FileConflictFullInfo(querier gitdomain.Querier, quickInfo FileConflictQuickInfo, parentLocation gitdomain.Location, mainBranch gitdomain.LocalBranchName) (FileConflictFullInfo, error) {
	mainBlob := None[BlobInfo]()
	parentBlob := None[BlobInfo]()
	if currentBranchBlobInfo, has := quickInfo.CurrentBranchChange.Get(); has {
		var err error
		mainBlob, err = self.ContentBlobInfo(querier, mainBranch.Location(), currentBranchBlobInfo.FilePath)
		if err != nil {
			return FileConflictFullInfo{}, err
		}
		parentBlob, err = self.ContentBlobInfo(querier, parentLocation, currentBranchBlobInfo.FilePath)
		if err != nil {
			return FileConflictFullInfo{}, err
		}
	}
	result := FileConflictFullInfo{
		Current: quickInfo.CurrentBranchChange,
		Main:    mainBlob,
		Parent:  parentBlob,
	}
	return result, nil
}

func (self *Commands) FileConflictFullInfos(querier gitdomain.Querier, quickInfos []FileConflictQuickInfo, parentLocation gitdomain.Location, mainBranch gitdomain.LocalBranchName) ([]FileConflictFullInfo, error) {
	result := make([]FileConflictFullInfo, len(quickInfos))
	for q, quickInfo := range quickInfos {
		fullInfo, err := self.FileConflictFullInfo(querier, quickInfo, parentLocation, mainBranch)
		if err != nil {
			return result, err
		}
		result[q] = fullInfo
	}
	return result, nil
}

func (self *Commands) FileConflictQuickInfos(querier gitdomain.Querier) ([]FileConflictQuickInfo, error) {
	output, err := querier.Query("git", "ls-files", "--unmerged")
	if err != nil {
		return []FileConflictQuickInfo{}, err
	}
	return ParseLsFilesUnmergedOutput(output)
}

// provides the commit message of the first commit in the branch with the given name
func (self *Commands) FirstCommitMessageInBranch(runner gitdomain.Querier, branch, parent gitdomain.BranchName) (Option[gitdomain.CommitMessage], error) {
	output, err := runner.QueryTrim("git", "log", fmt.Sprintf("%s..%s", parent, branch), "--format=%s", "--reverse")
	if err != nil {
		return None[gitdomain.CommitMessage](), err
	}
	lines := stringslice.Lines(output)
	if len(lines) == 0 {
		return None[gitdomain.CommitMessage](), nil
	}
	return Some(gitdomain.CommitMessage(lines[0])), nil
}

// provides the first branch in the given list of branch names that actually exists
func (self *Commands) FirstExistingBranch(runner gitdomain.Runner, branches ...gitdomain.LocalBranchName) Option[gitdomain.LocalBranchName] {
	for _, branch := range branches {
		if self.BranchExists(runner, branch) {
			return Some(branch)
		}
	}
	return None[gitdomain.LocalBranchName]()
}

func (self *Commands) ForcePushBranchSafely(runner gitdomain.Runner, noPushHook configdomain.NoPushHook, forceIfIncludes bool) error {
	args := []string{"push", "--force-with-lease"}
	if forceIfIncludes {
		args = append(args, "--force-if-includes")
	}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return runner.Run("git", args...)
}

func (self *Commands) GitVersion(querier gitdomain.Querier) (Version, error) {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\w+)`)
	output, err := querier.QueryTrim("git", "version")
	if err != nil {
		return EmptyVersion(), fmt.Errorf(messages.GitVersionProblem, err)
	}
	matches := versionRegexp.FindStringSubmatch(output)
	if matches == nil {
		return EmptyVersion(), fmt.Errorf(messages.GitVersionUnexpectedOutput, output)
	}
	majorVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		return EmptyVersion(), fmt.Errorf(messages.GitVersionMajorNotNumber, matches[1], err)
	}
	minorVersion, err := strconv.Atoi(matches[2])
	if err != nil {
		return EmptyVersion(), fmt.Errorf(messages.GitVersionMinorNotNumber, matches[2], err)
	}
	return Version{
		Major: majorVersion,
		Minor: minorVersion,
	}, nil
}

func (self *Commands) HasMergeInProgress(runner gitdomain.Runner) bool {
	err := runner.Run("git", "rev-parse", "--verify", "-q", "MERGE_HEAD")
	return err == nil
}

func (self *Commands) HasRebaseInProgress(querier gitdomain.Querier) (bool, error) {
	gitDir, err := self.gitDirectory(querier)
	if err != nil {
		return false, err
	}
	for _, rebaseDirName := range []string{"rebase-merge", "rebase-apply"} {
		rebaseDirPath := filepath.Join(gitDir, rebaseDirName)
		stat, err := os.Stat(rebaseDirPath)
		if err == nil && stat.IsDir() {
			return true, nil
		}
	}
	return false, nil
}

func (self *Commands) MergeBranchNoEdit(runner gitdomain.Runner, branch gitdomain.BranchName) error {
	return runner.Run("git", "merge", "--no-edit", "--ff", branch.String())
}

func (self *Commands) MergeFastForward(runner gitdomain.Runner, branch gitdomain.BranchName) error {
	return runner.Run("git", "merge", "--ff-only", branch.String())
}

func (self *Commands) MergeNoFastForward(runner gitdomain.Runner, useMessage configdomain.UseMessage, branch gitdomain.LocalBranchName) error {
	args := []string{"merge", "--no-ff"}
	switch {
	case useMessage.IsCustomMessage():
		message := useMessage.GetCustomMessageOrPanic()
		args = append(args, "-m", message.String())
	case useMessage.IsUseDefault():
		args = append(args, "--no-edit")
	case useMessage.IsEditDefault():
		// Unlike `git commit`, `git merge` only launches the editor in a tty.
		// Until cucumber tests run git subcommands in a tty we have to explicitly set
		// `--edit` mode to test commit message behaviour.
		args = append(args, "--edit")
	default:
		return fmt.Errorf("unhandled %#v case", useMessage)
	}
	// Add branch name as the last argument.
	args = append(args, "--", branch.String())
	return runner.Run("git", args...)
}

func (self *Commands) OriginHead(querier gitdomain.Querier) Option[gitdomain.LocalBranchName] {
	output, err := querier.QueryTrim("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {
		return None[gitdomain.LocalBranchName]()
	}
	output = strings.TrimSpace(output)
	if output == "" {
		return None[gitdomain.LocalBranchName]()
	}
	return Some(gitdomain.LocalBranchName(LastBranchInRef(output)))
}

func (self *Commands) PopStash(runner gitdomain.Runner) error {
	err := runner.Run("git", "stash", "pop")
	if err != nil {
		_ = runner.Run("git", "stash", "drop")
	}
	return err
}

// PreviouslyCheckedOutBranch provides the name of the branch that was checked out before the current branch was checked out.
func (self *Commands) PreviouslyCheckedOutBranch(querier gitdomain.Querier) Option[gitdomain.LocalBranchName] {
	output, err := querier.QueryTrim("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return None[gitdomain.LocalBranchName]()
	}
	if output == "" {
		return None[gitdomain.LocalBranchName]()
	}
	return gitdomain.NewLocalBranchNameOption(output)
}

func (self *Commands) Pull(runner gitdomain.Runner) error {
	return runner.Run("git", "pull")
}

// PushCurrentBranch pushes the current branch to its tracking branch.
func (self *Commands) PushCurrentBranch(runner gitdomain.Runner, noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return runner.Run("git", args...)
}

func (self *Commands) PushLocalBranch(runner gitdomain.Runner, localSHA gitdomain.SHA, branch gitdomain.LocalBranchName, remote gitdomain.Remote, noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, remote.String(), localSHA.String()+":refs/heads/"+branch.String())
	return runner.Run("git", args...)
}

// PushTags pushes new the Git tags to origin.
func (self *Commands) PushTags(runner gitdomain.Runner, noPushHook configdomain.NoPushHook) error {
	args := []string{"push", "--tags"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return runner.Run("git", args...)
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (self *Commands) Rebase(runner gitdomain.Runner, target gitdomain.BranchName) error {
	return runner.Run("git", "-c", "rebase.updateRefs=false", "rebase", target.String())
}

// Rebase initiates a Git rebase of the current branch onto the given branch.
func (self *Commands) RebaseOnto(runner gitdomain.Runner, branchToRebaseOnto gitdomain.Location, commitsToRemove gitdomain.Location, upstream Option[gitdomain.LocalBranchName]) error {
	args := []string{"-c", "rebase.updateRefs=false", "rebase", "--onto", branchToRebaseOnto.String()}
	if upstream, hasUpstream := upstream.Get(); hasUpstream {
		args = append(args, upstream.String())
	}
	args = append(args, commitsToRemove.String())
	return runner.Run("git", args...)
}

func (self *Commands) Remotes(querier gitdomain.Querier) (gitdomain.Remotes, error) {
	if !self.RemotesCache.Initialized() {
		remotes, err := self.RemotesUncached(querier)
		if err != nil {
			return remotes, err
		}
		self.RemotesCache.Set(&remotes)
	}
	return *self.RemotesCache.Value(), nil
}

func (self *Commands) RemotesUncached(querier gitdomain.Querier) (gitdomain.Remotes, error) {
	out, err := querier.QueryTrim("git", "remote")
	if err != nil {
		return gitdomain.Remotes{}, fmt.Errorf(messages.RemotesProblem, err)
	}
	if out == "" {
		return gitdomain.Remotes{}, nil
	}
	return gitdomain.NewRemotes(stringslice.Lines(out)...), nil
}

func (self *Commands) RemoveBitbucketAppPassword(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyBitbucketAppPassword.String())
}

func (self *Commands) RemoveBitbucketUsername(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyBitbucketUsername.String())
}

func (self *Commands) RemoveCodebergToken(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyCodebergToken.String())
}

// RemoveCommit removes the given commit from the current branch
func (self *Commands) RemoveCommit(runner gitdomain.Runner, commit gitdomain.SHA) error {
	return runner.Run("git", "-c", "rebase.updateRefs=false", "rebase", "--onto", commit.String()+"^", commit.String())
}

func (self *Commands) RemoveFile(runner gitdomain.Runner, fileName string) error {
	return runner.Run("git", "rm", fileName)
}

func (self *Commands) RemoveGitAlias(runner gitdomain.Runner, aliasableCommand configdomain.AliasableCommand) error {
	return runner.Run("git", "config", "--global", "--unset", aliasableCommand.Key().String())
}

func (self *Commands) RemoveGitHubToken(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyGithubToken.String())
}

func (self *Commands) RemoveGitLabToken(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyGitlabToken.String())
}

func (self *Commands) RemoveGiteaToken(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyGiteaToken.String())
}

func (self *Commands) RenameBranch(runner gitdomain.Runner, oldName, newName gitdomain.LocalBranchName) error {
	return runner.Run("git", "branch", "--move", oldName.String(), newName.String())
}

func (self *Commands) RepoStatus(backend gitdomain.RunnerQuerier) (gitdomain.RepoStatus, error) {
	output, err := backend.Query("git", "status", "-z", "--ignore-submodules")
	if err != nil {
		return gitdomain.RepoStatus{}, fmt.Errorf(messages.ConflictDetectionProblem, err)
	}
	statuses, err := ParseGitStatusZ(output)
	if err != nil {
		return gitdomain.RepoStatus{}, fmt.Errorf(messages.ConflictDetectionProblem, err)
	}
	hasConflicts := slices.ContainsFunc(statuses, FileStatusIsUnmerged)
	hasOpenChanges := len(statuses) > 0
	hasUntrackedChanges := slices.ContainsFunc(statuses, FileStatusIsUntracked)
	mergeInProgress := self.HasMergeInProgress(backend)
	rebaseInProgress, err := self.HasRebaseInProgress(backend)
	if err != nil {
		return gitdomain.RepoStatus{}, err
	}
	return gitdomain.RepoStatus{
		Conflicts:        hasConflicts,
		OpenChanges:      hasOpenChanges && !mergeInProgress && !rebaseInProgress,
		RebaseInProgress: rebaseInProgress,
		UntrackedChanges: hasUntrackedChanges,
	}, nil
}

func (self *Commands) ResetBranch(runner gitdomain.Runner, target gitdomain.BranchName) error {
	return runner.Run("git", "reset", "--soft", target.String())
}

func (self *Commands) ResetCurrentBranchToSHA(runner gitdomain.Runner, sha gitdomain.SHA) error {
	return runner.Run("git", "reset", "--hard", sha.String())
}

func (self *Commands) ResetRemoteBranchToSHA(runner gitdomain.Runner, branch gitdomain.RemoteBranchName, sha gitdomain.SHA) error {
	remote := branch.Remote()
	return runner.Run("git", "push", "--force-with-lease", remote.String(), sha.String()+":"+branch.LocalBranchName().String())
}

func (self *Commands) ResolveConflict(runner gitdomain.Runner, file string, resolution gitdomain.ConflictResolution) error {
	return runner.Run("git", "checkout", resolution.GitFlag(), file)
}

func (self *Commands) RevertCommit(runner gitdomain.Runner, sha gitdomain.SHA) error {
	return runner.Run("git", "revert", sha.String())
}

// RootDirectory provides the path of the root directory of the current repository.
func (self *Commands) RootDirectory(querier gitdomain.Querier) Option[gitdomain.RepoRootDir] {
	output, err := querier.QueryTrim("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return None[gitdomain.RepoRootDir]()
	}
	return Some(gitdomain.NewRepoRootDir(filepath.FromSlash(output)))
}

func (self *Commands) SHAForBranch(querier gitdomain.Querier, name gitdomain.BranchName) (gitdomain.SHA, error) {
	output, err := querier.QueryTrim("git", "rev-parse", name.String())
	if err != nil {
		return gitdomain.SHA(""), fmt.Errorf(messages.BranchLocalSHAProblem, name, err)
	}
	return gitdomain.NewSHA(output), nil
}

func (self *Commands) SetBitbucketAppPassword(runner gitdomain.Runner, value configdomain.BitbucketAppPassword, scope configdomain.ConfigScope) error {
	return runner.Run("git", "config", scope.GitFlag(), configdomain.KeyBitbucketAppPassword.String(), value.String())
}

func (self *Commands) SetBitbucketUsername(runner gitdomain.Runner, value configdomain.BitbucketUsername, scope configdomain.ConfigScope) error {
	return runner.Run("git", "config", scope.GitFlag(), configdomain.KeyBitbucketUsername.String(), value.String())
}

func (self *Commands) SetCodebergToken(runner gitdomain.Runner, value configdomain.CodebergToken, scope configdomain.ConfigScope) error {
	return runner.Run("git", "config", scope.GitFlag(), configdomain.KeyCodebergToken.String(), value.String())
}

func (self *Commands) SetForgeType(runner gitdomain.Runner, platform forgedomain.ForgeType) error {
	return runner.Run("git", "config", configdomain.KeyForgeType.String(), platform.String())
}

func (self *Commands) SetGitAlias(runner gitdomain.Runner, aliasableCommand configdomain.AliasableCommand) error {
	return runner.Run("git", "config", "--global", aliasableCommand.Key().String(), "town "+aliasableCommand.String())
}

func (self *Commands) SetGitHubToken(runner gitdomain.Runner, value configdomain.GitHubToken, scope configdomain.ConfigScope) error {
	return runner.Run("git", "config", scope.GitFlag(), configdomain.KeyGithubToken.String(), value.String())
}

func (self *Commands) SetGitLabToken(runner gitdomain.Runner, value configdomain.GitLabToken, scope configdomain.ConfigScope) error {
	return runner.Run("git", "config", scope.GitFlag(), configdomain.KeyGitlabToken.String(), value.String())
}

func (self *Commands) SetGiteaToken(runner gitdomain.Runner, value configdomain.GiteaToken, scope configdomain.ConfigScope) error {
	return runner.Run("git", "config", scope.GitFlag(), configdomain.KeyGiteaToken.String(), value.String())
}

func (self *Commands) SetOriginHostname(runner gitdomain.Runner, hostname configdomain.HostingOriginHostname) error {
	return runner.Run("git", "config", configdomain.KeyHostingOriginHostname.String(), hostname.String())
}

func (self *Commands) ShortenSHA(querier gitdomain.Querier, sha gitdomain.SHA) (gitdomain.SHA, error) {
	output, err := querier.QueryTrim("git", "rev-parse", "--short", sha.String())
	if err != nil {
		return gitdomain.SHA(""), fmt.Errorf(messages.BranchLocalSHAProblem, sha, err)
	}
	return gitdomain.NewSHA(output), nil
}

func (self *Commands) SquashMerge(runner gitdomain.Runner, branch gitdomain.LocalBranchName) error {
	return runner.Run("git", "merge", "--squash", "--ff", branch.String())
}

func (self *Commands) StageFiles(runner gitdomain.Runner, names ...string) error {
	args := append([]string{"add"}, names...)
	return runner.Run("git", args...)
}

func (self *Commands) Stash(runner gitdomain.Runner) error {
	err := runner.Run("git", "add", "-A")
	if err != nil {
		return err
	}
	return runner.Run("git", "stash", "-m", "Git Town WIP")
}

func (self *Commands) StashSize(querier gitdomain.Querier) (gitdomain.StashSize, error) {
	output, err := querier.QueryTrim("git", "stash", "list")
	return gitdomain.StashSize(len(stringslice.Lines(output))), err
}

func (self *Commands) UndoLastCommit(runner gitdomain.Runner) error {
	return runner.Run("git", "reset", "--soft", "HEAD~1")
}

func (self *Commands) UnstageAll(runner gitdomain.Runner) error {
	return runner.Run("git", "restore", "--staged", ".")
}

func LastBranchInRef(output string) string {
	index := strings.LastIndex(output, "/")
	return output[index+1:]
}

func NewUnmergedStage(value int) (UnmergedStage, error) {
	for _, stage := range UnmergedStages {
		if int(stage) == value {
			return stage, nil
		}
	}
	return 0, fmt.Errorf("unknown stage ID: %q", value)
}

type branchesQueryResult struct {
	BranchName     gitdomain.BranchName
	Head           bool
	RefName        string
	SHA            gitdomain.SHA
	Symref         bool
	Track          string
	UpstreamOption Option[gitdomain.RemoteBranchName] // the tracking branch name
	Worktree       bool
}

type branchesQueryResults []branchesQueryResult

func branchesQuery(querier gitdomain.Querier) (branchesQueryResults, error) {
	// WHAT DOES `:lstrip=2` DO?
	// A ref name looks like "refs/heads/branch-name" or
	// "refs/remotes/origin/branch-name". We want to remove the "refs/heads/" or
	// "refs/remotes" prefixes, so we use `:lstrip=2` to remove the first two path
	// components.
	// WHY NOT USE `:short`?
	// `:short` returns a "non-ambiguous" name. This means that if a branch and a
	// tag have the same name, it will return something like "heads/branch-name"
	// instead of "branch-name". We just want the branch name.
	forEachRefFormats := []string{
		"refname:%(refname)",                                  // full ref name
		"branchname:%(refname:lstrip=2)",                      // branch name
		"sha:%(objectname)",                                   // SHA of the commit the ref points to
		"head:%(if)%(HEAD)%(then)Y%(else)N%(end)",             // is the branch checked out in the current worktree? Y/N
		"worktree:%(if)%(worktreepath)%(then)Y%(else)N%(end)", // is the branch checked out in any worktree? Y/N
		"symref:%(if)%(symref)%(then)Y%(else)N%(end)",         // is the branch a symbolic ref? Y/N
		"upstream:%(upstream:lstrip=2)",                       // the tracking branch name
		// Leave `track` in the last position because it is the only one that contains spaces.
		// Then we can use SplitN to split the output correctly.
		"track:%(upstream:track,nobracket)", // e.g. "ahead 2", "behind 2", "ahead 2, behind 3", "gone"
	}
	output, err := querier.QueryTrim(
		"git", "for-each-ref",
		"--format="+strings.Join(forEachRefFormats, " "),
		"--sort=refname",
		"refs/heads/", "refs/remotes/", // local branches, remote branches
	)
	if err != nil {
		return branchesQueryResults{}, err
	}
	lines := stringslice.Lines(output)
	result := make(branchesQueryResults, len(lines))
	for l, line := range lines {
		parts := strings.SplitN(line, " ", len(forEachRefFormats))
		refname := strings.TrimPrefix(parts[0], "refname:")
		branchName := gitdomain.NewBranchName(strings.TrimPrefix(parts[1], "branchname:"))
		sha := gitdomain.NewSHA(strings.TrimPrefix(parts[2], "sha:"))
		head := parseYN(strings.TrimPrefix(parts[3], "head:"))
		worktree := parseYN(strings.TrimPrefix(parts[4], "worktree:"))
		symref := parseYN(strings.TrimPrefix(parts[5], "symref:"))
		upstreamOption := gitdomain.NewRemoteBranchNameOption(strings.TrimPrefix(parts[6], "upstream:")) // the tracking branch name
		track := strings.TrimPrefix(parts[7], "track:")
		result[l] = branchesQueryResult{
			BranchName:     branchName,
			Head:           head,
			RefName:        refname,
			SHA:            sha,
			Symref:         symref,
			Track:          track,
			UpstreamOption: upstreamOption,
			Worktree:       worktree,
		}
	}
	return result, nil
}

func determineSyncStatus(track string, upstream Option[gitdomain.RemoteBranchName]) gitdomain.SyncStatus {
	gone := track == "gone"
	ahead := strings.Contains(track, "ahead")
	behind := strings.Contains(track, "behind")
	switch {
	case gone:
		return gitdomain.SyncStatusDeletedAtRemote
	case ahead && behind:
		return gitdomain.SyncStatusNotInSync
	case ahead:
		return gitdomain.SyncStatusAhead
	case behind:
		return gitdomain.SyncStatusBehind
	case track == "":
		if upstream.IsSome() {
			return gitdomain.SyncStatusUpToDate
		}
		return gitdomain.SyncStatusLocalOnly
	default:
		panic(fmt.Sprintf(`unrecognized track "%s"`, track))
	}
}

// provides the path of the `.git` directory of the current repository.
func (self *Commands) gitDirectory(querier gitdomain.Querier) (string, error) {
	output, err := querier.QueryTrim("git", "rev-parse", "--absolute-git-dir")
	if err != nil {
		return "", fmt.Errorf(messages.GitDirMissing, err)
	}
	return output, nil
}

// indicates whether the ref is a local branch
func isLocalRefName(ref string) bool {
	return strings.HasPrefix(ref, "refs/heads/")
}

func makeBranchesSnapshotNewRepo(branch gitdomain.LocalBranchName) gitdomain.BranchesSnapshot {
	return gitdomain.BranchesSnapshot{
		Active: Some(branch),
		Branches: gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  Some(branch),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		},
	}
}

func parseYN(value string) bool {
	switch value {
	case "Y":
		return true
	case "N":
		return false
	default:
		panic(fmt.Sprintf("unrecognized value %q", value))
	}
}
