package configdomain

// Alias defines Git Town commands that can be aliased.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type Alias struct {
	name string
}

func (self Alias) String() string { return self.name }

var (
	AliasAppend       = Alias{"append"}        //nolint:gochecknoglobals
	AliasDiffParent   = Alias{"diff-parent"}   //nolint:gochecknoglobals
	AliasHack         = Alias{"hack"}          //nolint:gochecknoglobals
	AliasKill         = Alias{"kill"}          //nolint:gochecknoglobals
	AliasPrepend      = Alias{"prepend"}       //nolint:gochecknoglobals
	AliasPropose      = Alias{"propose"}       //nolint:gochecknoglobals
	AliasRenameBranch = Alias{"rename-branch"} //nolint:gochecknoglobals
	AliasRepo         = Alias{"repo"}          //nolint:gochecknoglobals
	AliasShip         = Alias{"ship"}          //nolint:gochecknoglobals
	AliasSync         = Alias{"sync"}          //nolint:gochecknoglobals
)

// Aliases provides all AliasType values.
func Aliases() []Alias {
	return []Alias{
		AliasAppend,
		AliasDiffParent,
		AliasHack,
		AliasKill,
		AliasPrepend,
		AliasPropose,
		AliasRenameBranch,
		AliasRepo,
		AliasShip,
		AliasSync,
	}
}
