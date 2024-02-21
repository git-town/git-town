package configdomain

import (
	"fmt"
	"regexp"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
)

// PerennialRegex contains the "branches.perennial-regex" setting.
type PerennialRegex string

// MatchBranch indicates whether the given branch matches this PerennialRegex.
func (self PerennialRegex) MatchBranch(branch gitdomain.LocalBranchName) bool {
	if self == "" {
		return false
	}
	re, err := regexp.Compile(string(self))
	if err != nil {
		fmt.Println(components.Red().Styled(fmt.Sprintf("Error in perennial regex %q: %s", self, err.Error())))
		return false
	}
	return re.MatchString(branch.String())
}

func (self PerennialRegex) String() string {
	return string(self)
}

func NewPerennialRegexRef(value string) *PerennialRegex {
	result := PerennialRegex(value)
	return &result
}
