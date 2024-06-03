package git

import (
	"os"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

type FrontendRunner interface {
	Run(executable string, args ...string) error
	RunMany(commands [][]string) error
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
