package forgedomain

import (
	"errors"
	"strings"

	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/kballard/go-shellquote"
)

// GitHubTokenScript is a shell script that when running provides the GitHubToken.
type GitHubTokenScript string

func (self GitHubTokenScript) Args() (binary string, args []string, err error) {
	parts, err := shellquote.Split(self.String())
	if err != nil {
		return "", []string{}, err
	}
	if len(parts) == 0 {
		return "", []string{}, errors.New("cannot split empty GitHubTokenScript")
	}
	return parts[0], parts[1:], nil
}

func (self GitHubTokenScript) String() string {
	return string(self)
}

func (self GitHubTokenScript) Load(querier subshelldomain.Querier) (Option[GitHubToken], error) {
	binary, args, err := self.Args()
	if err != nil {
		return None[GitHubToken](), err
	}
	output, err := querier.QueryTrim(binary, args...)
	return ParseGitHubToken(output), err
}

func NewGitHubTokenScript(text string) (GitHubTokenScript, error) {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return "", errors.New("entered GitHubTokenScript cannot be empty")
	}
	return GitHubTokenScript(text), nil
}
