package config

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/config/gitconfig"
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

type Runner interface {
	Query(executable string, args ...string) (string, error)
	Run(executable string, args ...string) error
}

type UnvalidatedConfig struct {
	UnvalidatedConfig configdomain.UnvalidatedConfig
	NormalConfig      NormalConfig
	GitConfig         gitconfig.Access // access to the Git configuration settings
	GitVersion        git.Version      // the version of the installed Git executable
}

func NewUnvalidatedConfig(args NewUnvalidatedConfigArgs) (UnvalidatedConfig, stringslice.Collector) {
	config := configdomain.NewUnvalidatedConfig(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	finalMessages := stringslice.NewCollector()
	return UnvalidatedConfig{
		UnvalidatedConfig: config,
		GitConfig:         args.Access,
		GitVersion:        args.GitVersion,
	}, finalMessages
}

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *UnvalidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return self.UnvalidatedConfig.IsMainBranch(branch) || self.NormalConfig.Config.IsPerennialBranch(branch)
}

func (self *UnvalidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	if mainBranch, hasMainBranch := self.UnvalidatedConfig.MainBranch.Get(); hasMainBranch {
		return append(gitdomain.LocalBranchNames{mainBranch}, self.NormalConfig.Config.PerennialBranches...)
	}
	return self.NormalConfig.Config.PerennialBranches
}

func (self *UnvalidatedConfig) RemoveMainBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(configdomain.KeyMainBranch)
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.UnvalidatedConfig.MainBranch = Some(branch)
	return self.GitConfig.SetLocalConfigValue(configdomain.KeyMainBranch, branch.String())
}

// UnvalidatedBranchesAndTypes provides the types for the given branches.
// This method's name startes with "Unvalidated" to indicate that the types might be incomplete,
// and you should use ValidatedConfig.BranchesAndTypes if possible.
func (self *UnvalidatedConfig) UnvalidatedBranchesAndTypes(branches gitdomain.LocalBranchNames) configdomain.BranchesAndTypes {
	result := make(configdomain.BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.UnvalidatedConfig.PartialBranchType(branch).GetOrElse(self.NormalConfig.Config.PartialBranchType(branch))
	}
	return result
}

type NewUnvalidatedConfigArgs struct {
	Access       gitconfig.Access
	ConfigFile   Option[configdomain.PartialConfig]
	DryRun       configdomain.DryRun
	GitVersion   git.Version
	GlobalConfig configdomain.PartialConfig
	LocalConfig  configdomain.PartialConfig
}
