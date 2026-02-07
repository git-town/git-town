package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/cache"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Commands are Git commands that Git Town executes to determine which frontend commands to run.
// They don't change the user's repo, execute instantaneously, and Git Town needs to know their output.
// They are invisible to the end user unless the "verbose" option is set.
type Commands struct {
	CurrentBranchCache *cache.WithPrevious[gitdomain.LocalBranchName] // caches the currently checked out Git branch
	RemotesCache       *cache.Cache[gitdomain.Remotes]                // caches Git remotes
}

func (self *Commands) AbortMerge(runner subshelldomain.Runner) error {
	return runner.Run("git", "merge", "--abort")
}

func (self *Commands) AbortRebase(runner subshelldomain.Runner) error {
	return runner.Run("git", "rebase", "--abort")
}

// BranchAuthors provides the user accounts that contributed to the given branch.
func (self *Commands) BranchAuthors(querier subshelldomain.Querier, branch, parent gitdomain.LocalBranchName) ([]gitdomain.Author, error) {
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

func (self *Commands) BranchContainsMerges(querier subshelldomain.Querier, branch, parent gitdomain.LocalBranchName) (bool, error) {
	output, err := querier.QueryTrim("git", "log", "--merges", "--format=%H", fmt.Sprintf("%s..%s", parent, branch))
	return len(output) > 0, err
}

func (self *Commands) BranchExists(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) bool {
	err := runner.Run("git", "rev-parse", "--verify", "-q", "refs/heads/"+branch.String())
	return err == nil
}

func (self *Commands) BranchExistsAtRemote(runner subshelldomain.Runner, branch gitdomain.LocalBranchName, remote gitdomain.Remote) bool {
	err := runner.Run("git", "ls-remote", remote.String(), branch.String())
	return err == nil
}

// BranchHasUnmergedChanges indicates whether the branch with the given name
// contains changes that were not merged into the main branch.
func (self *Commands) BranchHasUnmergedChanges(querier subshelldomain.Querier, branch, parent gitdomain.LocalBranchName) (bool, error) {
	out, err := querier.QueryTrim("git", "diff", "--shortstat", parent.String(), branch.String(), "--")
	return len(out) > 0, gohacks.WrapIfError(err, messages.BranchDiffProblem, branch)
}

func (self *Commands) BranchInSyncWithParent(querier subshelldomain.Querier, branch gitdomain.LocalBranchName, parent gitdomain.BranchName) (bool, error) {
	output, err := querier.QueryTrim("git", "log", "--no-merges", "--format=%H", parent.RefName(), "^"+branch.RefName())
	return len(output) == 0, err
}

// BranchInSyncWithTracking returns whether the local branch with the given name
// contains commits that have not been pushed to its tracking branch.
func (self *Commands) BranchInSyncWithTracking(querier subshelldomain.Querier, localBranch gitdomain.LocalBranchName, trackingBranch gitdomain.RemoteBranchName) (bool, error) {
	out, err := querier.QueryTrim("git", "rev-parse", localBranch.String(), trackingBranch.String())
	if err != nil {
		return false, fmt.Errorf(messages.DiffProblem, localBranch, trackingBranch, err)
	}
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		return false, fmt.Errorf("unexpected output of git rev-parse: expected 2 lines, got %d: %s", len(lines), strings.Join(lines, ", "))
	}
	branchSHA := strings.TrimSpace(lines[0])
	trackingSHA := strings.TrimSpace(lines[1])
	return branchSHA == trackingSHA, nil
}

func (self *Commands) BranchesAvailableInCurrentWorktree(querier subshelldomain.Querier) (gitdomain.LocalBranchNames, error) {
	branches, err := branchesQuery(querier)
	if err != nil {
		return gitdomain.LocalBranchNames{}, err
	}
	result := gitdomain.LocalBranchNames{}
	for _, branch := range branches {
		// Skip symbolic refs
		if branch.Symref {
			continue
		}
		// Only include local branches
		if !isLocalRefName(branch.RefName) {
			continue
		}
		// Include branches that are not checked out in other worktrees
		// A branch is available in the current worktree if:
		// - It's not checked out anywhere (!branch.Worktree), OR
		// - It's checked out in the current worktree (branch.Head)
		if !branch.Worktree || branch.Head {
			result = append(result, branch.BranchName.LocalName())
		}
	}
	return result, nil
}

