/*
This file contains functionality around storing configuration settings
inside Git's metadata storage for the repository.
*/

package git

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
)

// Config provides access to the Git Town configuration.
var Config *Configuration

// Configuration manages the Git Town configuration,
// stored in Git metadata in the local repo and in the global Git configuration.
type Configuration struct {

	// localConfig is a cache of the Git configuration in the local directory.
	localConfig map[string]string

	// globalConfig is a cache of the global Git configuration
	globalConfig map[string]string
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewConfiguration(dir string) *Configuration {
	result := &Configuration{
		localConfig:  map[string]string{},
		globalConfig: map[string]string{},
	}
	result.initializeCache(false, result.localConfig)
	result.initializeCache(true, result.globalConfig)
	return result
}

// AddAlias adds an alias for the given Git Town command.
func (c *Configuration) AddAlias(cmd string) {
	key := "alias." + cmd
	value := "town " + cmd
	c.globalConfig[key] = value
	command.New("git", "config", "--global", key, value).Output()
}

// AddToPerennialBranches adds the given branch as a perennial branch
func (c *Configuration) AddToPerennialBranches(branchName string) {
	c.SetPerennialBranches(append(c.GetPerennialBranches(), branchName))
}

// DeleteParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func (c *Configuration) DeleteParentBranch(branchName string) {
	c.removeLocalConfigurationValue("git-town-branch." + branchName + ".parent")
}

// EnsureIsFeatureBranch asserts that the given branch is a feature branch.
func (c *Configuration) EnsureIsFeatureBranch(branchName, errorSuffix string) {
	util.Ensure(c.IsFeatureBranch(branchName), fmt.Sprintf("The branch '%s' is not a feature branch. %s", branchName, errorSuffix))
}

