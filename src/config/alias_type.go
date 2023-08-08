package config

// AliasType defines Git Town commands that can be aliased.
type AliasType struct {
	name string
}

func (a AliasType) String() string { return a.name }

var (
	AliasTypeAppend         = AliasType{"append"}           //nolint:gochecknoglobals
	AliasTypeDiffParent     = AliasType{"diff-parent"}      //nolint:gochecknoglobals
	AliasTypeHack           = AliasType{"hack"}             //nolint:gochecknoglobals
	AliasTypeKill           = AliasType{"kill"}             //nolint:gochecknoglobals
	AliasTypeNewPullRequest = AliasType{"new-pull-request"} //nolint:gochecknoglobals
	AliasTypePrepend        = AliasType{"prepend"}          //nolint:gochecknoglobals
	AliasTypePruneBranches  = AliasType{"prune-branches"}   //nolint:gochecknoglobals
	AliasTypeRenameBranch   = AliasType{"rename-branch"}    //nolint:gochecknoglobals
	AliasTypeRepo           = AliasType{"repo"}             //nolint:gochecknoglobals
	AliasTypeShip           = AliasType{"ship"}             //nolint:gochecknoglobals
	AliasTypeSync           = AliasType{"sync"}             //nolint:gochecknoglobals
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