func (self *Commands) BranchesSnapshot(querier subshelldomain.Querier) (gitdomain.BranchesSnapshot, error) {
	branches, err := branchesQuery(querier)
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), err
	}
	if len(branches) == 0 {
		// We are in a brand-new repo.
		// Report the initial branch name (reported by `git branch --show-current`) as the current branch.
		currentBranchOpt, err := self.CurrentBranchUncached(querier)
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), err
		}
		currentBranch, hasCurrentBranch := currentBranchOpt.Get()
		if !hasCurrentBranch {
			return gitdomain.EmptyBranchesSnapshot(), nil
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
				Local:      Some(gitdomain.BranchData{Name: branch.BranchName.LocalName(), SHA: branch.SHA}),
				RemoteName: branch.UpstreamOption,
				RemoteSHA:  None[gitdomain.SHA](), // may be added later
				SyncStatus: gitdomain.SyncStatusOtherWorktree,
			})
		case isLocalRefName(branch.RefName):
			syncStatus := determineSyncStatus(branch.Track, branch.UpstreamOption)
			result = append(result, gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: branch.BranchName.LocalName(), SHA: branch.SHA}),
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
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(remoteBranchName),
					RemoteSHA:  Some(branch.SHA),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				})
			}
		}
	}
	detachedHead := false
	if currentBranchOpt.IsNone() {
		rebaseInProgress, err := self.HasRebaseInProgress(querier)
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), err
		}
		if !rebaseInProgress {
			// We are in a detached HEAD state. Use the current HEAD location as the branch name.
			detachedHead = true
			headSHA, err := self.CurrentSHA(querier)
			if err != nil {
				return gitdomain.EmptyBranchesSnapshot(), err
			}
			currentBranchOpt = gitdomain.NewLocalBranchNameOption(headSHA.String())
			// prepend to result
			result = slices.Insert(result, 0, gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: gitdomain.LocalBranchName(headSHA.String()), SHA: headSHA}),
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
		Branches:     result,
		Active:       currentBranchOpt,
		DetachedHead: detachedHead,
	}, nil
}

func (self *Commands) ChangeDir(dir gitdomain.RepoRootDir) error {
	return os.Chdir(dir.String())
}

func (self *Commands) CheckoutBranch(runner subshelldomain.Runner, name gitdomain.LocalBranchName, merge configdomain.SwitchUsingMerge) error {
	currentBranch, hasCurrentBranch := self.CurrentBranchCache.Get()
	if hasCurrentBranch && currentBranch == name {
		return nil
	}
	return self.CheckoutBranchUncached(runner, name, merge)
}

func (self *Commands) CheckoutBranchUncached(runner subshelldomain.Runner, name gitdomain.LocalBranchName, merge configdomain.SwitchUsingMerge) error {
	args := []string{"checkout", name.String()}
	if merge {
		args = append(args, "-m")
	}
	if err := runner.Run("git", args...); err != nil {
		return fmt.Errorf(messages.BranchCheckoutProblem, name, err)
	}
	if name.String() != "-" {
		self.CurrentBranchCache.Set(name)
	} else {
		self.CurrentBranchCache.Invalidate()
	}
	return nil
}

func (self *Commands) CherryPick(runner subshelldomain.Runner, sha gitdomain.SHA) error {
	return runner.Run("git", "cherry-pick", sha.String())
}

func (self *Commands) CherryPickAbort(runner subshelldomain.Runner) error {
	return runner.Run("git", "cherry-pick", "--abort")
}

func (self *Commands) CherryPickContinue(runner subshelldomain.Runner) error {
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
	commentOutSquashCommitMessageOnce.Do(func() {
		commentOutSquashCommitMessageRegex = regexp.MustCompile("(?m)^")
	})
	content = commentOutSquashCommitMessageRegex.ReplaceAllString(content, "# ")
	return os.WriteFile(squashMessageFile, []byte(content), 0o600)
}

