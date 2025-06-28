package validate

import (
	"errors"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

func GitUser(config configdomain.UnvalidatedConfigData) (gitdomain.GitUserEmail, gitdomain.GitUserName, error) {
	gitUserEmail, hasGitUserEmail := config.GitUserEmail.Get()
	if !hasGitUserEmail {
		return "", "", errors.New(messages.GitUserEmailMissing)
	}
	gitUserName, hasGitUserName := config.GitUserName.Get()
	if !hasGitUserName {
		return "", "", errors.New(messages.GitUserNameMissing)
	}
	return gitUserEmail, gitUserName, nil
}
