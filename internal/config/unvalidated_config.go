package config

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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

func (self *UnvalidatedConfig) Reload() (globalSnapshot, localSnapshot, unscopedSnapshot configdomain.SingleSnapshot) {
	globalSnapshot, _ = self.NormalConfig.GitIO.LoadSnapshot(Some(configdomain.ConfigScopeGlobal), configdomain.UpdateOutdatedNo) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	localSnapshot, _ = self.NormalConfig.GitIO.LoadSnapshot(Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedNo)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	unscopedSnapshot, _ = self.NormalConfig.GitIO.LoadSnapshot(None[configdomain.ConfigScope](), configdomain.UpdateOutdatedNo)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	unscopedGitConfig, _ := configdomain.NewPartialConfigFromSnapshot(unscopedSnapshot, false, nil)
	envConfig := envconfig.Load()
	unvalidatedConfig, normalConfig := mergeConfigs(mergeConfigsArgs{
		env:  envConfig,
		file: self.NormalConfig.ConfigFile,
		git:  unscopedGitConfig,
	})
	self.UnvalidatedConfig = unvalidatedConfig
	self.NormalConfig = NormalConfig{
		ConfigFile:       self.NormalConfig.ConfigFile,
		DryRun:           self.NormalConfig.DryRun,
		EnvConfig:        envConfig,
		GitConfig:        unscopedGitConfig,
		GitIO:            self.NormalConfig.GitIO,
		GitVersion:       self.NormalConfig.GitVersion,
		NormalConfigData: normalConfig,
	}
	return globalSnapshot, localSnapshot, unscopedSnapshot
}

func (self *UnvalidatedConfig) RemoveMainBranch() {
	_ = self.NormalConfig.GitIO.RemoveConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyMainBranch)
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.UnvalidatedConfig.MainBranch = Some(branch)
	return self.NormalConfig.GitIO.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyMainBranch, branch.String())
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

func DefaultUnvalidatedConfig(gitVersion git.Version) UnvalidatedConfig {
	return UnvalidatedConfig{
		NormalConfig: NormalConfig{
			ConfigFile:       None[configdomain.PartialConfig](),
			DryRun:           false,
			EnvConfig:        configdomain.EmptyPartialConfig(),
			GitConfig:        configdomain.EmptyPartialConfig(),
			GitVersion:       gitVersion,
			NormalConfigData: configdomain.DefaultNormalConfig(),
		},
		UnvalidatedConfig: configdomain.DefaultUnvalidatedConfig(),
	}
}

func NewUnvalidatedConfig(args NewUnvalidatedConfigArgs) UnvalidatedConfig {
	unvalidatedConfig, normalConfig := mergeConfigs(mergeConfigsArgs{
		env:  args.EnvConfig,
		file: args.ConfigFile,
		git:  args.GitConfig,
	})
	return UnvalidatedConfig{
		NormalConfig: NormalConfig{
			ConfigFile:       args.ConfigFile,
			DryRun:           args.DryRun,
			EnvConfig:        args.EnvConfig,
			GitConfig:        args.GitConfig,
			GitVersion:       args.GitVersion,
			NormalConfigData: normalConfig,
		},
		UnvalidatedConfig: unvalidatedConfig,
	}
}

type NewUnvalidatedConfigArgs struct {
	ConfigFile    Option[configdomain.PartialConfig]
	DryRun        configdomain.DryRun
	EnvConfig     configdomain.PartialConfig
	FinalMessages stringslice.Collector
	GitConfig     configdomain.PartialConfig
	GitVersion    git.Version
}

func mergeConfigs(args mergeConfigsArgs) (configdomain.UnvalidatedConfigData, configdomain.NormalConfigData) {
	result := configdomain.EmptyPartialConfig()
	if configFile, hasConfigFile := args.file.Get(); hasConfigFile {
		result = result.Merge(configFile)
	}
	result = result.Merge(args.git)
	result = result.Merge(args.env)
	return result.ToUnvalidatedConfig(), result.ToNormalConfig(configdomain.DefaultNormalConfig())
}

type mergeConfigsArgs struct {
	env  configdomain.PartialConfig         // configuration data taken from environment variables
	file Option[configdomain.PartialConfig] // data of the configuration file
	git  configdomain.PartialConfig         // data from the unscoped Git configuration
}
