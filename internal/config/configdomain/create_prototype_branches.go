package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v15/internal/gohacks"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// whether all created branches should be prototype
type CreatePrototypeBranches bool

func (self CreatePrototypeBranches) IsTrue() bool {
	return bool(self)
}

func (self CreatePrototypeBranches) String() string {
	return strconv.FormatBool(bool(self))
}

// deserializes the given Git configuration value into a CreatePrototypeBranches instance
func ParseCreatePrototypeBranches(value, source string) (Option[CreatePrototypeBranches], error) {
	parsedOpt, err := gohacks.ParseBool(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(CreatePrototypeBranches(parsed)), err
	}
	return None[CreatePrototypeBranches](), err
}
