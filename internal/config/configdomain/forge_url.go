package configdomain

import (
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// URL that Git Town should assume for Git server, used only to simulate a forge in tests
type ForgeURL string

func (self ForgeURL) String() string {
	return string(self)
}

func ParseRemoteForTesting(value stringss.Trimmed, _ string) (Option[ForgeURL], error) {
	if value == "" {
		return None[ForgeURL](), nil
	}
	return Some(ForgeURL(value)), nil
}
