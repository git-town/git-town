package configdomain

import (
	"errors"
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

type GitUserName string

func (self GitUserName) String() string {
	return string(self)
}

func NewGitUserName(value string) (GitUserName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New(messages.GitUserNameMissing)
	}
	return GitUserName(value), nil
}

func NewGitUserNameOption(value string) (Option[GitUserName], error) {
	value = strings.TrimSpace(value)
	name, err := NewGitUserName(value)
	if err != nil {
		return None[GitUserName](), err
	}
	return Some(name), nil
}