// GetAncestorBranches returns the names of all parent branches for the given branch,
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func (c *Configuration) GetAncestorBranches(branchName string) (result []string) {
	parentBranchMap := c.GetParentBranchMap()
	current := branchName
	for {
		if c.IsMainBranch(current) || c.IsPerennialBranch(current) {
			return
		}
		parent := parentBranchMap[current]
		if parent == "" {
			return
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// GetCodeHostingDriver provides the type of driver to use to talk to the code hosting service.
func (c *Configuration) GetCodeHostingDriver() string {
	return c.localConfig["git-town.code-hosting-driver"]
}

// GetCodeHostingOriginHostname provides the hostname of the code hosting server to use.
func (c *Configuration) GetCodeHostingOriginHostname() string {
	return c.localConfig["git-town.code-hosting-origin-hostname"]
}

// GetGithubAPIToken provides the API token to talk to the GitHub API.
func (c *Configuration) GetGithubAPIToken() string {
	return c.localConfig["git-town.github-token"]
}

// GetParentBranchMap returns a map from branch name to its parent branch
func (c *Configuration) GetParentBranchMap() map[string]string {
	result := map[string]string{}
	for _, key := range c.getConfigurationKeysMatching("^git-town-branch\\..*\\.parent$") {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := c.localConfig[key]
		result[child] = parent
	}
	return result
}

// GetChildBranches returns the names of all branches for which the given branch
// is a parent.
func (c *Configuration) GetChildBranches(branchName string) (result []string) {
	for _, key := range c.getConfigurationKeysMatching("^git-town-branch\\..*\\.parent$") {
		parent := c.localConfig[key]
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	return
}

// GetMainBranch returns the name of the main branch.
func (c *Configuration) GetMainBranch() string {
	return c.localConfig["git-town.main-branch-name"]
}

// GetParentBranch returns the name of the parent branch of the given branch.
func (c *Configuration) GetParentBranch(branchName string) string {
	return c.localConfig["git-town-branch."+branchName+".parent"]
}

// GetPerennialBranches returns all branches that are marked as perennial.
func (c *Configuration) GetPerennialBranches() []string {
	result := c.localConfig["git-town.perennial-branch-names"]
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// GetPullBranchStrategy returns the currently configured pull branch strategy.
func (c *Configuration) GetPullBranchStrategy() string {
	return c.getConfigurationValueWithDefault("git-town.pull-branch-strategy", "rebase")
}

// GetRemoteOriginURL returns the URL for the "origin" remote.
// In tests this value can be stubbed.
func (c *Configuration) GetRemoteOriginURL() string {
	if os.Getenv("GIT_TOWN_ENV") == "test" {
		mockRemoteURL := c.localConfig["git-town.testing.remote-url"]
		if mockRemoteURL != "" {
			return mockRemoteURL
		}
	}
	return command.New("git", "remote", "get-url", "origin").Output()
}

// GetRemoteUpstreamURL returns the URL of the "upstream" remote.
func (c *Configuration) GetRemoteUpstreamURL() string {
	return command.New("git", "remote", "get-url", "upstream").Output()
}

// GetSyncUpstream indicates whether this repository is configured to sync to its upstream remote.
func (c *Configuration) GetSyncUpstream() bool {
	return c.localConfig["git-town.sync-upstream"] != "false"
}

// GetURLHostname returns the hostname contained within the given Git URL.
func (c *Configuration) GetURLHostname(url string) string {
	hostnameRegex := regexp.MustCompile("(^[^:]*://([^@]*@)?|git@)([^/:]+).*")
	matches := hostnameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return matches[3]
}

// GetURLRepositoryName returns the repository name contains within the given Git URL.
func (c *Configuration) GetURLRepositoryName(url string) string {
	hostname := c.GetURLHostname(url)
	repositoryNameRegex, err := regexp.Compile(".*" + hostname + "[/:](.+)")
	exit.IfWrap(err, "Error compiling repository name regular expression")
	matches := repositoryNameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return strings.TrimSuffix(matches[1], ".git")
}

// HasGlobalConfigurationValue returns whether there is a global git configuration for the given key
func (c *Configuration) HasGlobalConfigurationValue(key string) bool {
	return command.New("git", "config", "-l", "--global", "--name").OutputContainsLine(key)
}

// HasParentBranch returns whether or not the given branch has a parent
func (c *Configuration) HasParentBranch(branchName string) bool {
	return c.GetParentBranch(branchName) != ""
}

// IsAncestorBranch returns whether the given branch is an ancestor of the other given branch.
func (c *Configuration) IsAncestorBranch(branchName, ancestorBranchName string) bool {
	ancestorBranches := c.GetAncestorBranches(branchName)
	return util.DoesStringArrayContain(ancestorBranches, ancestorBranchName)
}

// HasRemote returns whether the current repository contains a Git remote
// with the given name.
func (c *Configuration) HasRemote(name string) bool {
	return util.DoesStringArrayContain(c.getRemotes(), name)
}

// IsFeatureBranch returns whether the branch with the given name is
// a feature branch.
func (c *Configuration) IsFeatureBranch(branchName string) bool {
	return !c.IsMainBranch(branchName) && !c.IsPerennialBranch(branchName)
}

// IsMainBranch returns whether the branch with the given name
// is the main branch of the repository.
func (c *Configuration) IsMainBranch(branchName string) bool {
	return branchName == c.GetMainBranch()
}

// IsPerennialBranch returns whether the branch with the given name is
// a perennial branch.
func (c *Configuration) IsPerennialBranch(branchName string) bool {
	perennialBranches := c.GetPerennialBranches()
	return util.DoesStringArrayContain(perennialBranches, branchName)
}

// RemoveAlias removes the global alias for the given Git Town command.
func (c *Configuration) RemoveAlias(cmd string) {
	key := "alias." + cmd
	previousAlias := c.globalConfig[key]
	if previousAlias == "town "+cmd {
		command.New("git", "config", "--global", "--unset", key).Output()
	}
}

// RemoveAllConfiguration removes all Git Town configuration
func (c *Configuration) RemoveAllConfiguration() {
	command.New("git", "config", "--remove-section", "git-town").Output()
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration
func (c *Configuration) RemoveOutdatedConfiguration() {
	for child, parent := range c.GetParentBranchMap() {
		if !HasBranch(child) || !HasBranch(parent) {
			c.DeleteParentBranch(child)
		}
	}
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch
func (c *Configuration) RemoveFromPerennialBranches(branchName string) {
	c.SetPerennialBranches(util.RemoveStringFromSlice(c.GetPerennialBranches(), branchName))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (c *Configuration) SetMainBranch(branchName string) {
	c.setConfigurationValue("git-town.main-branch-name", branchName)
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (c *Configuration) SetParentBranch(branchName, parentBranchName string) {
	c.setConfigurationValue("git-town-branch."+branchName+".parent", parentBranchName)
}

// SetPerennialBranches marks the given branches as perennial branches
func (c *Configuration) SetPerennialBranches(branchNames []string) {
	c.setConfigurationValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (c *Configuration) SetPullBranchStrategy(strategy string) {
	c.setConfigurationValue("git-town.pull-branch-strategy", strategy)
}

// ShouldNewBranchPush returns whether the current repository is configured to push
// freshly created branches up to the origin remote.
func (c *Configuration) ShouldNewBranchPush() bool {
	return util.StringToBool(c.getConfigurationValueWithDefault("git-town.new-branch-push-flag", "false"))
}

// GetGlobalNewBranchPushFlag returns the global configuration for to push
// freshly created branches up to the origin remote.
func (c *Configuration) GetGlobalNewBranchPushFlag() string {
	return c.getGlobalConfigurationValueWithDefault("git-town.new-branch-push-flag", "false")
}

// UpdateOffline updates whether Git Town is in offline mode
func (c *Configuration) UpdateOffline(value bool) {
	c.setGlobalConfigurationValue("git-town.offline", strconv.FormatBool(value))
}

// UpdateShouldNewBranchPush updates whether the current repository is configured to push
// freshly created branches up to the origin remote.
func (c *Configuration) UpdateShouldNewBranchPush(value bool) {
	c.setConfigurationValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
}

// UpdateGlobalShouldNewBranchPush updates global whether to push
// freshly created branches up to the origin remote.
func (c *Configuration) UpdateGlobalShouldNewBranchPush(value bool) {
	c.setGlobalConfigurationValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
}

// Helpers

func (c *Configuration) getGlobalConfigurationValueWithDefault(key, defaultValue string) string {
	value := c.globalConfig[key]
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Configuration) getConfigurationValueWithDefault(key, defaultValue string) string {
	value := c.localConfig[key]
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Configuration) getConfigurationKeysMatching(toMatch string) (result []string) {
	re := regexp.MustCompile(toMatch)
	for key := range c.localConfig {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return result
}

func (c *Configuration) initializeCache(global bool, cache map[string]string) {
	cmdArgs := []string{"git", "config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	}
	cmd := command.New(cmdArgs...)
	if cmd.Err() != nil && strings.Contains(cmd.Output(), "No such file or directory") {
		return
	}
	exit.If(cmd.Err())
	if cmd.Output() == "" {
		return
	}
	for _, line := range strings.Split(cmd.Output(), "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		cache[key] = value
	}
}

func (c *Configuration) setConfigurationValue(key, value string) {
	command.New("git", "config", key, value).Run()
	c.localConfig[key] = value
}

func (c *Configuration) setGlobalConfigurationValue(key, value string) {
	command.New("git", "config", "--global", key, value).Run()
	c.globalConfig[key] = value
	c.localConfig = map[string]string{} // Need to reset config in case it was inheriting
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (c *Configuration) removeLocalConfigurationValue(key string) {
	command.New("git", "config", "--unset", key).Run()
	delete(c.localConfig, key)
}

// Remotes are cached in order to minimize the number of git commands run
var remotes []string
var remotesInitialized bool

func (c *Configuration) getRemotes() []string {
	if !remotesInitialized {
		remotes = command.New("git", "remote").OutputLines()
		remotesInitialized = true
	}
	return remotes
}

// Init

func init() {
	Config = NewConfiguration(".")
}
