package config

// Alias defines Git Town commands that can be aliased.
type Alias struct {
	name string
}

func (a Alias) String() string { return a.name }

var (
	AliasAppend         = Alias{"append"}           //nolint:gochecknoglobals
	AliasDiffParent     = Alias{"diff-parent"}      //nolint:gochecknoglobals
	AliasHack           = Alias{"hack"}             //nolint:gochecknoglobals
	AliasKill           = Alias{"kill"}             //nolint:gochecknoglobals
	AliasNewPullRequest = Alias{"new-pull-request"} //nolint:gochecknoglobals
	AliasPrepend        = Alias{"prepend"}          //nolint:gochecknoglobals
	AliasPruneBranches  = Alias{"prune-branches"}   //nolint:gochecknoglobals
	AliasRenameBranch   = Alias{"rename-branch"}    //nolint:gochecknoglobals
	AliasRepo           = Alias{"repo"}             //nolint:gochecknoglobals
	AliasShip           = Alias{"ship"}             //nolint:gochecknoglobals
	AliasSync           = Alias{"sync"}             //nolint:gochecknoglobals
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
