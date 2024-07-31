package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

// whether all created branches should be prototype
type CreatePrototypeBranches bool

func (self CreatePrototypeBranches) Bool() bool {
	return bool(self)
}

func (self CreatePrototypeBranches) IsTrue() bool {
	return self.Bool()
}

func (self CreatePrototypeBranches) String() string {
	return strconv.FormatBool(self.Bool())
}

func NewCreatePrototypeBranchesFromGitConfig(valueStr, source string) (Option[CreatePrototypeBranches], error) {
	if valueStr == "" {
		return None[CreatePrototypeBranches](), nil
	}
	valueBool, err := gohacks.ParseBool(valueStr)
	if err != nil {
		return None[CreatePrototypeBranches](), fmt.Errorf(messages.ValueInvalid, source, valueStr)
	}
	return Some(CreatePrototypeBranches(valueBool)), nil
}
