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

func NewCreatePrototypeBranches(value, source string) (CreatePrototypeBranches, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return CreatePrototypeBranches(true), fmt.Errorf(messages.ValueInvalid, source, value)
	}
	result := CreatePrototypeBranches(parsed)
	return result, nil
}

func NewCreatePrototypeBranchesFromGitConfig(configStr string) (Option[CreatePrototypeBranches], error) {
	if configStr == "" {
		return None[CreatePrototypeBranches](), nil
	}
	result, err := gohacks.ParseBool(configStr)
	return Some(CreatePrototypeBranches(result)), err
}
