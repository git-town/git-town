package config

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/confighelpers"
	"github.com/git-town/git-town/v14/src/config/envconfig"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
)

type UnvalidatedConfig struct {
	Config          configdomain.UnvalidatedConfig     // the merged configuration data
	ConfigFile      Option[configdomain.PartialConfig] // content of git-town.toml, nil = no config file exists
	DryRun          bool
	GitConfig       gitconfig.Access            // access to the Git configuration settings
	GlobalGitConfig configdomain.PartialConfig  // content of the global Git configuration
	LocalGitConfig  configdomain.PartialConfig  // content of the local Git configuration
	originURLCache  configdomain.OriginURLCache // TODO: remove if unused
}

func NewUnvalidatedConfig(args NewUnvalidatedConfigArgs) (UnvalidatedConfig, *stringslice.Collector) {
	config := configdomain.NewUnvalidatedConfig(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	configAccess := gitconfig.Access{Runner: args.Runner}
	finalMessages := stringslice.Collector{}
	return UnvalidatedConfig{
		Config:          config,
		ConfigFile:      args.ConfigFile,
		DryRun:          args.DryRun,
		GitConfig:       configAccess,
		GlobalGitConfig: args.GlobalConfig,
		LocalGitConfig:  args.LocalConfig,
		originURLCache:  configdomain.OriginURLCache{},
	}, &finalMessages
}

// OriginURL provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *UnvalidatedConfig) OriginURL() Option[giturl.Parts] {
	text := self.OriginURLString()
	if text == "" {
		return None[giturl.Parts]()
	}
	return confighelpers.DetermineOriginURL(text, self.Config.HostingOriginHostname, self.originURLCache)
}

// OriginURLString provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *UnvalidatedConfig) OriginURLString() string {
	remoteOverride := envconfig.OriginURLOverride()
	if remoteOverride != "" {
		return remoteOverride
	}
	return self.GitConfig.OriginRemote()
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration.
func (self *UnvalidatedConfig) RemoveOutdatedConfiguration(localBranches gitdomain.LocalBranchNames) error {
	for child, parent := range self.Config.Lineage {
		hasChildBranch := localBranches.Contains(child)
		hasParentBranch := localBranches.Contains(parent)
		if !hasChildBranch || !hasParentBranch {
			self.RemoveParent(child)
		}
	}
	return nil
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *UnvalidatedConfig) RemoveParent(branch gitdomain.LocalBranchName) {
	if self.LocalGitConfig.Lineage != nil {
		self.LocalGitConfig.Lineage.RemoveBranch(branch)
	}
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.NewParentKey(branch))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.Config.MainBranch = Some(branch)
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyMainBranch, branch.String())
}

// SetOffline updates whether Git Town is in offline mode.
func (self *UnvalidatedConfig) SetOffline(value configdomain.Offline) error {
	self.Config.Offline = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyOffline, value.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Config.Lineage[branch] = parentBranch
	return self.GitConfig.SetLocalConfigValue(gitconfig.NewParentKey(branch), parentBranch.String())
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *UnvalidatedConfig) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.PerennialBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialBranches, branches.Join(" "))
}

type NewUnvalidatedConfigArgs struct {
	ConfigFile   Option[configdomain.PartialConfig]
	DryRun       bool
	GlobalConfig configdomain.PartialConfig
	LocalConfig  configdomain.PartialConfig
	Runner       gitconfig.Runner
}
