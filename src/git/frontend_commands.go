package git

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/messages"
)

type FrontendRunner interface {
	Run(executable string, args ...string) error
	RunMany(commands [][]string) error
}

// FrontendCommands are Git commands that Git Town executes for the user to change the user's repository.
// They can take a while to execute (fetch, push) and stream their output to the user.
// Git Town only needs to know the exit code of frontend commands.
type FrontendCommands struct {
	FrontendRunner
	SetCachedCurrentBranch SetCachedCurrentBranchFunc
}

type SetCachedCurrentBranchFunc func(domain.LocalBranchName)

// AbortMerge cancels a currently ongoing Git merge operation.
func (self *FrontendCommands) AbortMerge() error {
	return self.Run("git", "merge", "--abort")
}

// AbortRebase cancels a currently ongoing Git rebase operation.
func (self *FrontendCommands) AbortRebase() error {
	return self.Run("git", "rebase", "--abort")
}

// AddGitAlias sets the given Git alias.
func (self *FrontendCommands) AddGitAlias(alias configdomain.Alias) error {
	aliasKey := configdomain.NewAliasKey(alias)
	return self.Run("git", "config", "--global", aliasKey.String(), "town "+alias.String())
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (self *FrontendCommands) CheckoutBranch(name domain.LocalBranchName) error {
	err := self.Run("git", "checkout", name.String())
	if err != nil {
		return fmt.Errorf(messages.BranchCheckoutProblem, name, err)
	}
	self.SetCachedCurrentBranch(name)
	return nil
}

// Commit performs a commit of the staged changes with an optional custom message and author.
func (self *FrontendCommands) Commit(message, author string) error {
	gitArgs := []string{"commit"}
	if message != "" {
		gitArgs = append(gitArgs, "-m", message)
	}
	if author != "" {
		gitArgs = append(gitArgs, "--author", author)
	}
	return self.Run("git", gitArgs...)
}

// CommitNoEdit commits all staged files with the default commit message.
func (self *FrontendCommands) CommitNoEdit() error {
	return self.Run("git", "commit", "--no-edit")
}

// CommitStagedChanges commits the currently staged changes.
func (self *FrontendCommands) CommitStagedChanges(message string) error {
	if message != "" {
		return self.Run("git", "commit", "-m", message)
	}
	return self.Run("git", "commit", "--no-edit")
}

// ContinueRebase continues the currently ongoing rebase.
func (self *FrontendCommands) ContinueRebase() error {
	return self.Run("git", "rebase", "--continue")
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (self *FrontendCommands) CreateBranch(name domain.LocalBranchName, parent domain.Location) error {
	return self.Run("git", "branch", name.String(), parent.String())
}

// CreateRemoteBranch creates a remote branch from the given local SHA.
func (self *FrontendCommands) CreateRemoteBranch(localSHA domain.SHA, branch domain.LocalBranchName, noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, domain.OriginRemote.String(), localSHA.String()+":refs/heads/"+branch.String())
	return self.Run("git", args...)
}

// PushBranch pushes the branch with the given name to origin.
func (self *FrontendCommands) CreateTrackingBranch(branch domain.LocalBranchName, remote domain.Remote, noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, "-u", remote.String())
	args = append(args, branch.String())
	return self.Run("git", args...)
}

// DeleteLastCommit resets HEAD to the previous commit.
func (self *FrontendCommands) DeleteLastCommit() error {
	return self.Run("git", "reset", "--hard", "HEAD~1")
}

// DeleteLocalBranch removes the local branch with the given name.
func (self *FrontendCommands) DeleteLocalBranch(name domain.LocalBranchName, force bool) error {
	args := []string{"branch", "-d", name.String()}
	if force {
		args[1] = "-D"
	}
	return self.Run("git", args...)
}

