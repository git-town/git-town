package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v14/src/gohacks"
)

// GiteaToken is a bearer token to use with the Gitea API.
type GiteaToken string

func (self GiteaToken) String() string {
	return string(self)
}

func NewGiteaToken(value string) GiteaToken {
	value = strings.TrimSpace(value)
	if value == "" {
		panic("empty Gitea token")
	}
	return GiteaToken(value)
}

func NewGiteaTokenOption(value string) gohacks.Option[GiteaToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return gohacks.NewOptionNone[GiteaToken]()
	}
	return gohacks.NewOption(NewGiteaToken(value))
}