var (
	commentOutSquashCommitMessageOnce  sync.Once
	commentOutSquashCommitMessageRegex *regexp.Regexp
)

func (self *Commands) Commit(runner subshelldomain.Runner, useMessage configdomain.UseMessage, author Option[gitdomain.Author], commitHook configdomain.CommitHook) error {
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

func (self *Commands) CommitMessage(querier subshelldomain.Querier, sha gitdomain.SHA) (gitdomain.CommitMessage, error) {
	output, err := querier.QueryTrim("git", "show", "--no-patch", "--format=%B", sha.String())
	return gitdomain.CommitMessage(strings.TrimSpace(output)), err
}

func (self *Commands) CommitStart(runner subshelldomain.Runner) error {
	return runner.Run("git", "commit")
}

func (self *Commands) CommitsInBranch(querier subshelldomain.Querier, branch gitdomain.LocalBranchName, parent Option[gitdomain.LocalBranchName]) (gitdomain.Commits, error) {
	if parent, hasParent := parent.Get(); hasParent {
		return self.CommitsInFeatureBranch(querier, branch, parent.BranchName())
	}
	return self.CommitsInPerennialBranch(querier)
}

func (self *Commands) CommitsInFeatureBranch(querier subshelldomain.Querier, branch gitdomain.LocalBranchName, parent gitdomain.BranchName) (gitdomain.Commits, error) {
	output, err := querier.QueryTrim("git", "log", "--format=%H %s", fmt.Sprintf("%s..%s", parent.String(), branch.String()))
	if err != nil {
		return gitdomain.Commits{}, err
	}
	lines := stringslice.NonEmptyLines(output)
	result := make(gitdomain.Commits, 0, len(lines))
	for _, line := range lines {
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

func (self *Commands) CommitsInPerennialBranch(querier subshelldomain.Querier) (gitdomain.Commits, error) {
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

func (self *Commands) ContentBlobInfo(querier subshelldomain.Querier, branch gitdomain.Location, filePath string) (Option[Blob], error) {
	output, err := querier.QueryTrim("git", "ls-tree", branch.String(), filePath)
	if err != nil || len(output) == 0 {
		return None[Blob](), err
	}
	sha, err := ParseLsTreeOutput(output)
	return Some(sha), err
}

func (self *Commands) ContinueRebase(runner subshelldomain.Runner) error {
	return runner.RunWithEnv([]string{"GIT_EDITOR=true"}, "git", "rebase", "--continue")
}

// CreateAndCheckoutBranch creates a new branch with the given name and checks it out using a single Git operation.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *Commands) CreateAndCheckoutBranch(runner subshelldomain.Runner, name gitdomain.LocalBranchName) error {
	err := runner.Run("git", "checkout", "-b", name.String())
	if err == nil {
		self.CurrentBranchCache.Set(name)
	}
	return err
}

// CreateAndCheckoutBranchWithParent creates a new branch with the given name and parent and checks it out using a single Git operation.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *Commands) CreateAndCheckoutBranchWithParent(runner subshelldomain.Runner, name gitdomain.LocalBranchName, parent gitdomain.Location) error {
	args := []string{"checkout", "-b", name.String(), parent.String()}
	if parent.IsRemoteBranchName() {
		args = append(args, "--no-track")
	}
	err := runner.Run("git", args...)
	if err == nil {
		self.CurrentBranchCache.Set(name)
	}
	return err
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *Commands) CreateBranch(runner subshelldomain.Runner, name gitdomain.LocalBranchName, parent gitdomain.Location) error {
	return runner.Run("git", "branch", name.String(), parent.String())
}

func (self *Commands) CreateTrackingBranch(runner subshelldomain.Runner, branch gitdomain.LocalBranchName, remote gitdomain.Remote, pushHook configdomain.PushHook) error {
	args := []string{"push"}
	if !pushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, "-u", remote.String())
	args = append(args, branch.String())
	return runner.Run("git", args...)
}

// CurrentBranch provides the name of the current branch.
// Provides (None, nil) if there is no current branch, e.g. when the Git HEAD is detached.
func (self *Commands) CurrentBranch(querier subshelldomain.Querier) (Option[gitdomain.LocalBranchName], error) {
	if cachedCurrentBranch, hasCachedCurrentBranch := self.CurrentBranchCache.Get(); hasCachedCurrentBranch {
		return Some(cachedCurrentBranch), nil
	}
	currentBranchOpt, err := self.CurrentBranchUncached(querier)
	if currentBranch, hasCurrentBranch := currentBranchOpt.Get(); hasCurrentBranch {
		self.CurrentBranchCache.Set(currentBranch)
	}
	return currentBranchOpt, err
}

func (self *Commands) CurrentBranchDuringRebase(querier subshelldomain.Querier) (Option[gitdomain.LocalBranchName], error) {
	gitDir, err := self.gitDirectory(querier)
	if err != nil {
		return None[gitdomain.LocalBranchName](), err
	}
	for _, rebaseHeadFileName := range []string{"rebase-merge/head-name", "rebase-apply/head-name"} {
		rebaseHeadFilePath := filepath.Join(gitDir, rebaseHeadFileName)
		content, err := os.ReadFile(rebaseHeadFilePath)
		if err != nil {
			continue
		}
		refName := strings.TrimSpace(string(content))
		if branchName, isBranchName := strings.CutPrefix(refName, "refs/heads/"); isBranchName {
			return Some(gitdomain.NewLocalBranchName(branchName)), nil
		}
		// rebase head name is not a branch name
		break
	}
	return None[gitdomain.LocalBranchName](), nil
}

func (self *Commands) CurrentBranchHasTrackingBranch(runner subshelldomain.Runner) bool {
	err := runner.Run("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	return err == nil
}

func (self *Commands) CurrentBranchUncached(querier subshelldomain.Querier) (Option[gitdomain.LocalBranchName], error) {
	// first try to detect the current branch the normal way
	output, err := querier.QueryTrim("git", "branch", "--show-current")
	if err != nil {
		return None[gitdomain.LocalBranchName](), fmt.Errorf(messages.BranchCurrentProblem, err)
	}
	if output != "" {
		return Some(gitdomain.NewLocalBranchName(output)), nil
	}
	// here we couldn't detect the current branch the normal way --> assume we are in a rebase and try the rebase way
	return self.CurrentBranchDuringRebase(querier)
}

// CurrentSHA provides the SHA of the currently checked out branch/commit.
func (self *Commands) CurrentSHA(querier subshelldomain.Querier) (gitdomain.SHA, error) {
	return self.SHAForBranch(querier, "HEAD")
}

// DeleteLastCommit resets HEAD to the previous commit.
func (self *Commands) DeleteLastCommit(runner subshelldomain.Runner) error {
	return runner.Run("git", "reset", "--hard", "HEAD~1")
}

func (self *Commands) DeleteLocalBranch(runner subshelldomain.Runner, name gitdomain.LocalBranchName) error {
	return runner.Run("git", "branch", "-D", name.String())
}

func (self *Commands) DeleteTrackingBranch(runner subshelldomain.Runner, name gitdomain.RemoteBranchName) error {
	remote, localBranchName := name.Parts()
	return runner.Run("git", "push", remote.String(), ":"+localBranchName.String())
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (self *Commands) DiffParent(runner subshelldomain.Runner, branch, parentBranch gitdomain.LocalBranchName, nameOnly configdomain.NameOnly) error {
	args := []string{"diff", "--merge-base", parentBranch.String(), branch.String()}
	if nameOnly {
		args = append(args, "--name-only")
	}
	return runner.Run("git", args...)
}

func (self *Commands) DiscardOpenChanges(runner subshelldomain.Runner) error {
	return runner.Run("git", "reset", "--hard")
}

func (self *Commands) DropMostRecentStash(runner subshelldomain.Runner) error {
	return runner.Run("git", "stash", "drop")
}

func (self *Commands) Fetch(runner subshelldomain.Runner, syncTags configdomain.SyncTags) error {
	if syncTags.ShouldSyncTags() {
		return runner.Run("git", "fetch", "--prune", "--tags")
	}
	return runner.Run("git", "fetch", "--prune", "--no-tags")
}

func (self *Commands) FetchUpstream(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) error {
	return runner.Run("git", "fetch", gitdomain.RemoteUpstream.String(), branch.String())
}

func (self *Commands) FileConflicts(querier subshelldomain.Querier) (FileConflicts, error) {
	output, err := querier.Query("git", "ls-files", "--unmerged")
	if err != nil {
		return FileConflicts{}, err
	}
	return ParseLsFilesUnmergedOutput(output)
}

// FirstCommitMessageInBranch provides the commit message of the first commit in the branch with the given name.
func (self *Commands) FirstCommitMessageInBranch(runner subshelldomain.Querier, branch, parent gitdomain.BranchName) (Option[gitdomain.CommitMessage], error) {
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

// FirstExistingBranch provides the first branch in the given list of branch names that actually exists.
func (self *Commands) FirstExistingBranch(runner subshelldomain.Runner, branches ...gitdomain.LocalBranchName) Option[gitdomain.LocalBranchName] {
	for _, branch := range branches {
		if self.BranchExists(runner, branch) {
			return Some(branch)
		}
	}
	return None[gitdomain.LocalBranchName]()
}

func (self *Commands) ForcePushBranchSafely(runner subshelldomain.Runner, pushHook configdomain.PushHook, forceIfIncludes bool) error {
	args := []string{"push", "--force-with-lease"}
	if forceIfIncludes {
		args = append(args, "--force-if-includes")
	}
	if !pushHook {
		args = append(args, "--no-verify")
	}
	return runner.Run("git", args...)
}

func (self *Commands) GitVersion(querier subshelldomain.Querier) (Version, error) {
	output, err := querier.QueryTrim("git", "version")
	if err != nil {
		return EmptyVersion(), fmt.Errorf(messages.GitVersionProblem, err)
	}
	gitVersionOnce.Do(func() {
		gitVersionRegex = regexp.MustCompile(`git version (\d+).(\d+).(\w+)`)
	})
	matches := gitVersionRegex.FindStringSubmatch(output)
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

var (
	gitVersionOnce  sync.Once
	gitVersionRegex *regexp.Regexp
)

func (self *Commands) HasMergeInProgress(runner subshelldomain.Runner) bool {
	err := runner.Run("git", "rev-parse", "--verify", "-q", "MERGE_HEAD")
	return err == nil
}

func (self *Commands) HasRebaseInProgress(querier subshelldomain.Querier) (bool, error) {
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

func (self *Commands) MergeBranchNoEdit(runner subshelldomain.Runner, branch gitdomain.BranchName) error {
	return runner.Run("git", "merge", "--no-edit", "--ff", branch.String())
}

// MergeConflicts loads the information needed to determine which of the given file conflicts are phantom merge conflicts.
func (self *Commands) MergeConflicts(querier subshelldomain.Querier, fileConflicts FileConflicts, parentLocation gitdomain.Location, rootBranch gitdomain.LocalBranchName) (MergeConflicts, error) {
	result := make(MergeConflicts, len(fileConflicts))
	for f, fileConflict := range fileConflicts {
		rootBlob := None[Blob]()
		parentBlob := None[Blob]()
		if currentBranchBlob, has := fileConflict.CurrentBranchChange.Get(); has {
			var err error
			rootBlob, err = self.ContentBlobInfo(querier, rootBranch.Location(), currentBranchBlob.FilePath)
			if err != nil {
				return result, err
			}
			parentBlob, err = self.ContentBlobInfo(querier, parentLocation, currentBranchBlob.FilePath)
			if err != nil {
				return result, err
			}
		}
		result[f] = MergeConflict{
			Current: fileConflict.CurrentBranchChange,
			Parent:  parentBlob,
			Root:    rootBlob,
		}
	}
	return result, nil
}

func (self *Commands) MergeFastForward(runner subshelldomain.Runner, branch gitdomain.BranchName) error {
	return runner.Run("git", "merge", "--ff-only", branch.String())
}

func (self *Commands) MergeNoFastForward(runner subshelldomain.Runner, useMessage configdomain.UseMessage, branch gitdomain.LocalBranchName) error {
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

func (self *Commands) OriginHead(querier subshelldomain.Querier) Option[gitdomain.LocalBranchName] {
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

func (self *Commands) PopStash(runner subshelldomain.Runner) error {
	err := runner.Run("git", "stash", "pop")
	if err != nil {
		_ = runner.Run("git", "stash", "drop")
	}
	return err
}

// PreviouslyCheckedOutBranch provides the name of the branch that was checked out before the current branch was checked out.
func (self *Commands) PreviouslyCheckedOutBranch(querier subshelldomain.Querier) Option[gitdomain.LocalBranchName] {
	output, err := querier.QueryTrim("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return None[gitdomain.LocalBranchName]()
	}
	if output == "" {
		return None[gitdomain.LocalBranchName]()
	}
	return gitdomain.NewLocalBranchNameOption(output)
}

func (self *Commands) Pull(runner subshelldomain.Runner) error {
	return runner.Run("git", "pull")
}

// PushCurrentBranch pushes the current branch to its tracking branch.
func (self *Commands) PushCurrentBranch(runner subshelldomain.Runner, pushHook configdomain.PushHook) error {
	args := []string{"push"}
	if !pushHook {
		args = append(args, "--no-verify")
	}
	return runner.Run("git", args...)
}

func (self *Commands) PushLocalBranch(runner subshelldomain.Runner, localSHA gitdomain.SHA, branch gitdomain.LocalBranchName, remote gitdomain.Remote, pushHook configdomain.PushHook) error {
	args := []string{"push"}
	if !pushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, remote.String(), localSHA.String()+":refs/heads/"+branch.String())
	return runner.Run("git", args...)
}

// PushTags pushes new the Git tags to origin.
func (self *Commands) PushTags(runner subshelldomain.Runner, pushHook configdomain.PushHook) error {
	args := []string{"push", "--tags"}
	if !pushHook {
		args = append(args, "--no-verify")
	}
	return runner.Run("git", args...)
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (self *Commands) Rebase(runner subshelldomain.Runner, target gitdomain.BranchName) error {
	return runner.Run("git", "-c", "rebase.updateRefs=false", "rebase", target.String())
}

// RebaseOnto initiates a Git rebase of the current branch onto the given branch.
func (self *Commands) RebaseOnto(runner subshelldomain.Runner, branchToRebaseOnto gitdomain.Location, commitsToRemove gitdomain.Location) error {
	return runner.Run("git", "-c", "rebase.updateRefs=false", "rebase", "--onto", branchToRebaseOnto.String(), commitsToRemove.String())
}

func (self *Commands) Remotes(querier subshelldomain.Querier) (gitdomain.Remotes, error) {
	if cachedRemotes, hasCachedRemotes := self.RemotesCache.Get(); hasCachedRemotes {
		return cachedRemotes, nil
	}
	remotes, err := self.RemotesUncached(querier)
	if err == nil {
		self.RemotesCache.Set(remotes)
	}
	return remotes, err
}

func (self *Commands) RemotesUncached(querier subshelldomain.Querier) (gitdomain.Remotes, error) {
	out, err := querier.QueryTrim("git", "remote")
	if err != nil {
		return gitdomain.Remotes{}, fmt.Errorf(messages.RemotesProblem, err)
	}
	if out == "" {
		return gitdomain.Remotes{}, nil
	}
	return gitdomain.NewRemotes(stringslice.Lines(out)...), nil
}

// RemoveCommit removes the given commit from the current branch
func (self *Commands) RemoveCommit(runner subshelldomain.Runner, commit gitdomain.SHA) error {
	return runner.Run("git", "-c", "rebase.updateRefs=false", "rebase", "--onto", commit.String()+"^", commit.String())
}

func (self *Commands) RemoveFile(runner subshelldomain.Runner, fileName string) error {
	return runner.Run("git", "rm", fileName)
}

func (self *Commands) RenameBranch(runner subshelldomain.Runner, oldName, newName gitdomain.LocalBranchName) error {
	return runner.Run("git", "branch", "--move", oldName.String(), newName.String())
}

func (self *Commands) RepoStatus(backend subshelldomain.RunnerQuerier) (gitdomain.RepoStatus, error) {
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

func (self *Commands) ResetBranch(runner subshelldomain.Runner, target gitdomain.BranchName) error {
	return runner.Run("git", "reset", "--soft", target.String(), "--")
}

func (self *Commands) ResetCurrentBranchToSHA(runner subshelldomain.Runner, sha gitdomain.SHA) error {
	return runner.Run("git", "reset", "--hard", sha.String())
}

func (self *Commands) ResetRemoteBranchToSHA(runner subshelldomain.Runner, branch gitdomain.RemoteBranchName, sha gitdomain.SHA) error {
	remote := branch.Remote()
	return runner.Run("git", "push", "--force-with-lease", remote.String(), sha.String()+":"+branch.LocalBranchName().String())
}

func (self *Commands) ResolveConflict(runner subshelldomain.Runner, file string, resolution gitdomain.ConflictResolution) error {
	return runner.Run("git", "checkout", resolution.GitFlag(), file)
}

func (self *Commands) RevertCommit(runner subshelldomain.Runner, sha gitdomain.SHA) error {
	return runner.Run("git", "revert", sha.String())
}

// RootDirectory provides the path of the root directory of the current repository.
func (self *Commands) RootDirectory(querier subshelldomain.Querier) Option[gitdomain.RepoRootDir] {
	output, err := querier.QueryTrim("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return None[gitdomain.RepoRootDir]()
	}
	return Some(gitdomain.NewRepoRootDir(filepath.FromSlash(output)))
}

func (self *Commands) SHAForBranch(querier subshelldomain.Querier, name gitdomain.BranchName) (gitdomain.SHA, error) {
	output, err := querier.QueryTrim("git", "rev-parse", name.String())
	return gitdomain.NewSHA(output), gohacks.WrapIfError(err, messages.BranchLocalSHAProblem, name)
}

func (self *Commands) ShortenSHA(querier subshelldomain.Querier, sha gitdomain.SHA) (gitdomain.SHA, error) {
	output, err := querier.QueryTrim("git", "rev-parse", "--short", sha.String())
	return gitdomain.NewSHA(output), gohacks.WrapIfError(err, messages.BranchLocalSHAProblem, sha)
}

func (self *Commands) SquashMerge(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) error {
	return runner.Run("git", "merge", "--squash", "--ff", branch.String())
}

func (self *Commands) StageFiles(runner subshelldomain.Runner, names ...string) error {
	args := append([]string{"add"}, names...)
	return runner.Run("git", args...)
}

// StandardBranch determines the branch that is configured in Git as the default branch.
func (self *Commands) StandardBranch(querier subshelldomain.Querier) Option[gitdomain.LocalBranchName] {
	if defaultBranch, has := gitconfig.DefaultBranch(querier).Get(); has {
		return Some(defaultBranch)
	}
	return self.OriginHead(querier)
}

func (self *Commands) Stash(runner subshelldomain.Runner) error {
	if err := runner.Run("git", "add", "-A"); err != nil {
		return err
	}
	return runner.Run("git", "stash", "-m", "Git Town WIP")
}

func (self *Commands) StashSize(querier subshelldomain.Querier) (gitdomain.StashSize, error) {
	output, err := querier.QueryTrim("git", "stash", "list")
	return gitdomain.StashSize(len(stringslice.Lines(output))), err
}

func (self *Commands) UndoLastCommit(runner subshelldomain.Runner) error {
	return runner.Run("git", "reset", "--soft", "HEAD~1")
}

func (self *Commands) UnstageAll(runner subshelldomain.Runner) error {
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

func branchesQuery(querier subshelldomain.Querier) (branchesQueryResults, error) {
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
func (self *Commands) gitDirectory(querier subshelldomain.Querier) (string, error) {
	output, err := querier.QueryTrim("git", "rev-parse", "--absolute-git-dir")
	return output, gohacks.WrapIfError(err, messages.GitDirMissing)
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
				Local:      Some(gitdomain.BranchData{Name: branch, SHA: "0000000"}), // brand-new repos witout any commits don't have a SHA
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		},
		DetachedHead: false,
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
