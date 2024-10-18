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
	Config          Mutable[configdomain.UnvalidatedConfig] // the merged configuration data
	ConfigFile      Option[configdomain.PartialConfig]      // content of git-town.toml, nil = no config file exists
	GitConfig       gitconfig.Access                        // access to the Git configuration settings
	GitVersion      git.Version                             // the version of the installed Git executable
	GlobalGitConfig configdomain.PartialConfig              // content of the global Git configuration
}

func NewUnvalidatedConfig(args NewUnvalidatedConfigArgs) (NormalConfig, UnvalidatedConfig, stringslice.Collector) {
	normalConfigData, unvalidatedConfigData := configdomain.NewUnvalidatedConfig(args.ConfigFile, args.GlobalConfig, args.LocalConfig)
	normalConfig := NormalConfig{
		Config:         NewMutable(&normalConfigData),
		DryRun:         false,
		GitConfig:      gitconfig.Access{},
		LocalGitConfig: configdomain.PartialConfig{},
	}
	unvalidatedConfig := UnvalidatedConfig{
		Config:          NewMutable(&unvalidatedConfigData),
		ConfigFile:      args.ConfigFile,
		GitConfig:       args.Access,
		GitVersion:      args.GitVersion,
		GlobalGitConfig: args.GlobalConfig,
	}
	finalMessages := stringslice.NewCollector()
	return normalConfig, unvalidatedConfig, finalMessages
}

func (self *UnvalidatedConfig) RemoveMainBranch() {
	_ = self.GitConfig.RemoveLocalConfigValue(configdomain.KeyMainBranch)
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *UnvalidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.Config.Value.MainBranch = Some(branch)
	return self.GitConfig.SetLocalConfigValue(configdomain.KeyMainBranch, branch.String())
}

type NewUnvalidatedConfigArgs struct {
	Access       gitconfig.Access
	ConfigFile   Option[configdomain.PartialConfig]
	DryRun       configdomain.DryRun
	GitVersion   git.Version
	GlobalConfig configdomain.PartialConfig
	LocalConfig  configdomain.PartialConfig
}
