package gitconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

func RemoveBranchTypeOverride(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) error {
	key := configdomain.NewBranchTypeOverrideKeyForBranch(branch)
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, key.Key)
}

func RemoveDeprecatedCreatePrototypeBranches(runner subshelldomain.Runner) {
	_ = RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.KeyDeprecatedCreatePrototypeBranches)
}

func RemoveDevRemote(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyDevRemote)
}

func RemoveFeatureRegex(runner subshelldomain.Runner, scope configdomain.ConfigScope) error {
	return RemoveConfigValue(runner, scope, configdomain.KeyFeatureRegex)
}

func RemoveParent(runner subshelldomain.Runner, child gitdomain.LocalBranchName) error {
	return RemoveConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(child))
}

func SetParent(runner subshelldomain.Runner, child, parent gitdomain.LocalBranchName) error {
	return SetConfigValue(runner, configdomain.ConfigScopeLocal, configdomain.NewParentKey(child), parent.String())
}
