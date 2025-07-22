package setup

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
)

type Data struct {
	Backend        subshelldomain.RunnerQuerier
	Config         config.UnvalidatedConfig
	ConfigSnapshot undoconfig.ConfigSnapshot
	DialogInputs   dialogcomponents.TestInputs
	Git            git.Commands
	LocalBranches  gitdomain.LocalBranchNames
	Remotes        gitdomain.Remotes
}
