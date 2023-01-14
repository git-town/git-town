// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/git-town/git-town/v7/src/stringslice"
)

// GitTown provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type GitTown struct {
	Storage Git
}

func NewGitTown(shell run.Shell) GitTown {
	return GitTown{
		Storage: NewGit(shell),
	}
}

// AddGitAlias sets the given Git alias.
func (gt *GitTown) AddGitAlias(command string) (*run.Result, error) {
	return gt.Storage.SetGlobalConfigValue("alias."+command, "town "+command)
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (gt *GitTown) AddToPerennialBranches(branchNames ...string) error {
	return gt.SetPerennialBranches(append(gt.PerennialBranches(), branchNames...))
}

// AncestorBranches provides the names of all parent branches for the given branch,
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func (gt *GitTown) AncestorBranches(branchName string) []string {
	parentBranchMap := gt.ParentBranchMap()
	current := branchName
	result := []string{}
	for {
		if gt.IsMainBranch(current) || gt.IsPerennialBranch(current) {
			return result
		}
		parent := parentBranchMap[current]
		if parent == "" {
			return result
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// BranchAncestryRoots provides the branches with children and no parents.
func (gt *GitTown) BranchAncestryRoots() []string {
	parentMap := gt.ParentBranchMap()
	roots := []string{}
	for _, parent := range parentMap {
		if _, ok := parentMap[parent]; !ok && !stringslice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}

// ChildBranches provides the names of all branches for which the given branch
// is a parent.
func (gt *GitTown) ChildBranches(branchName string) []string {
	result := []string{}
	for _, key := range gt.Storage.LocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		parent := gt.Storage.LocalConfigValue(key)
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	return result
}

func (gt *GitTown) DeprecatedNewBranchPushFlagGlobal() string {
	return gt.Storage.globalConfigCache["git-town.new-branch-push-flag"]
}

func (gt *GitTown) DeprecatedNewBranchPushFlagLocal() string {
	return gt.Storage.localConfigCache["git-town.new-branch-push-flag"]
}

// GitAlias provides the currently set alias for the given Git Town command.
func (gt *GitTown) GitAlias(command string) string {
	return gt.Storage.GlobalConfigValue("alias." + command)
}

// GitHubToken provides the content of the GitHub API token stored in the local or global Git Town configuration.
func (gt *GitTown) GitHubToken() string {
	return gt.Storage.LocalOrGlobalConfigValue("git-town.github-token")
}

// GitLabToken provides the content of the GitLab API token stored in the local or global Git Town configuration.
func (gt *GitTown) GitLabToken() string {
	return gt.Storage.LocalOrGlobalConfigValue("git-town.gitlab-token")
}

// GiteaToken provides the content of the Gitea API token stored in the local or global Git Town configuration.
func (gt *GitTown) GiteaToken() string {
	return gt.Storage.LocalOrGlobalConfigValue("git-town.gitea-token")
}

// HasBranchInformation indicates whether this configuration contains any branch hierarchy entries.
func (gt *GitTown) HasBranchInformation() bool {
	for key := range gt.Storage.localConfigCache {
		if strings.HasPrefix(key, "git-town-branch.") {
			return true
		}
	}
	return false
}

// HasParentBranch returns whether or not the given branch has a parent.
func (gt *GitTown) HasParentBranch(branchName string) bool {
	return gt.ParentBranch(branchName) != ""
}

// HostingService provides the name of the code hosting driver to use.
func (gt *GitTown) HostingService() string {
	return gt.Storage.LocalOrGlobalConfigValue("git-town.code-hosting-driver")
}

// IsAncestorBranch indicates whether the given branch is an ancestor of the other given branch.
func (gt *GitTown) IsAncestorBranch(branchName, ancestorBranchName string) bool {
	ancestorBranches := gt.AncestorBranches(branchName)
	return stringslice.Contains(ancestorBranches, ancestorBranchName)
}

// IsFeatureBranch indicates whether the branch with the given name is
// a feature branch.
func (gt *GitTown) IsFeatureBranch(branchName string) bool {
	return !gt.IsMainBranch(branchName) && !gt.IsPerennialBranch(branchName)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (gt *GitTown) IsMainBranch(branchName string) bool {
	return branchName == gt.MainBranch()
}

// IsOffline indicates whether Git Town is currently in offline mode.
func (gt *GitTown) IsOffline() bool {
	config := gt.Storage.GlobalConfigValue("git-town.offline")
	if config == "" {
		return false
	}
	result, err := strconv.ParseBool(config)
	if err != nil {
		fmt.Printf("Invalid value for git-town.offline: %q. Please provide either \"yes\" or \"no\". Considering \"no\" for now.", config)
		fmt.Println()
		return false
	}
	return result
}

// IsPerennialBranch indicates whether the branch with the given name is
// a perennial branch.
func (gt *GitTown) IsPerennialBranch(branchName string) bool {
	perennialBranches := gt.PerennialBranches()
	return stringslice.Contains(perennialBranches, branchName)
}

// MainBranch provides the name of the main branch.
func (gt *GitTown) MainBranch() string {
	return gt.Storage.LocalOrGlobalConfigValue("git-town.main-branch-name")
}

// OriginOverride provides the override for the origin hostname from the Git Town configuration.
func (gt *GitTown) OriginOverride() string {
	return gt.Storage.LocalConfigValue("git-town.code-hosting-origin-hostname")
}

// OriginURL provides the URL for the "origin" remote.
// In tests this value can be stubbed.
func (gt *GitTown) OriginURL() string {
	remote := os.Getenv("GIT_TOWN_REMOTE")
	if remote != "" {
		return remote
	}
	res, _ := gt.Storage.shell.Run("git", "remote", "get-url", "origin")
	return res.OutputSanitized()
}

// ParentBranchMap returns a map from branch name to its parent branch.
func (gt *GitTown) ParentBranchMap() map[string]string {
	result := map[string]string{}
	for _, key := range gt.Storage.LocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := gt.Storage.LocalConfigValue(key)
		result[child] = parent
	}
	return result
}

// ParentBranch provides the name of the parent branch of the given branch.
func (gt *GitTown) ParentBranch(branchName string) string {
	return gt.Storage.LocalConfigValue("git-town-branch." + branchName + ".parent")
}

// PerennialBranches returns all branches that are marked as perennial.
func (gt *GitTown) PerennialBranches() []string {
	result := gt.Storage.LocalOrGlobalConfigValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// PullBranchStrategy provides the currently configured pull branch strategy.
func (gt *GitTown) PullBranchStrategy() string {
	config := gt.Storage.LocalOrGlobalConfigValue("git-town.pull-branch-strategy")
	if config != "" {
		return config
	}
	return "rebase"
}

// PushHook provides the currently configured push-hook setting.
func (gt *GitTown) PushHook() (bool, error) {
	setting := gt.Storage.LocalOrGlobalConfigValue("git-town.push-hook")
	if setting == "" {
		return true, nil
	}
	result, err := cli.ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf("invalid value for git-town.push-hook: %q. Please provide either \"true\" or \"false\"", setting)
	}
	return result, nil
}

// PushHook provides the currently configured push-hook setting.
func (gt *GitTown) PushHookGlobal() (bool, error) {
	setting := gt.Storage.GlobalConfigValue("git-town.push-hook")
	if setting == "" {
		return true, nil
	}
	result, err := cli.ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf("invalid value for git-town.push-hook: %q. Please provide either \"true\" or \"false\"", setting)
	}
	return result, nil
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (gt *GitTown) RemoveFromPerennialBranches(branchName string) error {
	return gt.SetPerennialBranches(stringslice.Remove(gt.PerennialBranches(), branchName))
}

// RemoveGitAlias removes the given Git alias.
func (gt *GitTown) RemoveGitAlias(command string) (*run.Result, error) {
	return gt.Storage.removeGlobalConfigValue("alias." + command)
}

// RemoveLocalGitConfiguration removes all Git Town configuration.
func (gt *GitTown) RemoveLocalGitConfiguration() error {
	result, err := gt.Storage.shell.Run("git", "config", "--remove-section", "git-town")
	if err != nil {
		if result.ExitCode() == 128 {
			// Git returns exit code 128 when trying to delete a non-existing config section.
			// This is not an error condition in this workflow so we can ignore it here.
			return nil
		}
		return fmt.Errorf("unexpected error while removing the 'git-town' section from the Git configuration: %w", err)
	}
	return nil
}

// RemoveMainBranchConfiguration removes the configuration entry for the main branch name.
func (gt *GitTown) RemoveMainBranchConfiguration() error {
	return gt.Storage.removeLocalConfigValue("git-town.main-branch-name")
}

// RemoveParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func (gt *GitTown) RemoveParentBranch(branchName string) error {
	return gt.Storage.removeLocalConfigValue("git-town-branch." + branchName + ".parent")
}

// RemovePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (gt *GitTown) RemovePerennialBranchConfiguration() error {
	return gt.Storage.removeLocalConfigValue("git-town.perennial-branch-names")
}

// SetCodeHostingDriver sets the "github.code-hosting-driver" setting.
func (gt *GitTown) SetCodeHostingDriver(value string) error {
	const key = "git-town.code-hosting-driver"
	gt.Storage.localConfigCache[key] = value
	_, err := gt.Storage.shell.Run("git", "config", key, value)
	return err
}

// SetCodeHostingOriginHostname sets the "github.code-hosting-driver" setting.
func (gt *GitTown) SetCodeHostingOriginHostname(value string) error {
	const key = "git-town.code-hosting-origin-hostname"
	gt.Storage.localConfigCache[key] = value
	_, err := gt.Storage.shell.Run("git", "config", key, value)
	return err
}

// SetColorUI configures whether Git output contains color codes.
func (gt *GitTown) SetColorUI(value string) error {
	_, err := gt.Storage.shell.Run("git", "config", "color.ui", value)
	return err
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (gt *GitTown) SetMainBranch(branchName string) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town.main-branch-name", branchName)
	return err
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (gt *GitTown) SetNewBranchPush(value bool, global bool) error {
	setting := strconv.FormatBool(value)
	if global {
		_, err := gt.Storage.SetGlobalConfigValue("git-town.push-new-branches", setting)
		return err
	}
	_, err := gt.Storage.SetLocalConfigValue("git-town.push-new-branches", setting)
	return err
}

// SetOffline updates whether Git Town is in offline mode.
func (gt *GitTown) SetOffline(value bool) error {
	_, err := gt.Storage.SetGlobalConfigValue("git-town.offline", strconv.FormatBool(value))
	return err
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (gt *GitTown) SetParentBranch(branchName, parentBranchName string) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town-branch."+branchName+".parent", parentBranchName)
	return err
}

// SetPerennialBranches marks the given branches as perennial branches.
func (gt *GitTown) SetPerennialBranches(branchNames []string) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
	return err
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (gt *GitTown) SetPullBranchStrategy(strategy string) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town.pull-branch-strategy", strategy)
	return err
}

// SetPushHookLocally updates the configured pull branch strategy.
func (gt *GitTown) SetPushHookLocally(value bool) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town.push-hook", strconv.FormatBool(value))
	return err
}

// SetPushHook updates the configured pull branch strategy.
func (gt *GitTown) SetPushHookGlobally(value bool) error {
	_, err := gt.Storage.SetGlobalConfigValue("git-town.push-hook", strconv.FormatBool(value))
	return err
}

// SetShouldShipDeleteRemoteBranch updates the configured pull branch strategy.
func (gt *GitTown) SetShouldShipDeleteRemoteBranch(value bool) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town.ship-delete-remote-branch", strconv.FormatBool(value))
	return err
}

// SetShouldSyncUpstream updates the configured pull branch strategy.
func (gt *GitTown) SetShouldSyncUpstream(value bool) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town.sync-upstream", strconv.FormatBool(value))
	return err
}