// DeleteRemoteBranch removes the remote branch of the given local branch.
func (self *FrontendCommands) DeleteRemoteBranch(name domain.RemoteBranchName) error {
	remote, localBranchName := name.Parts()
	return self.Run("git", "push", remote.String(), ":"+localBranchName.String())
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (self *FrontendCommands) DiffParent(branch, parentBranch domain.LocalBranchName) error {
	return self.Run("git", "diff", parentBranch.String()+".."+branch.String())
}

// DiscardOpenChanges deletes all uncommitted changes.
func (self *FrontendCommands) DiscardOpenChanges() error {
	return self.Run("git", "reset", "--hard")
}

// Fetch retrieves the updates from the origin repo.
func (self *FrontendCommands) Fetch() error {
	return self.Run("git", "fetch", "--prune", "--tags")
}

// FetchUpstream fetches updates from the upstream remote.
func (self *FrontendCommands) FetchUpstream(branch domain.LocalBranchName) error {
	return self.Run("git", "fetch", domain.UpstreamRemote.String(), branch.String())
}

// PushBranch pushes the branch with the given name to origin.
func (self *FrontendCommands) ForcePushBranch(noPushHook configdomain.NoPushHook) error {
	args := []string{"push", "--force-with-lease"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return self.Run("git", args...)
}

// MergeBranchNoEdit merges the given branch into the current branch,
// using the default commit message.
func (self *FrontendCommands) MergeBranchNoEdit(branch domain.BranchName) error {
	return self.Run("git", "merge", "--no-edit", branch.String())
}

// NavigateToDir changes into the root directory of the current repository.
func (self *FrontendCommands) NavigateToDir(dir domain.RepoRootDir) error {
	return os.Chdir(dir.String())
}

// PopStash restores stashed-away changes into the workspace.
func (self *FrontendCommands) PopStash() error {
	return self.Run("git", "stash", "pop")
}

// Pull fetches updates from origin and updates the currently checked out branch.
func (self *FrontendCommands) Pull() error {
	return self.Run("git", "pull")
}

// PushCurrentBranch pushes the current branch to its tracking branch.
func (self *FrontendCommands) PushCurrentBranch(noPushHook configdomain.NoPushHook) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return self.Run("git", args...)
}

// PushTags pushes new the Git tags to origin.
func (self *FrontendCommands) PushTags() error {
	return self.Run("git", "push", "--tags")
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (self *FrontendCommands) Rebase(target domain.BranchName) error {
	return self.Run("git", "rebase", target.String())
}

// RemoveGitAlias removes the given Git alias.
func (self *FrontendCommands) RemoveGitAlias(alias configdomain.Alias) error {
	return self.Run("git", "config", "--global", "--unset", "alias."+alias.String())
}

// ResetCurrentBranchToSHA undoes all commits on the current branch all the way until the given SHA.
func (self *FrontendCommands) ResetCurrentBranchToSHA(sha domain.SHA, hard bool) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	args = append(args, sha.String())
	return self.Run("git", args...)
}

// ResetRemoteBranchToSHA sets the given remote branch to the given SHA.
func (self *FrontendCommands) ResetRemoteBranchToSHA(branch domain.RemoteBranchName, sha domain.SHA) error {
	return self.Run("git", "push", "--force-with-lease", domain.OriginRemote.String(), sha.String()+":"+branch.LocalBranchName().String())
}

// RevertCommit reverts the commit with the given SHA.
func (self *FrontendCommands) RevertCommit(sha domain.SHA) error {
	return self.Run("git", "revert", sha.String())
}

// SquashMerge squash-merges the given branch into the current branch.
func (self *FrontendCommands) SquashMerge(branch domain.LocalBranchName) error {
	return self.Run("git", "merge", "--squash", branch.String())
}

// StageFiles adds the file with the given name to the Git index.
func (self *FrontendCommands) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	return self.Run("git", args...)
}

// StartCommit starts a commit and stops at asking the user for the commit message.
func (self *FrontendCommands) StartCommit() error {
	return self.Run("git", "commit")
}

// Stash adds the current files to the Git stash.
func (self *FrontendCommands) Stash() error {
	return self.RunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
}

func (self *FrontendCommands) UndoLastCommit() error {
	return self.Run("git", "reset", "--soft", "HEAD^")
}
