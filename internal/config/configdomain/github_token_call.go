package configdomain

import (
	"errors"

	"github.com/kballard/go-shellquote"
)

// the shell call to load the GitHubToken from an external application
type GitHubTokenCall string

// provides this GitHubTokenCall in a callable format
func (self GitHubTokenCall) CallableFormat() (executable string, args []string, err error) {
	words, err := shellquote.Split(self.String())
	if err != nil {
		return "", []string{}, err
	}
	if len(words) == 0 {
		return "", []string{}, errors.New("github token call is empty")
	}
	return words[0], words[1:], nil
}

func (self GitHubTokenCall) String() string {
	return string(self)
}
