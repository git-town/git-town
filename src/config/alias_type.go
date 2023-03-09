package config

// AliasType defines Git Town commands that can be aliased.
type AliasType string

const (
	AliasTypeAppend         AliasType = "append"
	AliasTypeDiffParent     AliasType = "diff-parent"
	AliasTypeHack           AliasType = "hack"
	AliasTypeKill           AliasType = "kill"
	AliasTypeNewPullRequest AliasType = "new-pull-request"
	AliasTypePrepend        AliasType = "prepend"
	AliasTypePruneBranches  AliasType = "prune-branches"
	AliasTypeRenameBranch   AliasType = "rename-branch"
	AliasTypeRepo           AliasType = "repo"
	AliasTypeShip           AliasType = "ship"
	AliasTypeSync           AliasType = "sync"
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
