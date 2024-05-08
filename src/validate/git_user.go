package validate

import (
	"errors"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

func GitUser(config configdomain.UnvalidatedConfig) (configdomain.GitUserEmail, configdomain.GitUserName, error) {
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
