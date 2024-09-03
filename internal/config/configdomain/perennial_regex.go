package configdomain

import (
	"fmt"
	"regexp"

	"github.com/git-town/git-town/v16/internal/cli/colors"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// PerennialRegex contains the "branches.perennial-regex" setting.
type PerennialRegex struct {
	regex Option[*regexp.Regexp]
	text  string
}

// MatchesBranch indicates whether the given branch matches this PerennialRegex.
func (self PerennialRegex) MatchesBranch(branch gitdomain.LocalBranchName) bool {
	if self.text == "" {
		return false
	}
	re, err := regexp.Compile(self.text)
	if err != nil {
		fmt.Println(colors.Red().Styled(fmt.Sprintf("Error in perennial regex %q: %s", self, err.Error())))
		return false
	}
	return re.MatchString(branch.String())
}

func (self PerennialRegex) Regex() *regexp.Regexp {
	if self.regex.IsNone() {
		self.regex = Some(regexp.MustCompile(self.text))
	}
	return self.regex.GetOrPanic()
}

func (self PerennialRegex) String() string {
	return self.text
}

func ParsePerennialRegex(value string) Option[PerennialRegex] {
	if value == "" {
		return None[PerennialRegex]()
	}
	return Some(PerennialRegex{
		regex: None[*regexp.Regexp](),
		text:  value,
	})
}
