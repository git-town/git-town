package config

import (
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/config/gitconfig"
	"github.com/git-town/git-town/v18/internal/git"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

type UnvalidatedConfig struct {
	NormalConfig      NormalConfig
	UnvalidatedConfig configdomain.UnvalidatedConfigData
}

func (self *UnvalidatedConfig) BranchType(branch gitdomain.LocalBranchName) configdomain.BranchType {
	return self.UnvalidatedConfig.PartialBranchType(branch).GetOrElse(self.NormalConfig.PartialBranchType(branch))
}

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *UnvalidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	branchType := self.BranchType(branch)
	return branchType == configdomain.BranchTypeMainBranch || branchType == configdomain.BranchTypePerennialBranch
}

func (self *UnvalidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	if mainBranch, hasMainBranch := self.UnvalidatedConfig.MainBranch.Get(); hasMainBranch {
		return append(gitdomain.LocalBranchNames{mainBranch}, self.NormalConfig.PerennialBranches...)
	}
	return self.NormalConfig.PerennialBranches
}

func (self *UnvalidatedConfig) Reload() (globalSnapshot, localSnapshot configdomain.SingleSnapshot) {
	globalSnapshot, globalGitConfig, _ := self.NormalConfig.GitConfigAccess.LoadGlobal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	localSnapshot, localGitConfig, _ := self.NormalConfig.GitConfigAccess.LoadLocal(false)    // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	unvalidatedConfig, normalConfig := NewConfigs(self.NormalConfig.ConfigFile, globalGitConfig, localGitConfig)
	self.UnvalidatedConfig = unvalidatedConfig
	self.NormalConfig = NormalConfig{
		ConfigFile:       self.NormalConfig.ConfigFile,
		DryRun:           self.NormalConfig.DryRun,
		GitConfigAccess:  self.NormalConfig.GitConfigAccess,
		GitVersion:       self.NormalConfig.GitVersion,
		GlobalGitConfig:  globalGitConfig,
		LocalGitConfig:   localGitConfig,
		NormalConfigData: normalConfig,
	}
	return globalSnapshot, localSnapshot
}

func (self *UnvalidatedConfig) RemoveMainBranch() {
	_ = self.NormalConfig.GitConfigAccess.RemoveLocalConfigValue(configdomain.KeyMainBranch)
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.UnvalidatedConfig.MainBranch = Some(branch)
	return self.NormalConfig.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyMainBranch, branch.String())
}

// UnvalidatedBranchesAndTypes provides the types for the given branches.
// This method's name startes with "Unvalidated" to indicate that the types might be incomplete,
// and you should use ValidatedConfig.BranchesAndTypes if possible.
func (self *UnvalidatedConfig) UnvalidatedBranchesAndTypes(branches gitdomain.LocalBranchNames) configdomain.BranchesAndTypes {
	result := make(configdomain.BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.UnvalidatedConfig.PartialBranchType(branch).GetOrElse(self.NormalConfig.PartialBranchType(branch))
	}
	return result
}

func DefaultUnvalidatedConfig(gitAccess gitconfig.Access, gitVersion git.Version) UnvalidatedConfig {
	return UnvalidatedConfig{
		NormalConfig: NormalConfig{
			ConfigFile:       None[configdomain.PartialConfig](),
			DryRun:           false,
			GitConfigAccess:  gitAccess,
			GitVersion:       gitVersion,
			GlobalGitConfig:  configdomain.EmptyPartialConfig(),
			LocalGitConfig:   configdomain.EmptyPartialConfig(),
			NormalConfigData: configdomain.DefaultNormalConfig(),
		},
		UnvalidatedConfig: configdomain.DefaultUnvalidatedConfig(),
	}
}

func MergeConfigs(configFile Option[configdomain.PartialConfig], globalGitConfig, localGitConfig configdomain.PartialConfig) (configdomain.UnvalidatedConfigData, configdomain.NormalConfigData) {
	result := configdomain.EmptyPartialConfig()
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		result = result.Merge(configFile)
	}
	result = result.Merge(globalGitConfig)
	result = result.Merge(localGitConfig)
	return result.ToUnvalidatedConfig(), result.ToNormalConfig(configdomain.DefaultNormalConfig())
}

func NewConfigs(configFile Option[configdomain.PartialConfig], globalGitConfig, localGitConfig configdomain.PartialConfig) (configdomain.UnvalidatedConfigData, configdomain.NormalConfigData) {
	config := configdomain.EmptyPartialConfig()
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		config = config.Merge(configFile)
	}
	config = config.Merge(globalGitConfig)
	config = config.Merge(localGitConfig)
	normalConfig := config.ToNormalConfig(configdomain.DefaultNormalConfig())
	unvalidatedConfig := config.ToUnvalidatedConfig()
	return unvalidatedConfig, normalConfig
}

func NewUnvalidatedConfig(args NewUnvalidatedConfigArgs) UnvalidatedConfig {
	unvalidatedConfig, normalConfig := MergeConfigs(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	return UnvalidatedConfig{
		NormalConfig: NormalConfig{
			ConfigFile:       args.ConfigFile,
			DryRun:           args.DryRun,
			GitConfigAccess:  args.Access,
			GitVersion:       args.GitVersion,
			GlobalGitConfig:  args.GlobalConfig,
			LocalGitConfig:   args.LocalConfig,
			NormalConfigData: normalConfig,
		},
		UnvalidatedConfig: unvalidatedConfig,
	}
}

type NewUnvalidatedConfigArgs struct {
	Access        gitconfig.Access
	ConfigFile    Option[configdomain.PartialConfig]
	DryRun        configdomain.DryRun
	FinalMessages stringslice.Collector
	GitVersion    git.Version
	GlobalConfig  configdomain.PartialConfig
	LocalConfig   configdomain.PartialConfig
}
