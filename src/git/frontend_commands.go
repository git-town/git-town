package git

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

type FrontendRunner interface {
	Run(executable string, args ...string) error
	RunMany(commands [][]string) error
}

// FrontendCommands are Git commands that Git Town executes for the user to change the user's repository.
// They can take a while to execute (fetch, push) and stream their output to the user.
// Git Town only needs to know the exit code of frontend commands.
type FrontendCommands struct {
	Runner                 FrontendRunner
	SetCachedCurrentBranch SetCachedCurrentBranchFunc
}

type SetCachedCurrentBranchFunc func(gitdomain.LocalBranchName)

// AbortMerge cancels a currently ongoing Git merge operation.
func (self *FrontendCommands) AbortMerge() error {
	return self.Runner.Run("git", "merge", "--abort")
}

// AbortRebase cancels a currently ongoing Git rebase operation.
func (self *FrontendCommands) AbortRebase() error {
	return self.Runner.Run("git", "rebase", "--abort")
}

// CheckoutBranch checks out the Git branch with the given name in this repo,
// optionally using a three-way merge.
func (self *FrontendCommands) CheckoutBranch(name gitdomain.LocalBranchName, merge bool) error {
	args := []string{"checkout", name.String()}
	if merge {
		args = append(args, "-m")
	}
	err := self.Runner.Run("git", args...)
	self.SetCachedCurrentBranch(name)
	if err != nil {
		return fmt.Errorf(messages.BranchCheckoutProblem, name, err)
	}
	return nil
}

// Commit performs a commit of the staged changes with an optional custom message and author.
func (self *FrontendCommands) Commit(message gitdomain.CommitMessage, author gitdomain.Author) error {
	gitArgs := []string{"commit"}
	if message != "" {
		gitArgs = append(gitArgs, "-m", message.String())
	}
	if author != "" {
		gitArgs = append(gitArgs, "--author", author.String())
	}
	return self.Runner.Run("git", gitArgs...)
}

// CommitNoEdit commits all staged files with the default commit message.
func (self *FrontendCommands) CommitNoEdit() error {
	return self.Runner.Run("git", "commit", "--no-edit")
}

// CommitStagedChanges commits the currently staged changes.
func (self *FrontendCommands) CommitStagedChanges(message string) error {
	if message != "" {
		return self.Runner.Run("git", "commit", "-m", message)
	}
	return self.Runner.Run("git", "commit", "--no-edit")
}

// ContinueRebase continues the currently ongoing rebase.
func (self *FrontendCommands) ContinueRebase() error {
	return self.Runner.Run("git", "rebase", "--continue")
}

// CreateAndCheckoutBranch creates a new branch with the given name and checks it out using a single Git operation.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *FrontendCommands) CreateAndCheckoutBranch(name gitdomain.LocalBranchName) error {
	err := self.Runner.Run("git", "checkout", "-b", name.String())
	self.SetCachedCurrentBranch(name)
	return err
}

// CreateAndCheckoutBranchWithParent creates a new branch with the given name and parent and checks it out using a single Git operation.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *FrontendCommands) CreateAndCheckoutBranchWithParent(name gitdomain.LocalBranchName, parent gitdomain.Location) error {
	err := self.Runner.Run("git", "checkout", "-b", name.String(), parent.String())
	self.SetCachedCurrentBranch(name)
	return err
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *FrontendCommands) CreateBranch(name gitdomain.LocalBranchName, parent gitdomain.Location) error {
	return self.Runner.Run("git", "branch", name.String(), parent.String())
}

// CreateRemoteBranch creates a remote branch from the given local SHA.
func (self *FrontendCommands) CreateRemoteBranch(localSHA gitdomain.SHA, branch gitdomain.LocalBranchName, noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, gitdomain.RemoteOrigin.String(), localSHA.String()+":refs/heads/"+branch.String())
	return self.Runner.Run("git", args...)
}

// PushBranch pushes the branch with the given name to origin.
func (self *FrontendCommands) CreateTrackingBranch(branch gitdomain.LocalBranchName, remote gitdomain.Remote, noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, "-u", remote.String())
	args = append(args, branch.String())
	return self.Runner.Run("git", args...)
}