func (gt *GitTown) SetSyncStrategy(value string) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town.sync-strategy", value)
	return err
}

// SetTestOrigin sets the origin to be used for testing.
func (gt *GitTown) SetTestOrigin(value string) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town.testing.remote-url", value)
	return err
}

// ShouldNewBranchPush indicates whether the current repository is configured to push
// freshly created branches up to origin.
func (gt *GitTown) ShouldNewBranchPush() (bool, error) {
	oldLocalConfig := gt.Storage.LocalConfigValue("git-town.new-branch-push-flag")
	if oldLocalConfig != "" {
		fmt.Println("I found the deprecated local setting \"git-town.new-branch-push-flag\".")
		fmt.Println("I am upgrading this setting to the new format \"git-town.push-new-branches\".")
		err := gt.Storage.removeLocalConfigValue("git-town.new-branch-push-flag")
		if err != nil {
			return false, err
		}
		parsed, err := cli.ParseBool(oldLocalConfig)
		if err != nil {
			return false, err
		}
		err = gt.SetNewBranchPush(parsed, false)
		if err != nil {
			return false, err
		}
	}
	oldGlobalConfig := gt.Storage.GlobalConfigValue("git-town.new-branch-push-flag")
	if oldGlobalConfig != "" {
		fmt.Println("I found the deprecated global setting \"git-town.new-branch-push-flag\".")
		fmt.Println("I am upgrading this setting to the new format \"git-town.push-new-branches\".")
		_, err := gt.Storage.removeGlobalConfigValue("git-town.new-branch-push-flag")
		if err != nil {
			return false, err
		}
		parsed, err := cli.ParseBool(oldGlobalConfig)
		if err != nil {
			return false, err
		}
		err = gt.SetNewBranchPush(parsed, true)
		if err != nil {
			return false, err
		}
	}
	config := gt.Storage.LocalOrGlobalConfigValue("git-town.push-new-branches")
	if config == "" {
		return false, nil
	}
	value, err := cli.ParseBool(config)
	if err != nil {
		return false, fmt.Errorf("invalid value for git-town.push-new-branches: %q. Please provide either \"yes\" or \"no\"", config)
	}
	return value, nil
}

