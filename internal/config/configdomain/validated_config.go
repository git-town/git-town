package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
)

// ValidatedConfig is Git Town configuration where all essential values are guaranteed to exist and have meaningful values.
// This is ensured by querying from the user if needed.
// TODO: rename to ValidatedConfigData
type ValidatedConfig struct {
	GitUserEmail GitUserEmail
	GitUserName  GitUserName
	MainBranch   gitdomain.LocalBranchName
}

// Author provides the locally Git configured user.
func (self *ValidatedConfig) Author() gitdomain.Author {
	email := self.GitUserEmail
	name := self.GitUserName
	return gitdomain.Author(fmt.Sprintf("%s <%s>", name, email))
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *ValidatedConfig) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.MainBranch
}
