package config

// Alias defines Git Town commands that can be aliased.
type Alias string

const (
	AliasAppend         Alias = "append"
	AliasDiffParent     Alias = "diff-parent"
	AliasHack           Alias = "hack"
	AliasKill           Alias = "kill"
	AliasNewPullRequest Alias = "new-pull-request"
	AliasPrepend        Alias = "prepend"
	AliasPruneBranches  Alias = "prune-branches"
	AliasRenameBranch   Alias = "rename-branch"
	AliasRepo           Alias = "repo"
	AliasShip           Alias = "ship"
	AliasSync           Alias = "sync"
)

// Aliases provides all AliasType values.
func Aliases() []Alias {
	return []Alias{
		AliasAppend,
		AliasDiffParent,
		AliasHack,
		AliasKill,
		AliasNewPullRequest,
		AliasPrepend,
		AliasPruneBranches,
		AliasRenameBranch,
		AliasRepo,
		AliasShip,
		AliasSync,
	}
}
