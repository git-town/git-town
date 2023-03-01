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

const (
	CodeHostingDriver         = "git-town.code-hosting-driver"
	CodeHostingOriginHostname = "git-town.code-hosting-origin-hostname"
	GiteaToken                = "git-town.gitea-token"  //nolint:gosec
	GithubToken               = "git-town.github-token" //nolint:gosec
	GitlabToken               = "git-town.gitlab-token" //nolint:gosec
	MainBranch                = "git-town.main-branch-name"
	NewBranchPushFlag         = "git-town.new-branch-push-flag"
	Offline                   = "git-town.offline"
	PerennialBranches         = "git-town.perennial-branch-names"
	PullBranchStrategy        = "git-town.pull-branch-strategy"
	PushHook                  = "git-town.push-hook"
	PushNewBranches           = "git-town.push-new-branches"
	ShipDeleteRemoteBranch    = "git-town.ship-delete-remote-branch"
	SyncUpstream              = "git-town.sync-upstream"
	SyncStrategy              = "git-town.sync-strategy"
	TestingRemoteURL          = "git-town.testing.remote-url"
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
func (gt *GitTown) AddToPerennialBranches(branches ...string) error {
	return gt.SetPerennialBranches(append(gt.PerennialBranches(), branches...))
}

// AncestorBranches provides the names of all parent branches for the given branch,
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func (gt *GitTown) AncestorBranches(branch string) []string {
	parentBranchMap := gt.ParentBranchMap()
	current := branch
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
func (gt *GitTown) ChildBranches(branch string) []string {
	result := []string{}
	for _, key := range gt.Storage.LocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		parent := gt.Storage.LocalConfigValue(key)
		if parent == branch {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	sort.Strings(result)
	return result
}

func (gt *GitTown) DeprecatedNewBranchPushFlagGlobal() string {
	return gt.Storage.globalConfigCache[NewBranchPushFlag]
}

func (gt *GitTown) DeprecatedNewBranchPushFlagLocal() string {
	return gt.Storage.localConfigCache[NewBranchPushFlag]
}

// GitAlias provides the currently set alias for the given Git Town command.
func (gt *GitTown) GitAlias(command string) string {
	return gt.Storage.GlobalConfigValue("alias." + command)
}

// GitHubToken provides the content of the GitHub API token stored in the local or global Git Town configuration.
func (gt *GitTown) GitHubToken() string {
	return gt.Storage.LocalOrGlobalConfigValue(GithubToken)
}

// GitLabToken provides the content of the GitLab API token stored in the local or global Git Town configuration.
func (gt *GitTown) GitLabToken() string {
	return gt.Storage.LocalOrGlobalConfigValue(GitlabToken)
}

// GiteaToken provides the content of the Gitea API token stored in the local or global Git Town configuration.
func (gt *GitTown) GiteaToken() string {
	return gt.Storage.LocalOrGlobalConfigValue(GiteaToken)
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
func (gt *GitTown) HasParentBranch(branch string) bool {
	return gt.ParentBranch(branch) != ""
}

// HostingService provides the name of the code hosting connector to use.
func (gt *GitTown) HostingService() string {
	return gt.Storage.LocalOrGlobalConfigValue(CodeHostingDriver)
}

// IsAncestorBranch indicates whether the given branch is an ancestor of the other given branch.
func (gt *GitTown) IsAncestorBranch(branch, ancestorBranch string) bool {
	ancestorBranches := gt.AncestorBranches(branch)
	return stringslice.Contains(ancestorBranches, ancestorBranch)
}

// IsFeatureBranch indicates whether the branch with the given name is
// a feature branch.
func (gt *GitTown) IsFeatureBranch(branch string) bool {
	return !gt.IsMainBranch(branch) && !gt.IsPerennialBranch(branch)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (gt *GitTown) IsMainBranch(branch string) bool {
	return branch == gt.MainBranch()
}

// IsOffline indicates whether Git Town is currently in offline mode.
func (gt *GitTown) IsOffline() (bool, error) {
	config := gt.Storage.GlobalConfigValue(Offline)
	if config == "" {
		return false, nil
	}
	result, err := cli.ParseBool(config)
	if err != nil {
		return false, fmt.Errorf("invalid value for %s: %q. Please provide either \"true\" or \"false\"", Offline, config)
	}
	return result, nil
}

// IsPerennialBranch indicates whether the branch with the given name is
// a perennial branch.
func (gt *GitTown) IsPerennialBranch(branch string) bool {
	perennialBranches := gt.PerennialBranches()
	return stringslice.Contains(perennialBranches, branch)
}

// MainBranch provides the name of the main branch.
func (gt *GitTown) MainBranch() string {
	return gt.Storage.LocalOrGlobalConfigValue(MainBranch)
}

// MainBranch provides the name of the main branch, or the given default value if none is configured.
func (gt *GitTown) MainBranchOr(defaultValue string) string {
	configured := gt.Storage.LocalOrGlobalConfigValue(MainBranch)
	if configured != "" {
		return configured
	}
	return defaultValue
}

// OriginOverride provides the override for the origin hostname from the Git Town configuration.
func (gt *GitTown) OriginOverride() string {
	return gt.Storage.LocalConfigValue(CodeHostingOriginHostname)
}

// OriginURL provides the URL for the "origin" remote.
// In tests this value can be stubbed.
func (gt *GitTown) OriginURL() string {
	remote := os.Getenv("GIT_TOWN_REMOTE")
	if remote != "" {
		return remote
	}
	res, _ := gt.Storage.shell.Run("git", "remote", "get-url", OriginRemote)
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
func (gt *GitTown) ParentBranch(branch string) string {
	return gt.Storage.LocalConfigValue("git-town-branch." + branch + ".parent")
}

// PerennialBranches returns all branches that are marked as perennial.
func (gt *GitTown) PerennialBranches() []string {
	result := gt.Storage.LocalOrGlobalConfigValue(PerennialBranches)
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// PullBranchStrategy provides the currently configured pull branch strategy.
func (gt *GitTown) PullBranchStrategy() string {
	config := gt.Storage.LocalOrGlobalConfigValue(PullBranchStrategy)
	if config != "" {
		return config
	}
	return "rebase"
}

// PushHook provides the currently configured push-hook setting.
func (gt *GitTown) PushHook() (bool, error) {
	setting := gt.Storage.LocalOrGlobalConfigValue(PushHook)
	if setting == "" {
		return true, nil
	}
	result, err := cli.ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf("invalid value for %s: %q. Please provide either \"true\" or \"false\"", PushHook, setting)
	}
	return result, nil
}

// PushHook provides the currently configured push-hook setting.
func (gt *GitTown) PushHookGlobal() (bool, error) {
	setting := gt.Storage.GlobalConfigValue(PushHook)
	if setting == "" {
		return true, nil
	}
	result, err := cli.ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf("invalid value for global %s: %q. Please provide either \"true\" or \"false\"", PushHook, setting)
	}
	return result, nil
}

// Reload refreshes the cached configuration data.
func (gt *GitTown) Reload() {
	gt.Storage.Reload()
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (gt *GitTown) RemoveFromPerennialBranches(branch string) error {
	return gt.SetPerennialBranches(stringslice.Remove(gt.PerennialBranches(), branch))
}

// RemoveGitAlias removes the given Git alias.
func (gt *GitTown) RemoveGitAlias(command string) (*run.Result, error) {
	return gt.Storage.RemoveGlobalConfigValue("alias." + command)
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
	return gt.Storage.RemoveLocalConfigValue(MainBranch)
}

// RemoveParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func (gt *GitTown) RemoveParentBranch(branch string) error {
	return gt.Storage.RemoveLocalConfigValue("git-town-branch." + branch + ".parent")
}

// RemovePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (gt *GitTown) RemovePerennialBranchConfiguration() error {
	return gt.Storage.RemoveLocalConfigValue(PerennialBranches)
}

// SetCodeHostingDriver sets the "github.code-hosting-driver" setting.
func (gt *GitTown) SetCodeHostingDriver(value string) error {
	gt.Storage.localConfigCache[CodeHostingDriver] = value
	_, err := gt.Storage.shell.Run("git", "config", CodeHostingDriver, value)
	return err
}

// SetCodeHostingOriginHostname sets the "github.code-hosting-driver" setting.
func (gt *GitTown) SetCodeHostingOriginHostname(value string) error {
	gt.Storage.localConfigCache[CodeHostingOriginHostname] = value
	_, err := gt.Storage.shell.Run("git", "config", CodeHostingOriginHostname, value)
	return err
}

// SetColorUI configures whether Git output contains color codes.
func (gt *GitTown) SetColorUI(value string) error {
	_, err := gt.Storage.shell.Run("git", "config", "color.ui", value)
	return err
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (gt *GitTown) SetMainBranch(branch string) error {
	_, err := gt.Storage.SetLocalConfigValue(MainBranch, branch)
	return err
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (gt *GitTown) SetNewBranchPush(value bool, global bool) error {
	setting := strconv.FormatBool(value)
	if global {
		_, err := gt.Storage.SetGlobalConfigValue(PushNewBranches, setting)
		return err
	}
	_, err := gt.Storage.SetLocalConfigValue(PushNewBranches, setting)
	return err
}

// SetOffline updates whether Git Town is in offline mode.
func (gt *GitTown) SetOffline(value bool) error {
	_, err := gt.Storage.SetGlobalConfigValue(Offline, strconv.FormatBool(value))
	return err
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (gt *GitTown) SetParent(branch, parentBranch string) error {
	_, err := gt.Storage.SetLocalConfigValue("git-town-branch."+branch+".parent", parentBranch)
	return err
}

// SetPerennialBranches marks the given branches as perennial branches.
func (gt *GitTown) SetPerennialBranches(branch []string) error {
	_, err := gt.Storage.SetLocalConfigValue(PerennialBranches, strings.Join(branch, " "))
	return err
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (gt *GitTown) SetPullBranchStrategy(strategy string) error {
	_, err := gt.Storage.SetLocalConfigValue(PullBranchStrategy, strategy)
	return err
}

// SetPushHookLocally updates the configured pull branch strategy.
func (gt *GitTown) SetPushHookLocally(value bool) error {
	_, err := gt.Storage.SetLocalConfigValue(PushHook, strconv.FormatBool(value))
	return err
}

// SetPushHook updates the configured pull branch strategy.
func (gt *GitTown) SetPushHookGlobally(value bool) error {
	_, err := gt.Storage.SetGlobalConfigValue(PushHook, strconv.FormatBool(value))
	return err
}

// SetShouldShipDeleteRemoteBranch updates the configured pull branch strategy.
func (gt *GitTown) SetShouldShipDeleteRemoteBranch(value bool) error {
	_, err := gt.Storage.SetLocalConfigValue(ShipDeleteRemoteBranch, strconv.FormatBool(value))
	return err
}

// SetShouldSyncUpstream updates the configured pull branch strategy.
func (gt *GitTown) SetShouldSyncUpstream(value bool) error {
	_, err := gt.Storage.SetLocalConfigValue(SyncUpstream, strconv.FormatBool(value))
	return err
}

func (gt *GitTown) SetSyncStrategy(value string) error {
	_, err := gt.Storage.SetLocalConfigValue(SyncStrategy, value)
	return err
}

func (gt *GitTown) SetSyncStrategyGlobal(value string) error {
	_, err := gt.Storage.SetGlobalConfigValue(SyncStrategy, value)
	return err
}

// SetTestOrigin sets the origin to be used for testing.
func (gt *GitTown) SetTestOrigin(value string) error {
	_, err := gt.Storage.SetLocalConfigValue(TestingRemoteURL, value)
	return err
}

// ShouldNewBranchPush indicates whether the current repository is configured to push
// freshly created branches up to origin.
func (gt *GitTown) ShouldNewBranchPush() (bool, error) {
	oldLocalConfig := gt.Storage.LocalConfigValue(NewBranchPushFlag)
	if oldLocalConfig != "" {
		fmt.Printf("I found the deprecated local setting %q.\n", NewBranchPushFlag)
		fmt.Printf("I am upgrading this setting to the new format %q.\n", PushNewBranches)
		err := gt.Storage.RemoveLocalConfigValue(NewBranchPushFlag)
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
	oldGlobalConfig := gt.Storage.GlobalConfigValue(NewBranchPushFlag)
	if oldGlobalConfig != "" {
		fmt.Printf("I found the deprecated global setting %q.\n", NewBranchPushFlag)
		fmt.Printf("I am upgrading this setting to the new format %q.\n", PushNewBranches)
		_, err := gt.Storage.RemoveGlobalConfigValue("git-town.new-branch-push-flag")
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
	config := gt.Storage.LocalOrGlobalConfigValue(PushNewBranches)
	if config == "" {
		return false, nil
	}
	value, err := cli.ParseBool(config)
	if err != nil {
		return false, fmt.Errorf("invalid value for %s: %q. Please provide either \"yes\" or \"no\"", PushNewBranches, config)
	}
	return value, nil
}

// ShouldNewBranchPushGlobal indictes whether the global configuration requires to push
// freshly created branches to origin.
func (gt *GitTown) ShouldNewBranchPushGlobal() (bool, error) {
	config := gt.Storage.GlobalConfigValue(PushNewBranches)
	if config == "" {
		return false, nil
	}
	return cli.ParseBool(config)
}

// ShouldShipDeleteOriginBranch indicates whether to delete the remote branch after shipping.
func (gt *GitTown) ShouldShipDeleteOriginBranch() (bool, error) {
	setting := gt.Storage.LocalOrGlobalConfigValue(ShipDeleteRemoteBranch)
	if setting == "" {
		return true, nil
	}
	result, err := strconv.ParseBool(setting)
	if err != nil {
		return true, fmt.Errorf("invalid value for %s: %q. Please provide either \"true\" or \"false\"", ShipDeleteRemoteBranch, setting)
	}
	return result, nil
}

// ShouldSyncUpstream indicates whether this repo should sync with its upstream.
func (gt *GitTown) ShouldSyncUpstream() (bool, error) {
	text := gt.Storage.LocalOrGlobalConfigValue(SyncUpstream)
	if text == "" {
		return true, nil
	}
	return cli.ParseBool(text)
}

func (gt *GitTown) SyncStrategy() string {
	setting := gt.Storage.LocalOrGlobalConfigValue(SyncStrategy)
	if setting == "" {
		setting = "merge"
	}
	return setting
}

func (gt *GitTown) SyncStrategyGlobal() string {
	setting := gt.Storage.GlobalConfigValue(SyncStrategy)
	if setting == "" {
		setting = "merge"
	}
	return setting
}

// ValidateIsOnline asserts that Git Town is not in offline mode.
func (gt *GitTown) ValidateIsOnline() error {
	isOffline, err := gt.IsOffline()
	if err != nil {
		return err
	}
	if isOffline {
		return errors.New("this command requires an active internet connection")
	}
	return nil
}
