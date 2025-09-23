package setup

import (
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
)

type Data struct {
	Backend       subshelldomain.RunnerQuerier
	Config        config.UnvalidatedConfig
	Git           git.Commands
	Inputs        dialogcomponents.Inputs
	LocalBranches gitdomain.LocalBranchNames
	Remotes       gitdomain.Remotes
	Snapshot      configdomain.BeginConfigSnapshot
}
