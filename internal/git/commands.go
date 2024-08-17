package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks/cache"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v15/internal/messages"
)

// Commands are Git commands that Git Town executes to determine which frontend commands to run.
// They don't change the user's repo, execute instantaneously, and Git Town needs to know their output.
// They are invisible to the end user unless the "verbose" option is set.
type Commands struct {
	CurrentBranchCache *cache.LocalBranchWithPrevious // caches the currently checked out Git branch
	RemotesCache       *cache.Remotes                 // caches Git remotes
}

// AbortMerge cancels a currently ongoing Git merge operation.
func (self *Commands) AbortMerge(runner gitdomain.Runner) error {
	return runner.Run("git", "merge", "--abort")
}

// AbortRebase cancels a currently ongoing Git rebase operation.
func (self *Commands) AbortRebase(runner gitdomain.Runner) error {
	return runner.Run("git", "rebase", "--abort")
}

// BranchAuthors provides the user accounts that contributed to the given branch.
// Returns lines of "name <email>".
func (self *Commands) BranchAuthors(querier gitdomain.Querier, branch, parent gitdomain.LocalBranchName) ([]gitdomain.Author, error) {
	output, err := querier.QueryTrim("git", "shortlog", "-s", "-n", "-e", parent.String()+".."+branch.String())
	if err != nil {
		return []gitdomain.Author{}, err
	}
	result := []gitdomain.Author{}
	for _, line := range stringslice.Lines(output) {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		result = append(result, gitdomain.Author(parts[1]))
	}
	return result, nil
}

