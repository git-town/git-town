package config

import (
	"os"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/envconfig"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type UnvalidatedConfig struct {
	CLI               configdomain.PartialConfig // configuration received via CLI flags
	Defaults          NormalConfig               // default values
	Env               configdomain.PartialConfig // environment variables
	File              configdomain.PartialConfig // content of git-town.toml
	GitGlobal         configdomain.PartialConfig // global Git metadata
	GitLocal          configdomain.PartialConfig // local Git metadata
	GitUnscoped       configdomain.PartialConfig // unscoped Git metadata
	NormalConfig      NormalConfig
	UnvalidatedConfig configdomain.UnvalidatedConfigData
}

func (self *UnvalidatedConfig) BranchType(branch gitdomain.LocalBranchName) configdomain.BranchType {
	return self.UnvalidatedConfig.PartialBranchType(branch).GetOr(self.NormalConfig.PartialBranchType(branch))
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

func (self *UnvalidatedConfig) Reload(backend subshelldomain.RunnerQuerier) configdomain.BeginConfigSnapshot {
	globalSnapshot, _ := gitconfig.LoadSnapshot(backend, Some(configdomain.ConfigScopeGlobal), configdomain.UpdateOutdatedNo) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	localSnapshot, _ := gitconfig.LoadSnapshot(backend, Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedNo)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	unscopedSnapshot, _ := gitconfig.LoadSnapshot(backend, None[configdomain.ConfigScope](), configdomain.UpdateOutdatedNo)   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	unscopedGitConfig, _ := NewPartialConfigFromSnapshot(unscopedSnapshot, false, false, nil)
	envConfig, _ := envconfig.Load(envconfig.NewEnvVars(os.Environ()))
	unvalidatedConfig, normalConfig := mergeConfigs(mergeConfigsArgs{
		cli:      configdomain.EmptyPartialConfig(),
		defaults: DefaultNormalConfig(),
		env:      envConfig,
		file:     self.File,
		git:      unscopedGitConfig,
	})
	self.GitUnscoped = unscopedGitConfig
	self.UnvalidatedConfig = unvalidatedConfig
	self.NormalConfig = normalConfig
	return configdomain.BeginConfigSnapshot{
		Global:   globalSnapshot,
		Local:    localSnapshot,
		Unscoped: unscopedSnapshot,
	}
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName, runner subshelldomain.Runner) error {
	self.UnvalidatedConfig.MainBranch = Some(branch)
	return gitconfig.SetMainBranch(runner, branch, configdomain.ConfigScopeLocal)
}

// UnvalidatedBranchesAndTypes provides the types for the given branches.
// This method's name startes with "Unvalidated" to indicate that the types might be incomplete,
// and you should use ValidatedConfig.BranchesAndTypes if possible.
func (self *UnvalidatedConfig) UnvalidatedBranchesAndTypes(branches gitdomain.LocalBranchNames) configdomain.BranchesAndTypes {
	result := make(configdomain.BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.UnvalidatedConfig.PartialBranchType(branch).GetOr(self.NormalConfig.PartialBranchType(branch))
	}
	return result
}

func NewUnvalidatedConfig(args NewUnvalidatedConfigArgs) UnvalidatedConfig {
	unvalidatedConfig, normalConfig := mergeConfigs(mergeConfigsArgs{
		cli:      args.CliConfig,
		defaults: args.Defaults,
		env:      args.EnvConfig,
		file:     args.ConfigFile,
		git:      args.GitUnscoped,
	})
	return UnvalidatedConfig{
		CLI:               args.CliConfig,
		Defaults:          args.Defaults,
		Env:               args.EnvConfig,
		File:              args.ConfigFile,
		GitGlobal:         args.GitGlobal,
		GitLocal:          args.GitLocal,
		GitUnscoped:       args.GitUnscoped,
		NormalConfig:      normalConfig,
		UnvalidatedConfig: unvalidatedConfig,
	}
}

type NewUnvalidatedConfigArgs struct {
	CliConfig     configdomain.PartialConfig
	ConfigFile    configdomain.PartialConfig
	Defaults      NormalConfig
	EnvConfig     configdomain.PartialConfig
	FinalMessages stringslice.Collector
	GitGlobal     configdomain.PartialConfig
	GitLocal      configdomain.PartialConfig
	GitUnscoped   configdomain.PartialConfig
}

func mergeConfigs(args mergeConfigsArgs) (configdomain.UnvalidatedConfigData, NormalConfig) {
	result := configdomain.EmptyPartialConfig()
	result = result.Merge(args.file)
	result = result.Merge(args.git)
	result = result.Merge(args.env)
	result = result.Merge(args.cli)
	return result.ToUnvalidatedConfig(), NewNormalConfigFromPartial(result, args.defaults)
}

type mergeConfigsArgs struct {
	cli      configdomain.PartialConfig
	defaults NormalConfig
	env      configdomain.PartialConfig // configuration data taken from environment variables
	file     configdomain.PartialConfig // data of the configuration file
	git      configdomain.PartialConfig
}
