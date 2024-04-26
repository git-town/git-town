package configdomain

import (
	"errors"
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

type GitUserEmail string

func (self GitUserEmail) String() string {
	return string(self)
}

func NewGitUserEmail(value string) (GitUserEmail, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New(messages.GitUserEmailMissing)
	}
	return GitUserEmail(value), nil
}

func NewGitUserEmailOption(value string) (Option[GitUserEmail], error) {
	value = strings.TrimSpace(value)
	email, err := NewGitUserEmail(value)
	if err != nil {
		return None[GitUserEmail](), err
	}
	return Some(email), nil
}