// DeleteHostingPlatform removes the hosting platform config entry.
func (self *FrontendCommands) DeleteHostingPlatform() error {
	return self.Runner.Run("git", "config", "--unset", gitconfig.KeyHostingPlatform.String())
}

// DeleteLastCommit resets HEAD to the previous commit.
func (self *FrontendCommands) DeleteLastCommit() error {
	return self.Runner.Run("git", "reset", "--hard", "HEAD~1")
}

// DeleteLocalBranch removes the local branch with the given name.
func (self *FrontendCommands) DeleteLocalBranch(name gitdomain.LocalBranchName) error {
	return self.Runner.Run("git", "branch", "-D", name.String())
}

// DeleteOriginHostname removes the origin hostname override
func (self *FrontendCommands) DeleteOriginHostname() error {
	return self.Runner.Run("git", "config", "--unset", gitconfig.KeyHostingOriginHostname.String())
}

// DeleteTrackingBranch removes the tracking branch of the given local branch.
func (self *FrontendCommands) DeleteTrackingBranch(name gitdomain.RemoteBranchName) error {
	remote, localBranchName := name.Parts()
	return self.Runner.Run("git", "push", remote.String(), ":"+localBranchName.String())
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (self *FrontendCommands) DiffParent(branch, parentBranch gitdomain.LocalBranchName) error {
	return self.Runner.Run("git", "diff", parentBranch.String()+".."+branch.String())
}

// DiscardOpenChanges deletes all uncommitted changes.
func (self *FrontendCommands) DiscardOpenChanges() error {
	return self.Runner.Run("git", "reset", "--hard")
}

// Fetch retrieves the updates from the origin repo.
func (self *FrontendCommands) Fetch() error {
	return self.Runner.Run("git", "fetch", "--prune", "--tags")
}

// FetchUpstream fetches updates from the upstream remote.
func (self *FrontendCommands) FetchUpstream(branch gitdomain.LocalBranchName) error {
	return self.Runner.Run("git", "fetch", gitdomain.RemoteUpstream.String(), branch.String())
}

// PushBranch pushes the branch with the given name to origin.
func (self *FrontendCommands) ForcePushBranchSafely(noPushHook configdomain.NoPushHook) error {
	args := []string{"push", "--force-with-lease", "--force-if-includes"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return self.Runner.Run("git", args...)
}

// MergeBranchNoEdit merges the given branch into the current branch,
// using the default commit message.
func (self *FrontendCommands) MergeBranchNoEdit(branch gitdomain.BranchName) error {
	return self.Runner.Run("git", "merge", "--no-edit", "--ff", branch.String())
}

// NavigateToDir changes into the root directory of the current repository.
func (self *FrontendCommands) NavigateToDir(dir gitdomain.RepoRootDir) error {
	return os.Chdir(dir.String())
}

// PopStash restores stashed-away changes into the workspace.
func (self *FrontendCommands) PopStash() error {
	return self.Runner.Run("git", "stash", "pop")
}

// Pull fetches updates from origin and updates the currently checked out branch.
func (self *FrontendCommands) Pull() error {
	return self.Runner.Run("git", "pull")
}

// PushCurrentBranch pushes the current branch to its tracking branch.
func (self *FrontendCommands) PushCurrentBranch(noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return self.Runner.Run("git", args...)
}

// PushTags pushes new the Git tags to origin.
func (self *FrontendCommands) PushTags() error {
	return self.Runner.Run("git", "push", "--tags")
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (self *FrontendCommands) Rebase(target gitdomain.BranchName) error {
	return self.Runner.Run("git", "rebase", target.String())
}

// ResetCurrentBranchToSHA undoes all commits on the current branch all the way until the given SHA.
func (self *FrontendCommands) RemoveCommitsInCurrentBranch(parent gitdomain.LocalBranchName) error {
	return self.Runner.Run("git", "reset", "--soft", parent.String())
}

// RemoveGitAlias removes the given Git alias.
func (self *FrontendCommands) RemoveGitAlias(aliasableCommand configdomain.AliasableCommand) error {
	aliasKey := gitconfig.KeyForAliasableCommand(aliasableCommand)
	return self.Runner.Run("git", "config", "--global", "--unset", aliasKey.String())
}

// RemoveHubToken removes the stored token for the GitHub API.
func (self *FrontendCommands) RemoveGitHubToken() error {
	return self.Runner.Run("git", "config", "--unset", gitconfig.KeyGithubToken.String())
}

// RemoveHubToken removes the stored token for the GitHub API.
func (self *FrontendCommands) RemoveGitLabToken() error {
	return self.Runner.Run("git", "config", "--unset", gitconfig.KeyGitlabToken.String())
}

// RemoveHubToken removes the stored token for the GitHub API.
func (self *FrontendCommands) RemoveGiteaToken() error {
	return self.Runner.Run("git", "config", "--unset", gitconfig.KeyGiteaToken.String())
}

// ResetCurrentBranchToSHA undoes all commits on the current branch all the way until the given SHA.
func (self *FrontendCommands) ResetCurrentBranchToSHA(sha gitdomain.SHA, hard bool) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	args = append(args, sha.String())
	return self.Runner.Run("git", args...)
}

// ResetRemoteBranchToSHA sets the given remote branch to the given SHA.
func (self *FrontendCommands) ResetRemoteBranchToSHA(branch gitdomain.RemoteBranchName, sha gitdomain.SHA) error {
	return self.Runner.Run("git", "push", "--force-with-lease", gitdomain.RemoteOrigin.String(), sha.String()+":"+branch.LocalBranchName().String())
}

// RevertCommit reverts the commit with the given SHA.
func (self *FrontendCommands) RevertCommit(sha gitdomain.SHA) error {
	return self.Runner.Run("git", "revert", sha.String())
}

// SetGitAlias sets the given Git alias.
func (self *FrontendCommands) SetGitAlias(aliasableCommand configdomain.AliasableCommand) error {
	return self.Runner.Run("git", "config", "--global", gitconfig.KeyForAliasableCommand(aliasableCommand).String(), "town "+aliasableCommand.String())
}

// SetGitHubToken sets the given API token for the GitHub API.
func (self *FrontendCommands) SetGitHubToken(value configdomain.GitHubToken) error {
	return self.Runner.Run("git", "config", gitconfig.KeyGithubToken.String(), value.String())
}

// SetGitLabToken sets the given API token for the GitHub API.
func (self *FrontendCommands) SetGitLabToken(value configdomain.GitLabToken) error {
	return self.Runner.Run("git", "config", gitconfig.KeyGitlabToken.String(), value.String())
}

// SetGiteaToken sets the given API token for the Gitea API.
func (self *FrontendCommands) SetGiteaToken(value configdomain.GiteaToken) error {
	return self.Runner.Run("git", "config", gitconfig.KeyGiteaToken.String(), value.String())
}

// SetHostingPlatform sets the given code hosting platform.
func (self *FrontendCommands) SetHostingPlatform(platform configdomain.HostingPlatform) error {
	return self.Runner.Run("git", "config", gitconfig.KeyHostingPlatform.String(), platform.String())
}

// SetHostingPlatform sets the given code hosting platform.
func (self *FrontendCommands) SetOriginHostname(hostname configdomain.HostingOriginHostname) error {
	return self.Runner.Run("git", "config", gitconfig.KeyHostingOriginHostname.String(), hostname.String())
}

// SquashMerge squash-merges the given branch into the current branch.
func (self *FrontendCommands) SquashMerge(branch gitdomain.LocalBranchName) error {
	return self.Runner.Run("git", "merge", "--squash", "--ff", branch.String())
}

// StageFiles adds the file with the given name to the Git index.
func (self *FrontendCommands) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	return self.Runner.Run("git", args...)
}

// StartCommit starts a commit and stops at asking the user for the commit message.
func (self *FrontendCommands) StartCommit() error {
	return self.Runner.Run("git", "commit")
}

// Stash adds the current files to the Git stash.
func (self *FrontendCommands) Stash() error {
	return self.Runner.RunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
}

func (self *FrontendCommands) UndoLastCommit() error {
	return self.Runner.Run("git", "reset", "--soft", "HEAD~1")
}
