package config

// AliasType defines Git Town commands that can be aliased.
type AliasType string

const (
	AliasTypeAppend         = "append"
	AliasTypeDiffParent     = "diff-parent"
	AliasTypeHack           = "hack"
	AliasTypeKill           = "kill"
	AliasTypeNewPullRequest = "new-pull-request"
	AliasTypePrepend        = "prepend"
	AliasTypePruneBranches  = "prune-branches"
	AliasTypeRenameBranch   = "rename-branch"
	AliasTypeRepo           = "repo"
	AliasTypeShip           = "ship"
	AliasTypeSync           = "sync"
)

// AliasTypes provides all AliasType values.
func AliasTypes() []AliasType {
	return []AliasType{
		AliasTypeAppend,
		AliasTypeDiffParent,
		AliasTypeHack,
		AliasTypeKill,
		AliasTypeNewPullRequest,
		AliasTypePrepend,
		AliasTypePruneBranches,
		AliasTypeRenameBranch,
		AliasTypeRepo,
		AliasTypeShip,
		AliasTypeSync,
	}
}
