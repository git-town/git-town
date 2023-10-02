package git

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

type FrontendRunner interface {
	Run(executable string, args ...string) error
	RunMany([][]string) error
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
func (fcs *FrontendCommands) AbortMerge() error {
	return fcs.Run("git", "merge", "--abort")
}

// AbortRebase cancels a currently ongoing Git rebase operation.
func (fcs *FrontendCommands) AbortRebase() error {
	return fcs.Run("git", "rebase", "--abort")
}

// AddGitAlias sets the given Git alias.
func (fcs *FrontendCommands) AddGitAlias(alias config.Alias) error {
	aliasKey := config.NewAliasKey(alias)
	return fcs.Run("git", "config", "--global", aliasKey.String(), "town "+alias.String())
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (fcs *FrontendCommands) CheckoutBranch(name domain.LocalBranchName) error {
	err := fcs.Run("git", "checkout", name.String())
	if err != nil {
		return fmt.Errorf(messages.BranchCheckoutProblem, name, err)
	}
	fcs.SetCachedCurrentBranch(name)
	return nil
}

// CreateRemoteBranch creates a remote branch from the given local SHA.
func (fcs *FrontendCommands) CreateRemoteBranch(localSHA domain.SHA, branch domain.LocalBranchName, noPushHook bool) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, domain.OriginRemote.String(), localSHA.String()+":refs/heads/"+branch.String())
	return fcs.Run("git", args...)
}

// CommitNoEdit commits all staged files with the default commit message.
func (fcs *FrontendCommands) CommitNoEdit() error {
	return fcs.Run("git", "commit", "--no-edit")
}

// CommitStagedChanges commits the currently staged changes.
func (fcs *FrontendCommands) CommitStagedChanges(message string) error {
	if message != "" {
		return fcs.Run("git", "commit", "-m", message)
	}
	return fcs.Run("git", "commit", "--no-edit")
}

// Commit performs a commit of the staged changes with an optional custom message and author.
func (fcs *FrontendCommands) Commit(message, author string) error {
	gitArgs := []string{"commit"}
	if message != "" {
		gitArgs = append(gitArgs, "-m", message)
	}
	if author != "" {
		gitArgs = append(gitArgs, "--author", author)
	}
	return fcs.Run("git", gitArgs...)
}

// ContinueRebase continues the currently ongoing rebase.
func (fcs *FrontendCommands) ContinueRebase() error {
	return fcs.Run("git", "rebase", "--continue")
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (fcs *FrontendCommands) CreateBranch(name domain.LocalBranchName, parent domain.Location) error {
	return fcs.Run("git", "branch", name.String(), parent.String())
}

// DeleteLastCommit resets HEAD to the previous commit.
func (fcs *FrontendCommands) DeleteLastCommit() error {
	return fcs.Run("git", "reset", "--hard", "HEAD~1")
}

// PushBranch pushes the branch with the given name to origin.
func (fcs *FrontendCommands) CreateTrackingBranch(branch domain.LocalBranchName, remote domain.Remote, noPushHook bool) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, "-u", remote.String())
	args = append(args, branch.String())
	return fcs.Run("git", args...)
}

// DeleteLocalBranch removes the local branch with the given name.
func (fcs *FrontendCommands) DeleteLocalBranch(name domain.LocalBranchName, force bool) error {
	args := []string{"branch", "-d", name.String()}
	if force {
		args[1] = "-D"
	}
	return fcs.Run("git", args...)
}

