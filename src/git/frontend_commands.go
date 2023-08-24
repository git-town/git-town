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
func (fc *FrontendCommands) AbortMerge() error {
	return fc.Run("git", "merge", "--abort")
}

// AbortRebase cancels a currently ongoing Git rebase operation.
func (fc *FrontendCommands) AbortRebase() error {
	return fc.Run("git", "rebase", "--abort")
}

// AddGitAlias sets the given Git alias.
func (fc *FrontendCommands) AddGitAlias(alias config.Alias) error {
	aliasKey := config.NewAliasKey(alias)
	return fc.Run("git", "config", "--global", aliasKey.String(), "town "+alias.String())
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (fc *FrontendCommands) CheckoutBranch(name domain.LocalBranchName) error {
	err := fc.Run("git", "checkout", name.String())
	if err != nil {
		return fmt.Errorf(messages.BranchCheckoutProblem, name, err)
	}
	fc.SetCachedCurrentBranch(name)
	return nil
}

// CreateRemoteBranch creates a remote branch from the given local SHA.
func (fc *FrontendCommands) CreateRemoteBranch(localSHA domain.SHA, branch domain.LocalBranchName, noPushHook bool) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, domain.OriginRemote.String(), localSHA.String()+":refs/heads/"+branch.String())
	return fc.Run("git", args...)
}

// CommitNoEdit commits all staged files with the default commit message.
func (fc *FrontendCommands) CommitNoEdit() error {
	return fc.Run("git", "commit", "--no-edit")
}

// CommitStagedChanges commits the currently staged changes.
func (fc *FrontendCommands) CommitStagedChanges(message string) error {
	if message != "" {
		return fc.Run("git", "commit", "-m", message)
	}
	return fc.Run("git", "commit", "--no-edit")
}

// Commit performs a commit of the staged changes with an optional custom message and author.
func (fc *FrontendCommands) Commit(message, author string) error {
	gitArgs := []string{"commit"}
	if message != "" {
		gitArgs = append(gitArgs, "-m", message)
	}
	if author != "" {
		gitArgs = append(gitArgs, "--author", author)
	}
	return fc.Run("git", gitArgs...)
}

// ContinueRebase continues the currently ongoing rebase.
func (fc *FrontendCommands) ContinueRebase() error {
	return fc.Run("git", "rebase", "--continue")
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (fc *FrontendCommands) CreateBranch(name domain.LocalBranchName, parent domain.Location) error {
	return fc.Run("git", "branch", name.String(), parent.String())
}

// DeleteLastCommit resets HEAD to the previous commit.
func (fc *FrontendCommands) DeleteLastCommit() error {
	return fc.Run("git", "reset", "--hard", "HEAD~1")
}

// PushBranch pushes the branch with the given name to origin.
func (fc *FrontendCommands) CreateTrackingBranch(branch domain.LocalBranchName, remote domain.Remote, noPushHook bool) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, "-u", remote.String())
	args = append(args, branch.String())
	return fc.Run("git", args...)
}

// DeleteLocalBranch removes the local branch with the given name.
func (fc *FrontendCommands) DeleteLocalBranch(name domain.LocalBranchName, force bool) error {
	args := []string{"branch", "-d", name.String()}
	if force {
		args[1] = "-D"
	}
	return fc.Run("git", args...)
}

