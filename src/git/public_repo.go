package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/stringslice"
)

type PublicRunner interface {
	Run(executable string, args ...string) error
	RunMany([][]string) error
}

type PublicRepo struct {
	Public PublicRunner
	InternalRepo
}

// AbortMerge cancels a currently ongoing Git merge operation.
func (r *PublicRepo) AbortMerge() error {
	err := r.Public.Run("git", "merge", "--abort")
	if err != nil {
		return fmt.Errorf("cannot abort current merge: %w", err)
	}
	return nil
}

// AbortRebase cancels a currently ongoing Git rebase operation.
func (r *PublicRepo) AbortRebase() error {
	err := r.Public.Run("git", "rebase", "--abort")
	if err != nil {
		return fmt.Errorf("cannot abort current merge: %w", err)
	}
	return nil
}

// AddGitAlias sets the given Git alias.
func (r *PublicRepo) AddGitAlias(aliasType config.AliasType) error {
	return r.Public.Run("git", "config", "--global", "alias."+string(aliasType), "town "+string(aliasType))
}

// CheckoutBranch checks out the Git branch with the given name in this repo.
func (r *PublicRepo) CheckoutBranch(name string) error {
	err := r.Public.Run("git", "checkout", name)
	if err != nil {
		return fmt.Errorf("cannot check out branch %q: %w", name, err)
	}
	// if name != "-" {
	// 	r.CurrentBranchCache.Set(name)
	// } else {
	// 	r.CurrentBranchCache.Invalidate()
	// }
	return nil
}

// CreateRemoteBranch creates a remote branch from the given local SHA.
func (r *PublicRepo) CreateRemoteBranch(localSha, branch string, noPushHook bool) error {
	args := []string{"push"}
	if noPushHook {
		args = append(args, "--no-verify")
	}
	args = append(args, config.OriginRemote, localSha+":refs/heads/"+branch)
	err := r.Public.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot create remote branch for local SHA %q: %w", localSha, err)
	}
	return nil
}

// CommitNoEdit commits all staged files with the default commit message.
func (r *PublicRepo) CommitNoEdit() error {
	err := r.Public.Run("git", "commit", "--no-edit")
	if err != nil {
		return fmt.Errorf("cannot commit files: %w", err)
	}
	return nil
}

// CommitStagedChanges commits the currently staged changes.
func (r *PublicRepo) CommitStagedChanges(message string) error {
	var err error
	if message != "" {
		err = r.Public.Run("git", "commit", "-m", message)
	} else {
		err = r.Public.Run("git", "commit", "--no-edit")
	}
	if err != nil {
		return fmt.Errorf("cannot commit staged changes: %w", err)
	}
	return nil
}

// Commit performs a commit of the staged changes with an optional custom message and author.
func (r *PublicRepo) Commit(message, author string) error {
	gitArgs := []string{"commit"}
	if message != "" {
		gitArgs = append(gitArgs, "-m", message)
	}
	if author != "" {
		gitArgs = append(gitArgs, "--author", author)
	}
	err := r.Public.Run("git", gitArgs...)
	return err
}

// ContinueRebase continues the currently ongoing rebase.
func (r *PublicRepo) ContinueRebase() error {
	err := r.Public.Run("git", "rebase", "--continue")
	if err != nil {
		return fmt.Errorf("cannot continue rebase: %w", err)
	}
	return nil
}

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func (r *PublicRepo) CreateBranch(name, parent string) error {
	err := r.Public.Run("git", "branch", name, parent)
	if err != nil {
		return fmt.Errorf("cannot create branch %q: %w", name, err)
	}
	return nil
}

// DeleteLastCommit resets HEAD to the previous commit.
func (r *PublicRepo) DeleteLastCommit() error {
	err := r.Public.Run("git", "reset", "--hard", "HEAD~1")
	if err != nil {
		return fmt.Errorf("cannot delete last commit: %w", err)
	}
	return nil
}

// DeleteLocalBranch removes the local branch with the given name.
func (r *PublicRepo) DeleteLocalBranch(name string, force bool) error {
	args := []string{"branch", "-d", name}
	if force {
		args[1] = "-D"
	}
	err := r.Public.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot delete local branch %q: %w", name, err)
	}
	return nil
}

// DeleteRemoteBranch removes the remote branch of the given local branch.
func (r *PublicRepo) DeleteRemoteBranch(name string) error {
	err := r.Public.Run("git", "push", config.OriginRemote, ":"+name)
	if err != nil {
		return fmt.Errorf("cannot delete tracking branch for %q: %w", name, err)
	}
	return nil
}

// DiffParent displays the diff between the given branch and its given parent branch.
func (r *PublicRepo) DiffParent(branch, parentBranch string) error {
	err := r.Public.Run("git", "diff", parentBranch+".."+branch)
	if err != nil {
		return fmt.Errorf("cannot diff branch %q with its parent branch %q: %w", branch, parentBranch, err)
	}
	return nil
}

// DiscardOpenChanges deletes all uncommitted changes.
func (r *PublicRepo) DiscardOpenChanges() error {
	err := r.Public.Run("git", "reset", "--hard")
	if err != nil {
		return fmt.Errorf("cannot discard open changes: %w", err)
	}
	return nil
}

// Fetch retrieves the updates from the origin repo.
func (r *PublicRepo) Fetch() error {
	err := r.Public.Run("git", "fetch", "--prune", "--tags")
	if err != nil {
		return fmt.Errorf("cannot fetch: %w", err)
	}
	return nil
}