// DeleteRemoteBranch removes the remote branch of the given local branch.
func (fcs *FrontendCommands) DeleteRemoteBranch(name domain.RemoteBranchName) error {
	remote, localBranchName := name.Parts()
	return fcs.Run("git", "push", remote.String(), ":"+localBranchName.String())
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (fcs *FrontendCommands) DiffParent(branch, parentBranch domain.LocalBranchName) error {
	return fcs.Run("git", "diff", parentBranch.String()+".."+branch.String())
}

// DiscardOpenChanges deletes all uncommitted changes.
func (fcs *FrontendCommands) DiscardOpenChanges() error {
	return fcs.Run("git", "reset", "--hard")
}

// Fetch retrieves the updates from the origin repo.
func (fcs *FrontendCommands) Fetch() error {
	return fcs.Run("git", "fetch", "--prune", "--tags")
}

// FetchUpstream fetches updates from the upstream remote.
func (fcs *FrontendCommands) FetchUpstream(branch domain.LocalBranchName) error {
	return fcs.Run("git", "fetch", domain.UpstreamRemote.String(), branch.String())
}

// MergeBranchNoEdit merges the given branch into the current branch,
// using the default commit message.
func (fcs *FrontendCommands) MergeBranchNoEdit(branch domain.BranchName) error {
	err := fcs.Run("git", "merge", "--no-edit", branch.String())
	return err
}

// NavigateToDir changes into the root directory of the current repository.
func (fcs *FrontendCommands) NavigateToDir(dir domain.RepoRootDir) error {
	return os.Chdir(dir.String())
}

// PopStash restores stashed-away changes into the workspace.
func (fcs *FrontendCommands) PopStash() error {
	return fcs.Run("git", "stash", "pop")
}

// Pull fetches updates from origin and updates the currently checked out branch.
func (fcs *FrontendCommands) Pull() error {
	return fcs.Run("git", "pull")
}

// PushCurrentBranch pushes the current branch to its tracking branch.
func (fcs *FrontendCommands) PushCurrentBranch(noPushHook bool) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return fcs.Run("git", args...)
}

// PushBranch pushes the branch with the given name to origin.
func (fcs *FrontendCommands) ForcePushBranch(noPushHook bool) error {
	args := []string{"push", "--force-with-lease"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return fcs.Run("git", args...)
}

// PushTags pushes new the Git tags to origin.
func (fcs *FrontendCommands) PushTags() error {
	return fcs.Run("git", "push", "--tags")
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (fcs *FrontendCommands) Rebase(target domain.BranchName) error {
	return fcs.Run("git", "rebase", target.String())
}

// RemoveGitAlias removes the given Git alias.
func (fcs *FrontendCommands) RemoveGitAlias(alias config.Alias) error {
	return fcs.Run("git", "config", "--global", "--unset", "alias."+alias.String())
}

// ResetCurrentBranchToSHA undoes all commits on the current branch all the way until the given SHA.
func (fcs *FrontendCommands) ResetCurrentBranchToSHA(sha domain.SHA, hard bool) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	args = append(args, sha.String())
	return fcs.Run("git", args...)
}

// ResetRemoteBranchToSHA sets the given remote branch to the given SHA.
func (fcs *FrontendCommands) ResetRemoteBranchToSHA(branch domain.RemoteBranchName, sha domain.SHA) error {
	return fcs.Run("git", "push", "--force-with-lease", domain.OriginRemote.String(), sha.String()+":"+branch.LocalBranchName().String())
}

// RevertCommit reverts the commit with the given SHA.
func (fcs *FrontendCommands) RevertCommit(sha domain.SHA) error {
	return fcs.Run("git", "revert", sha.String())
}

// SquashMerge squash-merges the given branch into the current branch.
func (fcs *FrontendCommands) SquashMerge(branch domain.LocalBranchName) error {
	return fcs.Run("git", "merge", "--squash", branch.String())
}

// Stash adds the current files to the Git stash.
func (fcs *FrontendCommands) Stash() error {
	return fcs.RunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
}

// StageFiles adds the file with the given name to the Git index.
func (fcs *FrontendCommands) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	return fcs.Run("git", args...)
}

// StartCommit starts a commit and stops at asking the user for the commit message.
func (fcs *FrontendCommands) StartCommit() error {
	return fcs.Run("git", "commit")
}

func (fcs *FrontendCommands) UndoLastCommit() error {
	return fcs.Run("git", "reset", "--soft", "HEAD^")
}