func (self *Commands) BranchExists(runner gitdomain.Runner, branch gitdomain.LocalBranchName) bool {
	err := runner.Run("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch.String())
	return err == nil
}

// BranchHasUnmergedChanges indicates whether the branch with the given name
// contains changes that were not merged into the main branch.
func (self *Commands) BranchHasUnmergedChanges(querier gitdomain.Querier, branch, parent gitdomain.LocalBranchName) (bool, error) {
	out, err := querier.QueryTrim("git", "diff", parent.String()+".."+branch.String())
	if err != nil {
		return false, fmt.Errorf(messages.BranchDiffProblem, branch, err)
	}
	return out != "", nil
}

// BranchesSnapshot provides detailed information about the sync status of all branches.
func (self *Commands) BranchesSnapshot(querier gitdomain.Querier) (gitdomain.BranchesSnapshot, error) { //nolint:nonamedreturns
	output, err := querier.Query("git", "branch", "-vva", "--sort=refname")
	if err != nil {
		return gitdomain.EmptyBranchesSnapshot(), err
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

// CheckoutBranch checks out the Git branch with the given name.
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

// CheckoutBranch checks out the Git branch with the given name.
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

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// Adds the given prefix with the newline if provided.
func (self *Commands) CommentOutSquashCommitMessage(prefix string) error {
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

// Commit performs a commit of the staged changes with an optional custom message and author.
func (self *Commands) Commit(runner gitdomain.Runner, message Option[gitdomain.CommitMessage], author gitdomain.Author) error {
	gitArgs := []string{"commit"}
	if messageContent, has := message.Get(); has {
		gitArgs = append(gitArgs, "-m", messageContent.String())
	}
	if author != "" {
		gitArgs = append(gitArgs, "--author", author.String())
	}
	return runner.Run("git", gitArgs...)
}

// CommitNoEdit commits all staged files with the default commit message.
func (self *Commands) CommitNoEdit(runner gitdomain.Runner) error {
	return runner.Run("git", "commit", "--no-edit")
}

// CommitStagedChanges commits the currently staged changes.
func (self *Commands) CommitStagedChanges(runner gitdomain.Runner, message gitdomain.CommitMessage) error {
	if message != "" {
		return runner.Run("git", "commit", "-m", message.String())
	}
	return runner.Run("git", "commit", "--no-edit")
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
	output, err := querier.QueryTrim("git", "log", "--pretty=format:%h %s", "-10")
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

// ContinueRebase continues the currently ongoing rebase.
func (self *Commands) ContinueRebase(runner gitdomain.Runner) error {
	return runner.Run("git", "rebase", "--continue")
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
	err := runner.Run("git", "checkout", "-b", name.String(), parent.String())
	self.CurrentBranchCache.Set(name)
	return err
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *Commands) CreateBranch(runner gitdomain.Runner, name gitdomain.LocalBranchName, parent gitdomain.Location) error {
	return runner.Run("git", "branch", name.String(), parent.String())
}

// CreateRemoteBranch creates a remote branch from the given local SHA.
func (self *Commands) CreateRemoteBranch(runner gitdomain.Runner, localSHA gitdomain.SHA, branch gitdomain.LocalBranchName, noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, gitdomain.RemoteOrigin.String(), localSHA.String()+":refs/heads/"+branch.String())
	return runner.Run("git", args...)
}

// pushes the branch with the given name to origin
func (self *Commands) CreateTrackingBranch(runner gitdomain.Runner, branch gitdomain.LocalBranchName, remote gitdomain.Remote, noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, "-u", remote.String())
	args = append(args, branch.String())
	return runner.Run("git", args...)
}

// CurrentBranch provides the name of the currently checked out branch.
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

// CurrentBranch provides the currently checked out branch.
func (self *Commands) CurrentBranchUncached(querier gitdomain.Querier) (gitdomain.LocalBranchName, error) {
	// first try to detect the current branch the normal way
	output, err := querier.QueryTrim("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err == nil && output != "HEAD" {
		return gitdomain.NewLocalBranchName(output), nil
	}
	// here we couldn't detect the current branch the normal way --> assume we are in a rebase and try the rebase way
	return self.currentBranchDuringRebase(querier)
}

// CurrentSHA provides the SHA of the currently checked out branch/commit.
func (self *Commands) CurrentSHA(querier gitdomain.Querier) (gitdomain.SHA, error) {
	return self.SHAForBranch(querier, gitdomain.NewBranchName("HEAD"))
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

// DeleteHostingPlatform removes the hosting platform config entry.
func (self *Commands) DeleteHostingPlatform(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyHostingPlatform.String())
}

// DeleteLastCommit resets HEAD to the previous commit.
func (self *Commands) DeleteLastCommit(runner gitdomain.Runner) error {
	return runner.Run("git", "reset", "--hard", "HEAD~1")
}

// DeleteLocalBranch removes the local branch with the given name.
func (self *Commands) DeleteLocalBranch(runner gitdomain.Runner, name gitdomain.LocalBranchName) error {
	return runner.Run("git", "branch", "-D", name.String())
}

// DeleteOriginHostname removes the origin hostname override
func (self *Commands) DeleteOriginHostname(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyHostingOriginHostname.String())
}

// DeleteTrackingBranch removes the tracking branch of the given local branch.
func (self *Commands) DeleteTrackingBranch(runner gitdomain.Runner, name gitdomain.RemoteBranchName) error {
	remote, localBranchName := name.Parts()
	return runner.Run("git", "push", remote.String(), ":"+localBranchName.String())
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (self *Commands) DiffParent(runner gitdomain.Runner, branch, parentBranch gitdomain.LocalBranchName) error {
	return runner.Run("git", "diff", parentBranch.String()+".."+branch.String())
}

// DiscardOpenChanges deletes all uncommitted changes.
func (self *Commands) DiscardOpenChanges(runner gitdomain.Runner) error {
	return runner.Run("git", "reset", "--hard")
}

// Fetch retrieves the updates from the origin repo.
func (self *Commands) Fetch(runner gitdomain.Runner, syncTags configdomain.SyncTags) error {
	if syncTags.IsTrue() {
		return runner.Run("git", "fetch", "--prune", "--tags")
	}
	return runner.Run("git", "fetch", "--prune", "--no-tags")
}

// FetchUpstream fetches updates from the upstream remote.
func (self *Commands) FetchUpstream(runner gitdomain.Runner, branch gitdomain.LocalBranchName) error {
	return runner.Run("git", "fetch", gitdomain.RemoteUpstream.String(), branch.String())
}

// provides the commit message of the first commit in the branch with the given name
func (self *Commands) FirstCommitMessageInBranch(runner gitdomain.Querier, branch, mainBranch gitdomain.LocalBranchName) (Option[gitdomain.CommitMessage], error) {
	output, err := runner.QueryTrim("git", "log", fmt.Sprintf("%s..%s", mainBranch, branch), "--format=%h")
	if err != nil {
		return None[gitdomain.CommitMessage](), err
	}
	hashes := stringslice.Lines(output)
	if len(hashes) == 0 {
		return None[gitdomain.CommitMessage](), nil
	}
	hash := hashes[len(hashes)-1]
	message, err := runner.QueryTrim("git", "log", "--format=%B", "-n", "1", hash)
	if err != nil {
		return None[gitdomain.CommitMessage](), err
	}
	return Some(gitdomain.CommitMessage(message)), nil
}

func (self *Commands) FirstExistingBranch(runner gitdomain.Runner, branches ...gitdomain.LocalBranchName) Option[gitdomain.LocalBranchName] {
	for _, branch := range branches {
		if self.BranchExists(runner, branch) {
			return Some(branch)
		}
	}
	return None[gitdomain.LocalBranchName]()
}

// pushes the branch with the given name to origin.
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

// HasLocalBranch indicates whether this repo has a local branch with the given name.
func (self *Commands) HasLocalBranch(runner gitdomain.Runner, name gitdomain.LocalBranchName) bool {
	return runner.Run("git", "show-ref", "--quiet", "refs/heads/"+name.String()) == nil
}

// HasMergeInProgress indicates whether this Git repository currently has a merge in progress.
func (self *Commands) HasMergeInProgress(runner gitdomain.Runner) bool {
	err := runner.Run("git", "rev-parse", "-q", "--verify", "MERGE_HEAD")
	return err == nil
}

// HasShippableChanges indicates whether the given branch has changes
// not currently in the main branch.
func (self *Commands) HasShippableChanges(querier gitdomain.Querier, branch, mainBranch gitdomain.LocalBranchName) (bool, error) {
	out, err := querier.QueryTrim("git", "diff", mainBranch.String()+".."+branch.String())
	if err != nil {
		return false, fmt.Errorf(messages.ShippableChangesProblem, branch, err)
	}
	return out != "", nil
}

// LastCommitMessage provides the commit message for the last commit.
func (self *Commands) LastCommitMessage(querier gitdomain.Querier) (gitdomain.CommitMessage, error) {
	out, err := querier.QueryTrim("git", "log", "-1", "--format=%B")
	if err != nil {
		return "", fmt.Errorf(messages.CommitMessageProblem, err)
	}
	return gitdomain.CommitMessage(out), nil
}

// MergeBranchNoEdit merges the given branch into the current branch,
// using the default commit message.
func (self *Commands) MergeBranchNoEdit(runner gitdomain.Runner, branch gitdomain.BranchName) error {
	return runner.Run("git", "merge", "--no-edit", "--ff", branch.String())
}

// NavigateToDir changes into the root directory of the current repository.
func (self *Commands) NavigateToDir(dir gitdomain.RepoRootDir) error {
	return os.Chdir(dir.String())
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

// PopStash restores stashed-away changes into the workspace.
func (self *Commands) PopStash(runner gitdomain.Runner) error {
	return runner.Run("git", "stash", "pop")
}

// PreviouslyCheckedOutBranch provides the name of the branch that was previously checked out in this repo.
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

// Pull fetches updates from origin and updates the currently checked out branch.
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

// PushTags pushes new the Git tags to origin.
func (self *Commands) PushTags(runner gitdomain.Runner) error {
	return runner.Run("git", "push", "--tags")
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (self *Commands) Rebase(runner gitdomain.Runner, target gitdomain.BranchName) error {
	return runner.Run("git", "rebase", target.String())
}

// Remotes provides the names of all Git remotes in this repository.
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

// Remotes provides the names of all Git remotes in this repository.
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

// RemoveGitAlias removes the given Git alias.
func (self *Commands) RemoveGitAlias(runner gitdomain.Runner, aliasableCommand configdomain.AliasableCommand) error {
	return runner.Run("git", "config", "--global", "--unset", aliasableCommand.Key().String())
}

// RemoveHubToken removes the stored token for the GitHub API.
func (self *Commands) RemoveGitHubToken(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyGithubToken.String())
}

// RemoveHubToken removes the stored token for the GitHub API.
func (self *Commands) RemoveGitLabToken(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyGitlabToken.String())
}

// RemoveHubToken removes the stored token for the GitHub API.
func (self *Commands) RemoveGiteaToken(runner gitdomain.Runner) error {
	return runner.Run("git", "config", "--unset", configdomain.KeyGiteaToken.String())
}

// RepoStatus provides a summary of the state the current workspace is in right now: rebasing, has conflicts, has open changes, etc.
func (self *Commands) RepoStatus(querier gitdomain.Querier) (gitdomain.RepoStatus, error) {
	output, err := querier.QueryTrim("git", "status", "--long", "--ignore-submodules")
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

// ResetCurrentBranchToSHA undoes all commits on the current branch all the way until the given SHA.
func (self *Commands) ResetBranch(runner gitdomain.Runner, target gitdomain.BranchName) error {
	return runner.Run("git", "reset", "--soft", target.String())
}

// ResetCurrentBranchToSHA undoes all commits on the current branch all the way until the given SHA.
func (self *Commands) ResetCurrentBranchToSHA(runner gitdomain.Runner, sha gitdomain.SHA, hard bool) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	args = append(args, sha.String())
	return runner.Run("git", args...)
}

// ResetRemoteBranchToSHA sets the given remote branch to the given SHA.
func (self *Commands) ResetRemoteBranchToSHA(runner gitdomain.Runner, branch gitdomain.RemoteBranchName, sha gitdomain.SHA) error {
	return runner.Run("git", "push", "--force-with-lease", gitdomain.RemoteOrigin.String(), sha.String()+":"+branch.LocalBranchName().String())
}

// RevertCommit reverts the commit with the given SHA.
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

// SHAForBranch provides the SHA for the local branch with the given name.
func (self *Commands) SHAForBranch(querier gitdomain.Querier, name gitdomain.BranchName) (gitdomain.SHA, error) {
	output, err := querier.QueryTrim("git", "rev-parse", "--short", name.String())
	if err != nil {
		return gitdomain.SHA(""), fmt.Errorf(messages.BranchLocalSHAProblem, name, err)
	}
	return gitdomain.NewSHA(output), nil
}

// SetGitAlias sets the given Git alias.
func (self *Commands) SetGitAlias(runner gitdomain.Runner, aliasableCommand configdomain.AliasableCommand) error {
	return runner.Run("git", "config", "--global", aliasableCommand.Key().String(), "town "+aliasableCommand.String())
}

// SetGitHubToken sets the given API token for the GitHub API.
func (self *Commands) SetGitHubToken(runner gitdomain.Runner, value configdomain.GitHubToken) error {
	return runner.Run("git", "config", configdomain.KeyGithubToken.String(), value.String())
}

// SetGitLabToken sets the given API token for the GitHub API.
func (self *Commands) SetGitLabToken(runner gitdomain.Runner, value configdomain.GitLabToken) error {
	return runner.Run("git", "config", configdomain.KeyGitlabToken.String(), value.String())
}

// SetGiteaToken sets the given API token for the Gitea API.
func (self *Commands) SetGiteaToken(runner gitdomain.Runner, value configdomain.GiteaToken) error {
	return runner.Run("git", "config", configdomain.KeyGiteaToken.String(), value.String())
}

// SetHostingPlatform sets the given code hosting platform.
func (self *Commands) SetHostingPlatform(runner gitdomain.Runner, platform configdomain.HostingPlatform) error {
	return runner.Run("git", "config", configdomain.KeyHostingPlatform.String(), platform.String())
}

// SetHostingPlatform sets the given code hosting platform.
func (self *Commands) SetOriginHostname(runner gitdomain.Runner, hostname configdomain.HostingOriginHostname) error {
	return runner.Run("git", "config", configdomain.KeyHostingOriginHostname.String(), hostname.String())
}

// ShouldPushBranch returns whether the local branch with the given name
// contains commits that have not been pushed to its tracking branch.
func (self *Commands) ShouldPushBranch(querier gitdomain.Querier, branch gitdomain.LocalBranchName) (bool, error) {
	out, err := querier.QueryTrim("git", "rev-list", "--left-right", branch.String()+"..."+branch.TrackingBranch().String())
	if err != nil {
		return false, fmt.Errorf(messages.DiffProblem, branch, branch, err)
	}
	return out != "", nil
}

// SquashMerge squash-merges the given branch into the current branch.
func (self *Commands) SquashMerge(runner gitdomain.Runner, branch gitdomain.LocalBranchName) error {
	return runner.Run("git", "merge", "--squash", "--ff", branch.String())
}

// StageFiles adds the file with the given name to the Git index.
func (self *Commands) StageFiles(runner gitdomain.Runner, names ...string) error {
	args := append([]string{"add"}, names...)
	return runner.Run("git", args...)
}

// StartCommit starts a commit and stops at asking the user for the commit message.
func (self *Commands) StartCommit(runner gitdomain.Runner) error {
	return runner.Run("git", "commit")
}

// Stash adds the current files to the Git stash.
func (self *Commands) Stash(runner gitdomain.Runner) error {
	err := runner.Run("git", "add", "-A")
	if err != nil {
		return err
	}
	return runner.Run("git", "stash")
}

// StashSize provides the number of stashes in this repository.
func (self *Commands) StashSize(querier gitdomain.Querier) (gitdomain.StashSize, error) {
	output, err := querier.QueryTrim("git", "stash", "list")
	return gitdomain.StashSize(len(stringslice.Lines(output))), err
}

func (self *Commands) UndoLastCommit(runner gitdomain.Runner) error {
	return runner.Run("git", "reset", "--soft", "HEAD~1")
}

// Version indicates whether the needed Git version is installed.
func (self *Commands) Version(querier gitdomain.Querier) (major int, minor int, err error) {
	versionRegexp := regexp.MustCompile(`git version (\d+).(\d+).(\d+)`)
	output, err := querier.QueryTrim("git", "version")
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

// IsRemoteGone indicates whether the given remoteText indicates a deleted tracking branch.
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

func ParseActiveBranchDuringRebase(lineWithStar string) gitdomain.LocalBranchName {
	parts := strings.Split(lineWithStar, " ")
	partsWithBranchName := parts[4:]
	branchNameWithClosingParen := strings.Join(partsWithBranchName, " ")
	return gitdomain.NewLocalBranchName(branchNameWithClosingParen[:len(branchNameWithClosingParen)-1])
}

// ParseVerboseBranchesOutput provides the branches in the given Git output as well as the name of the currently checked out branch.
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
		if parts[0] == "(no" || parts[1] == "(no" || parts[1] == "branch," {
			continue
		}
		branchName := parts[0]
		sha := gitdomain.NewSHA(parts[1])
		remoteText := parts[2]
		if line[0] == '*' { // "(no" as in "(no branch, rebasing main)" is what we get when a rebase is active, in which case no branch is checked out
			checkedoutBranch = Some(gitdomain.LocalBranchName(branchName))
		}
		syncStatus, trackingBranchName := determineSyncStatus(branchName, remoteText)
		switch {
		case line[0] == '+':
			result = append(result, gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName(branchName)),
				LocalSHA:   Some(sha),
				SyncStatus: gitdomain.SyncStatusOtherWorktree,
				RemoteName: trackingBranchName,
				RemoteSHA:  None[gitdomain.SHA](),
			})
		case isLocalBranchName(branchName):
			result = append(result, gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName(branchName)),
				LocalSHA:   Some(sha),
				SyncStatus: syncStatus,
				RemoteName: trackingBranchName,
				RemoteSHA:  None[gitdomain.SHA](), // will be added later
			})
		default:
			remoteBranchName := gitdomain.NewRemoteBranchName(strings.TrimPrefix(branchName, "remotes/"))
			if existingBranchWithTracking, hasExistingBranchWithTracking := result.FindByRemoteName(remoteBranchName).Get(); hasExistingBranchWithTracking {
				existingBranchWithTracking.RemoteSHA = Some(sha)
			} else {
				result = append(result, gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(remoteBranchName),
					RemoteSHA:  Some(sha),
				})
			}
		}
	}
	return result, checkedoutBranch
}

func determineSyncStatus(branchName, remoteText string) (syncStatus gitdomain.SyncStatus, trackingBranchName Option[gitdomain.RemoteBranchName]) {
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
		return gitdomain.SyncStatusRemoteOnly, None[gitdomain.RemoteBranchName]()
	}
	return gitdomain.SyncStatusLocalOnly, None[gitdomain.RemoteBranchName]()
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
