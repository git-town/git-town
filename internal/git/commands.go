package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/gohacks/cache"
	"github.com/git-town/git-town/v19/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
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
	output, err := querier.QueryTrim("git", "log", "--merges", fmt.Sprintf("%s..%s", parent, branch))
	return len(output) > 0, err
}

func (self *Commands) BranchExists(runner gitdomain.Runner, branch gitdomain.LocalBranchName) bool {
	err := runner.Run("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch.String())
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

func (self *Commands) BranchInSyncWithParent(querier gitdomain.Querier, branch, parent gitdomain.LocalBranchName) (bool, error) {
	output, err := querier.QueryTrim("git", "log", "--no-merges", parent.String(), "^"+branch.String())
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
	output, err := querier.Query("git", "-c", "core.abbrev=40", "branch", "-vva", "--sort=refname")
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), err
	}
	if output == "" {
		// We are in a brand-new repo.
		// Even though `git branch` returns nothing, `git branch --show-current` will return the initial branch name.
		currentBranch, err := self.CurrentBranchUncached(querier)
		if err != nil {
			return gitdomain.EmptyBranchesSnapshot(), err
		}
		return gitdomain.BranchesSnapshot{
			Active: Some(currentBranch),
			Branches: gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(currentBranch),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			},
		}, nil
	}
	branches, currentBranchOpt := ParseVerboseBranchesOutput(output)
	currentBranch, hasCurrentBranch := currentBranchOpt.Get()
	if hasCurrentBranch {
		self.CurrentBranchCache.Set(currentBranch)
	}
	return gitdomain.BranchesSnapshot{
		Branches: branches,
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

func (self *Commands) CheckoutOurVersion(runner gitdomain.Runner, file string) error {
	return runner.Run("git", "checkout", "--ours", file)
}

func (self *Commands) CheckoutTheirVersion(runner gitdomain.Runner, file string) error {
	return runner.Run("git", "checkout", "--theirs", file)
}

func (self *Commands) CherryPick(runner gitdomain.Runner, sha gitdomain.SHA) error {
	return runner.Run("git", "cherry-pick", sha.String())
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

func (self *Commands) Commit(runner gitdomain.Runner, useMessage configdomain.UseMessage, author Option[gitdomain.Author]) error {
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
	return runner.Run("git", args...)
}

func (self *Commands) CommitMessage(querier gitdomain.Querier, sha gitdomain.SHA) (gitdomain.CommitMessage, error) {
	output, err := querier.QueryTrim("git", "log", "--format=%B", "-n", "1", sha.String())
	return gitdomain.CommitMessage(strings.TrimSpace(output)), err
}

func (self *Commands) CommitStart(runner gitdomain.Runner) error {
	return runner.Run("git", "commit")
}

func (self *Commands) CommitsInBranch(querier gitdomain.Querier, branch gitdomain.LocalBranchName, parent Option[gitdomain.LocalBranchName]) (gitdomain.Commits, error) {
	if parent, hasParent := parent.Get(); hasParent {
		return self.CommitsInFeatureBranch(querier, branch, parent)
	}
	return self.CommitsInPerennialBranch(querier)
}

func (self *Commands) CommitsInFeatureBranch(querier gitdomain.Querier, branch, parent gitdomain.LocalBranchName) (gitdomain.Commits, error) {
	output, err := querier.QueryTrim("git", "cherry", "-v", parent.String(), branch.String())
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

func (self *Commands) CommitsInPerennialBranch(querier gitdomain.Querier) (gitdomain.Commits, error) {
	output, err := querier.QueryTrim("git", "log", "--pretty=format:%H %s", "-10")
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
	return runner.Run("git", "-c", "core.editor=true", "rebase", "--continue")
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

func (self *Commands) CurrentBranchHasTrackingBranch(runner gitdomain.Runner) bool {
	err := runner.Run("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	return err == nil
}

func (self *Commands) CurrentBranchUncached(querier gitdomain.Querier) (gitdomain.LocalBranchName, error) {
	// first try to detect the current branch the normal way
	output, err := querier.QueryTrim("git", "branch", "--show-current")
	if err == nil && output != "" {
		return gitdomain.NewLocalBranchName(output), nil
	}
	// here we couldn't detect the current branch the normal way --> assume we are in a rebase and try the rebase way
	return self.currentBranchDuringRebase(querier)
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

func (self *Commands) HasLocalBranch(runner gitdomain.Runner, name gitdomain.LocalBranchName) bool {
	return runner.Run("git", "show-ref", "--quiet", "refs/heads/"+name.String()) == nil
}

func (self *Commands) HasMergeInProgress(querier gitdomain.Querier) bool {
	_, err := querier.Query("git", "rev-parse", "-q", "--verify", "MERGE_HEAD")
	return err == nil
}

func (self *Commands) HasRebaseInProgress(querier gitdomain.Querier) bool {
	gitDir, err := querier.QueryTrim("git", "rev-parse", "--absolute-git-dir")
	if err != nil {
		return false
	}
	for _, rebaseDirName := range []string{"rebase-merge", "rebase-apply"} {
		rebaseDirPath := filepath.Join(gitDir, rebaseDirName)
		stat, err := os.Stat(rebaseDirPath)
		if err == nil && stat.IsDir() {
			return true
		}
	}
	return false
}

// HeadCommitMessage provides the commit message for the last commit.
func (self *Commands) HeadCommitMessage(querier gitdomain.Querier) (gitdomain.CommitMessage, error) {
	out, err := querier.QueryTrim("git", "log", "-1", "--format=%B")
	if err != nil {
		return "", fmt.Errorf(messages.CommitMessageProblem, err)
	}
	return gitdomain.CommitMessage(out), nil
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
	return Some(gitdomain.NewLocalBranchName(output))
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
func (self *Commands) Rebase(runner gitdomain.Runner, target gitdomain.BranchName, version Version) error {
	args := []string{"rebase", target.String()}
	if version.HasRebaseUpdateRefs() {
		args = append(args, "--no-update-refs")
	}
	return runner.Run("git", args...)
}

// Rebase initiates a Git rebase of the current branch onto the given branch.
func (self *Commands) RebaseOnto(runner gitdomain.Runner, branchToRebaseOnto gitdomain.Location, commitsToRemove gitdomain.Location, upstream Option[gitdomain.LocalBranchName]) error {
	args := []string{"rebase", "--onto", branchToRebaseOnto.String()}
	if upstream, hasUpstream := upstream.Get(); hasUpstream {
		args = append(args, upstream.String())
	}
	args = append(args, commitsToRemove.String(), "--no-update-refs")
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
	return runner.Run("git", "rebase", "--onto", commit.String()+"^", commit.String(), "--no-update-refs")
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

func (self *Commands) RepoStatus(querier gitdomain.Querier) (gitdomain.RepoStatus, error) {
	output, err := querier.Query("git", "status", "-z", "--ignore-submodules")
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
	mergeInProgress := self.HasMergeInProgress(querier)
	rebaseInProgress := self.HasRebaseInProgress(querier)
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

func (self *Commands) ResetCurrentBranchToSHA(runner gitdomain.Runner, sha gitdomain.SHA, hard bool) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	args = append(args, sha.String())
	return runner.Run("git", args...)
}

func (self *Commands) ResetRemoteBranchToSHA(runner gitdomain.Runner, branch gitdomain.RemoteBranchName, sha gitdomain.SHA) error {
	remote := branch.Remote()
	return runner.Run("git", "push", "--force-with-lease", remote.String(), sha.String()+":"+branch.LocalBranchName().String())
}

func (self *Commands) RevertCommit(runner gitdomain.Runner, sha gitdomain.SHA) error {
	return runner.Run("git", "revert", sha.String())
}

// RootDirectory provides the path of the root directory of the current repository,
// i.e. the directory that contains the ".git" folder.
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

func (self *Commands) SetBitbucketAppPassword(runner gitdomain.Runner, value configdomain.BitbucketAppPassword) error {
	return runner.Run("git", "config", configdomain.KeyBitbucketAppPassword.String(), value.String())
}

func (self *Commands) SetBitbucketUsername(runner gitdomain.Runner, value configdomain.BitbucketUsername) error {
	return runner.Run("git", "config", configdomain.KeyBitbucketUsername.String(), value.String())
}

func (self *Commands) SetCodebergToken(runner gitdomain.Runner, value configdomain.CodebergToken) error {
	return runner.Run("git", "config", configdomain.KeyCodebergToken.String(), value.String())
}

func (self *Commands) SetForgeType(runner gitdomain.Runner, platform configdomain.ForgeType) error {
	return runner.Run("git", "config", configdomain.KeyForgeType.String(), platform.String())
}

func (self *Commands) SetGitAlias(runner gitdomain.Runner, aliasableCommand configdomain.AliasableCommand) error {
	return runner.Run("git", "config", "--global", aliasableCommand.Key().String(), "town "+aliasableCommand.String())
}

func (self *Commands) SetGitHubToken(runner gitdomain.Runner, value configdomain.GitHubToken) error {
	return runner.Run("git", "config", configdomain.KeyGithubToken.String(), value.String())
}

func (self *Commands) SetGitLabToken(runner gitdomain.Runner, value configdomain.GitLabToken) error {
	return runner.Run("git", "config", configdomain.KeyGitlabToken.String(), value.String())
}

func (self *Commands) SetGiteaToken(runner gitdomain.Runner, value configdomain.GiteaToken) error {
	return runner.Run("git", "config", configdomain.KeyGiteaToken.String(), value.String())
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

func (self *Commands) currentBranchDuringRebase(querier gitdomain.Querier) (gitdomain.LocalBranchName, error) {
	output, err := querier.QueryTrim("git", "branch", "--list")
	if err != nil {
		return "", err
	}
	lines := stringslice.Lines(output)
	linesWithStar := stringslice.LinesWithPrefix(lines, "* ")
	if len(linesWithStar) == 0 {
		return "", err
	}
	if len(linesWithStar) > 1 {
		panic("multiple lines with star found:\n " + output)
	}
	lineWithStar := linesWithStar[0]
	return ParseActiveBranchDuringRebase(lineWithStar), nil
}

func IsAhead(branchName, remoteText string) (bool, Option[gitdomain.RemoteBranchName]) {
	reText := fmt.Sprintf(`\[(\w+\/%s): ahead \d+\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, Some(gitdomain.NewRemoteBranchName(matches[1]))
	}
	return false, None[gitdomain.RemoteBranchName]()
}

func IsAheadAndBehind(branchName, remoteText string) (bool, Option[gitdomain.RemoteBranchName]) {
	reText := fmt.Sprintf(`\[(\w+\/%s): ahead \d+, behind \d+\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, Some(gitdomain.NewRemoteBranchName(matches[1]))
	}
	return false, None[gitdomain.RemoteBranchName]()
}

func IsBehind(branchName, remoteText string) (bool, Option[gitdomain.RemoteBranchName]) {
	reText := fmt.Sprintf(`\[(\w+\/%s): behind \d+\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, Some(gitdomain.NewRemoteBranchName(matches[1]))
	}
	return false, None[gitdomain.RemoteBranchName]()
}

func IsInSync(branchName, remoteText string) (bool, Option[gitdomain.RemoteBranchName]) {
	reText := fmt.Sprintf(`\[(\w+\/%s)\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, Some(gitdomain.NewRemoteBranchName(matches[1]))
	}
	return false, None[gitdomain.RemoteBranchName]()
}

// IsRemoteGone indicates whether the given part of "git branch -vva" indicates a deleted tracking branch.
func IsRemoteGone(branchName, remoteText string) (bool, Option[gitdomain.RemoteBranchName]) {
	reText := fmt.Sprintf(`^\[(\w+\/%s): gone\] `, regexp.QuoteMeta(branchName))
	re := regexp.MustCompile(reText)
	matches := re.FindStringSubmatch(remoteText)
	if len(matches) == 2 {
		return true, Some(gitdomain.NewRemoteBranchName(matches[1]))
	}
	return false, None[gitdomain.RemoteBranchName]()
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

func ParseActiveBranchDuringRebase(lineWithStar string) gitdomain.LocalBranchName {
	parts := strings.Split(lineWithStar, " ")
	partsWithBranchName := parts[4:]
	branchNameWithClosingParen := strings.Join(partsWithBranchName, " ")
	return gitdomain.NewLocalBranchName(branchNameWithClosingParen[:len(branchNameWithClosingParen)-1])
}

func ParseVerboseBranchesOutput(output string) (gitdomain.BranchInfos, Option[gitdomain.LocalBranchName]) {
	result := gitdomain.BranchInfos{}
	spaceRE := regexp.MustCompile(" +")
	lines := stringslice.Lines(output)
	checkedoutBranch := None[gitdomain.LocalBranchName]()
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
		if parts[0] == "(no" || parts[1] == "(no" || parts[1] == "branch," { // "(no" as in "(no branch, rebasing main)" is what we get when a rebase is active, in which case no branch is checked out
			continue
		}
		var branchName string
		var sha Option[gitdomain.SHA]
		if parts[1] == "detached" {
			parts := spaceRE.Split(line[2:], 6)
			branchName = parts[4]
			sha = Some(gitdomain.NewSHA(parts[4]))
		} else {
			branchName = parts[0]
			rawSHA := parts[1]
			if rawSHA == "->" {
				continue
			}
			sha = Some(gitdomain.NewSHA(rawSHA))
		}
		remoteText := parts[2]
		if line[0] == '*' {
			checkedoutBranch = Some(gitdomain.NewLocalBranchName(branchName))
		}
		syncStatus, trackingBranchName := determineSyncStatus(branchName, remoteText)
		switch {
		case line[0] == '+':
			result = append(result, gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName(branchName)),
				LocalSHA:   sha,
				SyncStatus: gitdomain.SyncStatusOtherWorktree,
				RemoteName: trackingBranchName,
				RemoteSHA:  None[gitdomain.SHA](),
			})
		case isLocalBranchName(branchName):
			result = append(result, gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName(branchName)),
				LocalSHA:   sha,
				SyncStatus: syncStatus,
				RemoteName: trackingBranchName,
				RemoteSHA:  None[gitdomain.SHA](), // will be added later
			})
		default:
			remoteBranchName := gitdomain.NewRemoteBranchName(strings.TrimPrefix(branchName, "remotes/"))
			if existingBranchWithTracking, hasExistingBranchWithTracking := result.FindByRemoteName(remoteBranchName).Get(); hasExistingBranchWithTracking {
				existingBranchWithTracking.RemoteSHA = sha
			} else {
				result = append(result, gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(remoteBranchName),
					RemoteSHA:  sha,
				})
			}
		}
	}
	return result, checkedoutBranch
}

func determineSyncStatus(branchName, remoteText string) (syncStatus gitdomain.SyncStatus, trackingBranchName Option[gitdomain.RemoteBranchName]) {
	if isInSync, trackingBranchName := IsInSync(branchName, remoteText); isInSync {
		return gitdomain.SyncStatusUpToDate, trackingBranchName
	}
	if isGone, trackingBranchName := IsRemoteGone(branchName, remoteText); isGone {
		return gitdomain.SyncStatusDeletedAtRemote, trackingBranchName
	}
	if isAhead, trackingBranchName := IsAhead(branchName, remoteText); isAhead {
		return gitdomain.SyncStatusAhead, trackingBranchName
	}
	if isBehind, trackingBranchName := IsBehind(branchName, remoteText); isBehind {
		return gitdomain.SyncStatusBehind, trackingBranchName
	}
	if isAheadAndBehind, trackingBranchName := IsAheadAndBehind(branchName, remoteText); isAheadAndBehind {
		return gitdomain.SyncStatusNotInSync, trackingBranchName
	}
	if strings.HasPrefix(branchName, "remotes/") {
		return gitdomain.SyncStatusRemoteOnly, None[gitdomain.RemoteBranchName]()
	}
	return gitdomain.SyncStatusLocalOnly, None[gitdomain.RemoteBranchName]()
}

// indicates whether the branch with the given name exists locally
func isLocalBranchName(branch string) bool {
	return !strings.HasPrefix(branch, "remotes/")
}
