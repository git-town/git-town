package git

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/slice"
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

type SetCachedCurrentBranchFunc func(string)

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
func (fc *FrontendCommands) CheckoutBranch(name string) error {
	err := fc.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf(messages.BranchCheckoutProblem, name, err)
	}
	fc.SetCachedCurrentBranch(name)
	return nil
}

// CreateRemoteBranch creates a remote branch from the given local SHA.
func (fc *FrontendCommands) CreateRemoteBranch(localSha SHA, branch string, noPushHook bool) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, config.OriginRemote, localSha.Content+":refs/heads/"+branch)
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
func (fc *FrontendCommands) CreateBranch(name, parent string) error {
	return fc.Run("git", "branch", name, parent)
}

// DeleteLastCommit resets HEAD to the previous commit.
func (fc *FrontendCommands) DeleteLastCommit() error {
	return fc.Run("git", "reset", "--hard", "HEAD~1")
}

// DeleteLocalBranch removes the local branch with the given name.
func (fc *FrontendCommands) DeleteLocalBranch(name string, force bool) error {
	args := []string{"branch", "-d", name}
	if force {
		args[1] = "-D"
	}
	return fc.Run("git", args...)
}

// DeleteRemoteBranch removes the remote branch of the given local branch.
func (fc *FrontendCommands) DeleteRemoteBranch(name string) error {
	return fc.Run("git", "push", config.OriginRemote, ":"+name)
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (fc *FrontendCommands) DiffParent(branch, parentBranch string) error {
	return fc.Run("git", "diff", parentBranch+".."+branch)
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
func (fc *FrontendCommands) FetchUpstream(branch string) error {
	return fc.Run("git", "fetch", "upstream", branch)
}

// MergeBranchNoEdit merges the given branch into the current branch,
// using the default commit message.
func (fc *FrontendCommands) MergeBranchNoEdit(branch string) error {
	err := fc.Run("git", "merge", "--no-edit", branch)
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

type PushArgs struct {
	Branch         string
	ForceWithLease bool `exhaustruct:"optional"`
	NoPushHook     bool `exhaustruct:"optional"`
	Remote         string
}

// PushBranch pushes the branch with the given name to origin.
func (fc *FrontendCommands) PushBranch(options ...PushArgs) error {
	option := slice.FirstElementOr(options, PushArgs{Branch: "", Remote: ""})
	args := []string{"push"}
	provideBranch := false
	if option.NoPushHook {
		args = append(args, "--no-verify")
	}
	if option.ForceWithLease {
		args = append(args, "--force-with-lease")
	}
	if option.Remote != "" {
		args = append(args, "-u", option.Remote)
		provideBranch = true
	}
	if option.Branch != "" && provideBranch {
		args = append(args, option.Branch)
	}
	return fc.Run("git", args...)
}

// PushTags pushes new the Git tags to origin.
func (fc *FrontendCommands) PushTags() error {
	return fc.Run("git", "push", "--tags")
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (fc *FrontendCommands) Rebase(target string) error {
	return fc.Run("git", "rebase", target)
}

// RemoveGitAlias removes the given Git alias.
func (fc *FrontendCommands) RemoveGitAlias(alias config.Alias) error {
	return fc.Run("git", "config", "--global", "--unset", "alias."+alias.String())
}

// ResetToSha undoes all commits on the current branch all the way until the given SHA.
func (fc *FrontendCommands) ResetToSha(sha SHA, hard bool) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	args = append(args, sha.Content)
	return fc.Run("git", args...)
}

// RevertCommit reverts the commit with the given SHA.
func (fc *FrontendCommands) RevertCommit(sha SHA) error {
	return fc.Run("git", "revert", sha.Content)
}

// SquashMerge squash-merges the given branch into the current branch.
func (fc *FrontendCommands) SquashMerge(branch string) error {
	return fc.Run("git", "merge", "--squash", branch)
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
