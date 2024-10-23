package config

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/config/gitconfig"
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// TODO: rename to UnvalidatedConfigData
type UnvalidatedConfig struct {
	NormalConfig      NormalConfig
	UnvalidatedConfig configdomain.UnvalidatedConfig
}

// TODO: delete?
type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

func (self *UnvalidatedConfig) BranchType(branch gitdomain.LocalBranchName) configdomain.BranchType {
	return self.UnvalidatedConfig.PartialBranchType(branch).GetOrElse(self.NormalConfig.PartialBranchType(branch))
}

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *UnvalidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return self.UnvalidatedConfig.IsMainBranch(branch) || self.NormalConfig.IsPerennialBranch(branch)
}

func (self *UnvalidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	if mainBranch, hasMainBranch := self.UnvalidatedConfig.MainBranch.Get(); hasMainBranch {
		return append(gitdomain.LocalBranchNames{mainBranch}, self.NormalConfig.PerennialBranches...)
	}
	return self.NormalConfig.PerennialBranches
}

func (self *UnvalidatedConfig) Reload() {
	_, globalGitConfig, _ := self.NormalConfig.GitConfig.LoadGlobal(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	_, localGitConfig, _ := self.NormalConfig.GitConfig.LoadLocal(false)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	fmt.Println("1111111111111111111111111111111111111111111111 LOCAL", localGitConfig)
	unvalidatedConfig, normalConfig := NewConfigs(self.NormalConfig.ConfigFile, self.NormalConfig.GlobalGitConfig, self.NormalConfig.LocalGitConfig)
	self.UnvalidatedConfig = unvalidatedConfig
	self.NormalConfig = NormalConfig{
		ConfigFile:      self.NormalConfig.ConfigFile,
		DryRun:          self.NormalConfig.DryRun,
		GitConfig:       self.NormalConfig.GitConfig,
		GitVersion:      self.NormalConfig.GitVersion,
		GlobalGitConfig: globalGitConfig,
		LocalGitConfig:  localGitConfig,
		NormalConfig:    normalConfig,
	}
}

func (self *UnvalidatedConfig) RemoveMainBranch() {
	_ = self.NormalConfig.GitConfig.RemoveLocalConfigValue(configdomain.KeyMainBranch)
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.UnvalidatedConfig.MainBranch = Some(branch)
	return self.NormalConfig.GitConfig.SetLocalConfigValue(configdomain.KeyMainBranch, branch.String())
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
			ConfigFile:      None[configdomain.PartialConfig](),
			DryRun:          false,
			GitConfig:       gitAccess,
			GitVersion:      gitVersion,
			GlobalGitConfig: configdomain.EmptyPartialConfig(),
			LocalGitConfig:  configdomain.EmptyPartialConfig(),
			NormalConfig:    configdomain.DefaultNormalConfig(),
		},
		UnvalidatedConfig: configdomain.DefaultUnvalidatedConfig(),
	}
}

func MergeConfigs(configFile Option[configdomain.PartialConfig], globalGitConfig, localGitConfig configdomain.PartialConfig) (configdomain.UnvalidatedConfig, configdomain.NormalConfig) {
	result := configdomain.EmptyPartialConfig()
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		result = result.Merge(configFile)
	}
	result = result.Merge(globalGitConfig)
	result = result.Merge(localGitConfig)
	return result.ToUnvalidatedConfig(), result.ToNormalConfig(configdomain.DefaultNormalConfig())
}

func NewConfigs(configFile Option[configdomain.PartialConfig], globalGitConfig, localGitConfig configdomain.PartialConfig) (configdomain.UnvalidatedConfig, configdomain.NormalConfig) {
	config := configdomain.EmptyPartialConfig()
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		config = config.Merge(configFile)
	}
	config = config.Merge(globalGitConfig)
	fmt.Println("222222222222222222222222222222222222", localGitConfig.Lineage)
	config = config.Merge(localGitConfig)
	fmt.Println("333333333333333333333333333333333333", config.Lineage)
	normalConfig := config.ToNormalConfig(configdomain.DefaultNormalConfig())
	unvalidatedConfig := config.ToUnvalidatedConfig()
	return unvalidatedConfig, normalConfig
}

func NewUnvalidatedConfig(args NewUnvalidatedConfigArgs) (UnvalidatedConfig, stringslice.Collector) {
	unvalidatedConfig, normalConfig := MergeConfigs(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	finalMessages := stringslice.NewCollector()
	return UnvalidatedConfig{
		NormalConfig: NormalConfig{
			ConfigFile:      args.ConfigFile,
			DryRun:          args.DryRun,
			GitConfig:       args.Access,
			GitVersion:      args.GitVersion,
			GlobalGitConfig: args.GlobalConfig,
			LocalGitConfig:  args.LocalConfig,
			NormalConfig:    normalConfig,
		},
		UnvalidatedConfig: unvalidatedConfig,
	}, finalMessages
}

type NewUnvalidatedConfigArgs struct {
	Access       gitconfig.Access
	ConfigFile   Option[configdomain.PartialConfig]
	DryRun       configdomain.DryRun
	GitVersion   git.Version
	GlobalConfig configdomain.PartialConfig
	LocalConfig  configdomain.PartialConfig
}
