package configdomain

import (
	"fmt"
	"regexp"

	"github.com/git-town/git-town/v15/internal/cli/colors"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// PerennialRegex contains the "branches.perennial-regex" setting.
type PerennialRegex string

// MatchesBranch indicates whether the given branch matches this PerennialRegex.
func (self PerennialRegex) MatchesBranch(branch gitdomain.LocalBranchName) bool {
	if self == "" {
		return false
	}
	re, err := regexp.Compile(string(self))
	if err != nil {
		fmt.Println(colors.Red().Styled(fmt.Sprintf("Error in perennial regex %q: %s", self, err.Error())))
		return false
	}
	return re.MatchString(branch.String())
}

func (self PerennialRegex) String() string {
	return string(self)
}

func ParsePerennialRegex(value string) Option[PerennialRegex] {
	if value == "" {
		return None[PerennialRegex]()
	}
	return Some(PerennialRegex(value))
}
