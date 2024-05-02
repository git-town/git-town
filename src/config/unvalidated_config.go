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
	ConfigFile      Option[configdomain.PartialConfig] // content of git-town.toml, nil = no config file exists
	DryRun          bool
	Config          configdomain.UnvalidatedConfig // the merged configuration data
	GitConfig       gitconfig.Access               // access to the Git configuration settings
	GlobalGitConfig configdomain.PartialConfig     // content of the global Git configuration
	LocalGitConfig  configdomain.PartialConfig     // content of the local Git configuration
	originURLCache  configdomain.OriginURLCache    // TODO: remove if unused
}

func NewUnvalidatedConfig(args NewConfigArgs) (UnvalidatedConfig, *stringslice.Collector, error) {
	config := configdomain.NewUnvalidatedConfig(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	configAccess := gitconfig.Access{Runner: args.Runner}
	finalMessages := stringslice.Collector{}
	err := cleanupPerennialParentEntries(config.Lineage, config.MainAndPerennials(), configAccess, &finalMessages)
	return UnvalidatedConfig{
		Config:          config,
		ConfigFile:      args.ConfigFile,
		DryRun:          args.DryRun,
		GitConfig:       configAccess,
		GlobalGitConfig: args.GlobalConfig,
		LocalGitConfig:  args.LocalConfig,
		originURLCache:  configdomain.OriginURLCache{},
	}, &finalMessages, err
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

// SetPerennialBranches marks the given branches as perennial branches.
func (self *UnvalidatedConfig) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.Config.PerennialBranches = branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialBranches, branches.Join(" "))
}
