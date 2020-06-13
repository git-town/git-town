package git

import "github.com/git-town/git-town/src/command"

// CachedRunner is a git.Runner that caches parts of the current state of its Git repository in memory
// to avoid unnecessary Git operations.
type CachedRunner struct {
	Runner
	remoteBranchCache *RemoteBranchCache // caches the remote branches of this Git repo
}

// NewCachedRunner provides CachedRunner instances.
func NewCachedRunner(shell command.Shell, config *Configuration, remotesCache *RemotesCache, remoteBranchCache *RemoteBranchCache) CachedRunner {
	return CachedRunner{NewRunner(shell, config, remotesCache), remoteBranchCache}
}

// RemoteBranches provides the names of the remote branches in this repo.
func (cr *CachedRunner) RemoteBranches() ([]string, error) {
	if cr.remoteBranchCache.Initialized() {
		return cr.remoteBranchCache.Get(), nil
	}
	result, err := cr.Runner.RemoteBranches()
	cr.remoteBranchCache.Set(result)
	return result, err
}