// DeleteRemoteBranch removes the remote branch of the given local branch.
// TODO: provide the actual domain.RemoteBranchName here and delete that branch instead of "origin/<localbranch>".
func (fc *FrontendCommands) DeleteRemoteBranch(name domain.LocalBranchName) error {
	return fc.Run("git", "push", domain.OriginRemote.String(), ":"+name.String())
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (fc *FrontendCommands) DiffParent(branch, parentBranch domain.LocalBranchName) error {
	return fc.Run("git", "diff", parentBranch.String()+".."+branch.String())
}

// DiscardOpenChanges deletes all uncommitted changes.
func (fc *FrontendCommands) DiscardOpenChanges() error {
	return fc.Run("git", "reset", "--hard")
}

// Fetch retrieves the updates from the origin repo.
func (fc *FrontendCommands) Fetch() error {
	return fc.Run("git", "fetch", "--prune", "--tags")
}

// FetchUpstream fetches updates from the upstream remote.
func (fc *FrontendCommands) FetchUpstream(branch domain.LocalBranchName) error {
	return fc.Run("git", "fetch", domain.UpstreamRemote.String(), branch.String())
}

// MergeBranchNoEdit merges the given branch into the current branch,
// using the default commit message.
func (fc *FrontendCommands) MergeBranchNoEdit(branch domain.BranchName) error {
	err := fc.Run("git", "merge", "--no-edit", branch.String())
	return err
}

// NavigateToDir changes into the root directory of the current repository.
func (fc *FrontendCommands) NavigateToDir(dir string) error {
	return os.Chdir(dir)
}

// PopStash restores stashed-away changes into the workspace.
func (fc *FrontendCommands) PopStash() error {
	return fc.Run("git", "stash", "pop")
}

// Pull fetches updates from origin and updates the currently checked out branch.
func (fc *FrontendCommands) Pull() error {
	return fc.Run("git", "pull")
}

// PushCurrentBranch pushes the current branch to its tracking branch.
func (fc *FrontendCommands) PushCurrentBranch(noPushHook bool) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return fc.Run("git", args...)
}

// PushBranch pushes the branch with the given name to origin.
func (fc *FrontendCommands) ForcePushBranch(noPushHook bool) error {
	args := []string{"push", "--force-with-lease"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	return fc.Run("git", args...)
}

// PushTags pushes new the Git tags to origin.
func (fc *FrontendCommands) PushTags() error {
	return fc.Run("git", "push", "--tags")
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (fc *FrontendCommands) Rebase(target domain.BranchName) error {
	return fc.Run("git", "rebase", target.String())
}

// RemoveGitAlias removes the given Git alias.
func (fc *FrontendCommands) RemoveGitAlias(alias config.Alias) error {
	return fc.Run("git", "config", "--global", "--unset", "alias."+alias.String())
}

// ResetCurrentBranchToSHA undoes all commits on the current branch all the way until the given SHA.
func (fc *FrontendCommands) ResetCurrentBranchToSHA(sha domain.SHA, hard bool) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	args = append(args, sha.String())
	return fc.Run("git", args...)
}

// ResetCurrentBranchToSHA undoes all commits on the current branch all the way until the given SHA.
func (fc *FrontendCommands) ResetRemoteBranchToSHA(remoteBranch domain.RemoteBranchName, shaToPush domain.SHA, shaThatMustExist domain.SHA) error {
	remote, localBranch := remoteBranch.Parts()
	return fc.Run("git", "push", fmt.Sprintf("--force-with-lease=%s:%s", localBranch.String(), shaThatMustExist.String()), remote.String(), shaToPush.String()+":"+localBranch.String())
}

// RevertCommit reverts the commit with the given SHA.
func (fc *FrontendCommands) RevertCommit(sha domain.SHA) error {
	return fc.Run("git", "revert", sha.String())
}

// SquashMerge squash-merges the given branch into the current branch.
func (fc *FrontendCommands) SquashMerge(branch domain.LocalBranchName) error {
	return fc.Run("git", "merge", "--squash", branch.String())
}

// Stash adds the current files to the Git stash.
func (fc *FrontendCommands) Stash() error {
	return fc.RunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
}

// StageFiles adds the file with the given name to the Git index.
func (fc *FrontendCommands) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	return fc.Run("git", args...)
}

// StartCommit starts a commit and stops at asking the user for the commit message.
func (fc *FrontendCommands) StartCommit() error {
	return fc.Run("git", "commit")
}