// FetchUpstream fetches updates from the upstream remote.
func (r *PublicRepo) FetchUpstream(branch string) error {
	err := r.Public.Run("git", "fetch", "upstream", branch)
	if err != nil {
		return fmt.Errorf("cannot fetch from upstream: %w", err)
	}
	return nil
}

// MergeBranchNoEdit merges the given branch into the current branch,
// using the default commit message.
func (r *PublicRepo) MergeBranchNoEdit(branch string) error {
	err := r.Public.Run("git", "merge", "--no-edit", branch)
	return err
}

// NavigateToRootIfNecessary changes into the root directory of the current repository.
func (r *PublicRepo) NavigateToRootIfNecessary() error {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current working directory: %w", err)
	}
	gitRootDirectory, err := r.RootDirectory()
	if err != nil {
		return err
	}
	if currentDirectory == gitRootDirectory {
		return nil
	}
	return os.Chdir(gitRootDirectory)
}

// PopStash restores stashed-away changes into the workspace.
func (r *PublicRepo) PopStash() error {
	err := r.Public.Run("git", "stash", "pop")
	if err != nil {
		return fmt.Errorf("cannot pop the stash: %w", err)
	}
	return nil
}

// Pull fetches updates from origin and updates the currently checked out branch.
func (r *PublicRepo) Pull() error {
	err := r.Public.Run("git", "pull")
	if err != nil {
		return fmt.Errorf("cannot pull updates: %w", err)
	}
	return nil
}

type PushArgs struct {
	Branch         string
	Force          bool `exhaustruct:"optional"`
	ForceWithLease bool `exhaustruct:"optional"`
	NoPushHook     bool `exhaustruct:"optional"`
	Remote         string
}

// PushBranch pushes the branch with the given name to origin.
// TODO: remove unused elements from PushArgs.
func (r *PublicRepo) PushBranch(options ...PushArgs) error {
	var option PushArgs
	if len(options) > 0 {
		option = options[0]
	} else {
		option = PushArgs{} //nolint:exhaustruct  // intentional zero-value object
	}
	args := []string{"push"}
	provideBranch := false
	if option.Force {
		args = append(args, "-f")
	}
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
	err := r.Public.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot push branch to origin: %w", err)
	}
	return nil
}

// PushTags pushes new the Git tags to origin.
func (r *PublicRepo) PushTags() error {
	err := r.Public.Run("git", "push", "--tags")
	if err != nil {
		return fmt.Errorf("cannot push branch in repo: %w", err)
	}
	return nil
}

// Rebase initiates a Git rebase of the current branch against the given branch.
func (r *PublicRepo) Rebase(target string) error {
	err := r.Public.Run("git", "rebase", target)
	if err != nil {
		return fmt.Errorf("cannot rebase against branch %q: %w", target, err)
	}
	return nil
}

// RemoveGitAlias removes the given Git alias.
func (r *PublicRepo) RemoveGitAlias(aliasType config.AliasType) error {
	return r.Public.Run("git", "config", "--global", "--unset", "alias."+string(aliasType), "town "+string(aliasType))
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (r *PublicRepo) RemoveOutdatedConfiguration() error {
	branches, err := r.LocalAndOriginBranches(r.Config.MainBranch())
	if err != nil {
		return err
	}
	for child, parent := range r.Config.ParentBranchMap() {
		hasChildBranch := stringslice.Contains(branches, child)
		hasParentBranch := stringslice.Contains(branches, parent)
		if !hasChildBranch || !hasParentBranch {
			err = r.Config.RemoveParentBranch(child)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ResetToSha undoes all commits on the current branch all the way until the given SHA.
func (r *PublicRepo) ResetToSha(sha string, hard bool) error {
	args := []string{"reset"}
	if hard {
		args = append(args, "--hard")
	}
	args = append(args, sha)
	err := r.Public.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot reset to SHA %q: %w", sha, err)
	}
	return nil
}

// RevertCommit reverts the commit with the given SHA.
func (r *PublicRepo) RevertCommit(sha string) error {
	err := r.Public.Run("git", "revert", sha)
	if err != nil {
		return fmt.Errorf("cannot revert commit %q: %w", sha, err)
	}
	return nil
}

// SquashMerge squash-merges the given branch into the current branch.
func (r *PublicRepo) SquashMerge(branch string) error {
	err := r.Public.Run("git", "merge", "--squash", branch)
	if err != nil {
		return fmt.Errorf("cannot squash-merge branch %q: %w", branch, err)
	}
	return nil
}

// Stash adds the current files to the Git stash.
func (r *PublicRepo) Stash() error {
	err := r.Public.RunMany([][]string{
		{"git", "add", "-A"},
		{"git", "stash"},
	})
	if err != nil {
		return fmt.Errorf("cannot stash: %w", err)
	}
	return nil
}

// StageFiles adds the file with the given name to the Git index.
func (r *PublicRepo) StageFiles(names ...string) error {
	args := append([]string{"add"}, names...)
	err := r.Public.Run("git", args...)
	if err != nil {
		return fmt.Errorf("cannot stage files %s: %w", strings.Join(names, ", "), err)
	}
	return nil
}

// StartCommit starts a commit and stops at asking the user for the commit message.
func (r *PublicRepo) StartCommit() error {
	err := r.Public.Run("git", "commit")
	if err != nil {
		return fmt.Errorf("cannot start commit: %w", err)
	}
	return nil
}
