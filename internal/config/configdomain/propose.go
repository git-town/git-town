package configdomain

import (
	"fmt"
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// indicates whether a Git Town command should propose the branch it creates
type Propose bool

func (self Propose) ShouldPropose() bool {
	return bool(self)
}

const (
	ProposeTitleFirst  ProposeTitle = "first"  // use the title of the first commit
	ProposeTitleNative ProposeTitle = "native" // delegate to external tooling (default)
	ProposeTitleSelect ProposeTitle = "select" // let the user select from commit titles
)

type ProposeTitle string

func (self ProposeTitle) String() string {
	return string(self)
}

func ParseProposeTitle(text string, source string) (Option[ProposeTitle], error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return None[ProposeTitle](), nil
	}
	text = strings.ToLower(text)
	for _, proposeTitle := range ProposeTitles() {
		if proposeTitle.String() == text {
			return Some(proposeTitle), nil
		}
	}
	return None[ProposeTitle](), fmt.Errorf("unknown propose.title value %q (source: %s), allowed values: first, select, native", text, source)
}

func ProposeTitles() []ProposeTitle {
	return []ProposeTitle{
		ProposeTitleFirst,
		ProposeTitleNative,
		ProposeTitleSelect,
	}
}