// ShouldNewBranchPushGlobal indictes whether the global configuration requires to push
// freshly created branches to origin.
func (gt *GitTown) ShouldNewBranchPushGlobal() (bool, error) {
	config := gt.Storage.GlobalConfigValue("git-town.push-new-branches")
	if config == "" {
		return false, nil
	}
	return cli.ParseBool(config)
}

// ShouldShipDeleteOriginBranch indicates whether to delete the remote branch after shipping.
func (gt *GitTown) ShouldShipDeleteOriginBranch() bool {
	setting := gt.Storage.LocalOrGlobalConfigValue("git-town.ship-delete-remote-branch")
	if setting == "" {
		return true
	}
	result, err := strconv.ParseBool(setting)
	if err != nil {
		fmt.Printf("Invalid value for git-town.ship-delete-remote-branch: %q. Please provide either true or false. Considering true for now.\n", setting)
		return true
	}
	return result
}

// ShouldSyncUpstream indicates whether this repo should sync with its upstream.
func (gt *GitTown) ShouldSyncUpstream() bool {
	return gt.Storage.LocalOrGlobalConfigValue("git-town.sync-upstream") != "false"
}

func (gt *GitTown) SyncStrategy() string {
	setting := gt.Storage.LocalOrGlobalConfigValue("git-town.sync-strategy")
	if setting == "" {
		setting = "merge"
	}
	return setting
}

// ValidateIsOnline asserts that Git Town is not in offline mode.
func (gt *GitTown) ValidateIsOnline() error {
	if gt.IsOffline() {
		return errors.New("this command requires an active internet connection")
	}
	return nil
}
