// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/git-town/git-town/v7/src/run"
)

type Config struct {
	gitConfig         *gitConfig
	PerennialBranches *PerennialBranches
	Ancestry          *Ancestry
	Hosting           *Hosting
	Offline           *Offline
}

func NewConfiguration(shell run.Shell) Config {
	gitConfig := gitConfig{
		localConfigCache:  loadGitConfig(shell, false),
		globalConfigCache: loadGitConfig(shell, true),
		shell:             shell,
	}
	config := Config{
		gitConfig: &gitConfig,
		PerennialBranches: &PerennialBranches{
			gc: &gitConfig,
		},
		Offline: &Offline{
			gitConfig: &gitConfig,
		},
		Hosting: &Hosting{
			gitConfig: &gitConfig,
		},
	}
	config.Ancestry = &Ancestry{
		gitConfig: &gitConfig,
		config:    &config,
	}
	return config
}

// Reload refreshes the cached configuration information.
func (c *Config) Reload() {
	c.gitConfig.Reload()
}

func (c *Config) RemoveLocalGitConfiguration() error {
	return c.gitConfig.RemoveLocalGitConfiguration()
}

func (c *Config) SetLocalConfigValueInTests(key, value string) (*run.Result, error) {
	if flag.Lookup("test.v") == nil {
		panic("this function is only for tests")
	}
	return c.gitConfig.SetLocalConfigValue(key, value)
}

func (c *Config) SetGlobalConfigValueInTests(key, value string) (*run.Result, error) {
	if flag.Lookup("test.v") == nil {
		panic("this function is only for tests")
	}
	return c.gitConfig.SetGlobalConfigValue(key, value)
}

// AddGitAlias sets the given Git alias.
func (c *Config) AddGitAlias(command string) (*run.Result, error) {
	return c.gitConfig.SetGlobalConfigValue("alias."+command, "town "+command)
}

// DeleteMainBranchConfiguration removes the configuration entry for the main branch name.
func (c *Config) DeleteMainBranchConfiguration() error {
	return c.gitConfig.removeLocalConfigValue("git-town.main-branch-name")
}

// GitAlias provides the currently set alias for the given Git Town command.
func (c *Config) GitAlias(command string) string {
	return c.gitConfig.globalConfigValue("alias." + command)
}

// IsFeatureBranch indicates whether the branch with the given name is
// a feature branch.
func (c *Config) IsFeatureBranch(branchName string) bool {
	return !c.IsMainBranch(branchName) && !c.PerennialBranches.Is(branchName)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (c *Config) IsMainBranch(branchName string) bool {
	return branchName == c.MainBranch()
}

// localConfigKeysMatching provides the names of the Git Town configuration keys matching the given RegExp string.
func (c *Config) localConfigKeysMatching(toMatch string) []string {
	result := []string{}
	re := regexp.MustCompile(toMatch)
	for key := range c.gitConfig.localConfigCache {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return result
}

// MainBranch provides the name of the main branch.
func (c *Config) MainBranch() string {
	return c.gitConfig.localOrGlobalConfigValue("git-town.main-branch-name")
}

// PullBranchStrategy provides the currently configured pull branch strategy.
func (c *Config) PullBranchStrategy() string {
	config := c.gitConfig.localOrGlobalConfigValue("git-town.pull-branch-strategy")
	if config != "" {
		return config
	}
	return "rebase"
}

// PushVerify provides the currently configured pull branch strategy.
func (c *Config) PushVerify() bool {
	config := c.gitConfig.localOrGlobalConfigValue("git-town.push-verify")
	if config == "" {
		return true
	}
	result, err := strconv.ParseBool(config)
	if err != nil {
		fmt.Printf("Invalid value for git-town.push-verify: %q. Please provide either true or false. Considering true for now.", config)
		fmt.Println()
		return true
	}
	return result
}

// OriginURL provides the URL for the "origin" remote.
// In tests this value can be stubbed.
func (c *Config) OriginURL() string {
	remote := os.Getenv("GIT_TOWN_REMOTE")
	if remote != "" {
		return remote
	}
	res, _ := c.gitConfig.shell.Run("git", "remote", "get-url", "origin")
	return res.OutputSanitized()
}

// RemoveGitAlias removes the given Git alias.
func (c *Config) RemoveGitAlias(command string) (*run.Result, error) {
	return c.gitConfig.removeGlobalConfigValue("alias." + command)
}

// SetColorUI configures whether Git output contains color codes.
func (c *Config) SetColorUI(value string) error {
	_, err := c.gitConfig.shell.Run("git", "config", "color.ui", value)
	return err
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (c *Config) SetMainBranch(branchName string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.main-branch-name", branchName)
	return err
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (c *Config) SetNewBranchPush(value bool, global bool) error {
	if global {
		_, err := c.gitConfig.SetGlobalConfigValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
		return err
	}
	_, err := c.gitConfig.SetLocalConfigValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
	return err
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (c *Config) SetPullBranchStrategy(strategy string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.pull-branch-strategy", strategy)
	return err
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (c *Config) SetPushVerify(strategy string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.push-verify", strategy)
	return err
}

// SetShouldShipDeleteRemoteBranch updates the configured pull branch strategy.
func (c *Config) SetShouldShipDeleteRemoteBranch(value bool) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.ship-delete-remote-branch", strconv.FormatBool(value))
	return err
}

// SetShouldSyncUpstream updates the configured pull branch strategy.
func (c *Config) SetShouldSyncUpstream(value bool) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.sync-upstream", strconv.FormatBool(value))
	return err
}

func (c *Config) SetSyncStrategy(value string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.sync-strategy", value)
	return err
}

// ShouldNewBranchPush indicates whether the current repository is configured to push
// freshly created branches up to origin.
func (c *Config) ShouldNewBranchPush() bool {
	config := c.gitConfig.localOrGlobalConfigValue("git-town.new-branch-push-flag")
	if config == "" {
		return false
	}
	value, err := strconv.ParseBool(config)
	if err != nil {
		fmt.Printf("Invalid value for git-town.new-branch-push-flag: %q. Please provide either true or false. Considering false for now.\n", config)
		return false
	}
	return value
}

// ShouldNewBranchPushGlobal indictes whether the global configuration requires to push
// freshly created branches to origin.
func (c *Config) ShouldNewBranchPushGlobal() bool {
	config := c.gitConfig.globalConfigValue("git-town.new-branch-push-flag")
	return config == "true"
}

// ShouldShipDeleteOriginBranch indicates whether to delete the remote branch after shipping.
func (c *Config) ShouldShipDeleteOriginBranch() bool {
	setting := c.gitConfig.localOrGlobalConfigValue("git-town.ship-delete-remote-branch")
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
func (c *Config) ShouldSyncUpstream() bool {
	return c.gitConfig.localOrGlobalConfigValue("git-town.sync-upstream") != "false"
}

func (c *Config) SyncStrategy() string {
	setting := c.gitConfig.localOrGlobalConfigValue("git-town.sync-strategy")
	if setting == "" {
		setting = "merge"
	}
	return setting
}
