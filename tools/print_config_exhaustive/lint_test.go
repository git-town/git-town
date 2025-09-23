package main_test

import (
	"testing"

	main "github.com/git-town/git-town/tools/print_config_exhaustive"
	"github.com/shoenig/test/must"
)

func TestDefinitionFields(t *testing.T) {
	t.Parallel()

	//nolint:dupword
	give := `
	package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/git-town/git-town/v22/pkg/set"
)

// configuration settings that exist in both UnvalidatedConfig and ValidatedConfig
type NormalConfig struct {
	Aliases                  Aliases
	BitbucketAppPassword     Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername        Option[forgedomain.BitbucketUsername]
	BranchTypeOverrides      BranchTypeOverrides
	ContributionRegex        Option[ContributionRegex]
	DevRemote                gitdomain.Remote
	FeatureRegex             Option[FeatureRegex]
	ForgejoToken            Option[forgedomain.ForgejoToken]
	ForgeType                Option[forgedomain.ForgeType] // None = auto-detect
	GitHubConnectorType      Option[forgedomain.GitHubConnectorType]
	GitHubToken              Option[forgedomain.GitHubToken]
	GitLabConnectorType      Option[forgedomain.GitLabConnectorType]
	GitLabToken              Option[forgedomain.GitLabToken]
	GiteaToken               Option[forgedomain.GiteaToken]
	HostingOriginHostname    Option[HostingOriginHostname]
	Lineage                  Lineage
	NewBranchType            Option[BranchType]
	ObservedRegex            Option[ObservedRegex]
	Offline                  Offline
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[PerennialRegex]
	PushHook                 PushHook
	ShareNewBranches         ShareNewBranches
	ShipDeleteTrackingBranch ShipDeleteTrackingBranch
	ShipStrategy             ShipStrategy
	SyncFeatureStrategy      SyncFeatureStrategy
	SyncPerennialStrategy    SyncPerennialStrategy
	SyncPrototypeStrategy    SyncPrototypeStrategy
	SyncTags                 SyncTags
	SyncUpstream             SyncUpstream
	UnknownBranchType        BranchType
}
`
	have := main.FindDefinedFields(give)
	want := []string{
		"Aliases",
		"BitbucketAppPassword",
		"BitbucketUsername",
		"BranchTypeOverrides",
		"ContributionRegex",
		"DevRemote",
		"FeatureRegex",
		"ForgejoToken",
		"ForgeType",
		"GitHubConnectorType",
		"GitHubToken",
		"GitLabConnectorType",
		"GitLabToken",
		"GiteaToken",
		"HostingOriginHostname",
		"Lineage",
		"NewBranchType",
		"ObservedRegex",
		"Offline",
		"PerennialBranches",
		"PerennialRegex",
		"PushHook",
		"ShareNewBranches",
		"ShipDeleteTrackingBranch",
		"ShipStrategy",
		"SyncFeatureStrategy",
		"SyncPerennialStrategy",
		"SyncPrototypeStrategy",
		"SyncTags",
		"SyncUpstream",
		"UnknownBranchType",
	}
	must.Eq(t, want, have)
}

func TestFindUnprinted(t *testing.T) {
	t.Parallel()
	fields := []string{
		"ContributionRegex",
		"FeatureRegex",
		"MainBranch",
		"ObservedRegex",
		"ForgeType",
	}
	whiteList := []string{
		"ParkedBranch",
	}
	printText := `
	print.Entry("contribution regex", format.OptionalStringerSetting(config.NormalConfig.ContributionRegex))
	print.Entry("feature regex", format.OptionalStringerSetting(config.NormalConfig.FeatureRegex))
	print.Entry("main branch", format.OptionalStringerSetting(config.UnvalidatedConfig.MainBranch))
	print.Entry("observed branches", format.BranchNames(config.NormalConfig.PartialBranchesOfType(configdomain.BranchTypeObservedBranch)))
	`
	have := main.FindUnprintedFields(fields, printText, whiteList)
	want := []string{
		"ObservedRegex",
		"ForgeType",
	}
	must.Eq(t, want, have)
}

func TestParsePrintFile(t *testing.T) {
	t.Parallel()
	give := `
// Package config implements Git Town's "config" command.
package config

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/format"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/spf13/cobra"
)

const configDesc = "Display your Git Town configuration"

func RootCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	configCmd := cobra.Command{
		Use:     "config",
		GroupID: cmdhelpers.GroupIDSetup,
		Args:    cobra.NoArgs,
		Short:   configDesc,
		Long:    cmdhelpers.Long(configDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeDisplayConfig(verbose)
		},
	}
	addVerboseFlag(&configCmd)
	configCmd.AddCommand(getParentCommand())
	configCmd.AddCommand(removeConfigCommand())
	configCmd.AddCommand(SetupCommand())
	return &configCmd
}

func executeDisplayConfig(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	printConfig(repo.UnvalidatedConfig)
	return nil
}

func printConfig(config config.UnvalidatedConfig) {
	fmt.Println()
	print.Header("Branches")
	print.Entry("contribution branches", format.BranchNames(config.NormalConfig.PartialBranchesOfType(configdomain.BranchTypeContributionBranch)))
	print.Entry("contribution regex", format.OptionalStringerSetting(config.NormalConfig.ContributionRegex))
	print.Entry("feature regex", format.OptionalStringerSetting(config.NormalConfig.FeatureRegex))
	fmt.Println()
	if config.NormalConfig.Lineage.Len() > 0 {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(config.NormalConfig.Lineage))
	}
}
`
	have := main.FindPrintFunc(give)
	want := `
	fmt.Println()
	print.Header("Branches")
	print.Entry("contribution branches", format.BranchNames(config.NormalConfig.PartialBranchesOfType(configdomain.BranchTypeContributionBranch)))
	print.Entry("contribution regex", format.OptionalStringerSetting(config.NormalConfig.ContributionRegex))
	print.Entry("feature regex", format.OptionalStringerSetting(config.NormalConfig.FeatureRegex))
	fmt.Println()
	if config.NormalConfig.Lineage.Len() > 0 {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(config.NormalConfig.Lineage))
	`
	must.EqOp(t, want, have)
}
