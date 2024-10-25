package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
)

// ValidatedConfigData is Git Town configuration where all essential values are guaranteed to exist and have meaningful values.
// This is ensured by querying from the user if needed.
type ValidatedConfigData struct {
	GitUserEmail GitUserEmail
	GitUserName  GitUserName
	MainBranch   gitdomain.LocalBranchName
}

// Author provides the locally Git configured user.
func (self *ValidatedConfigData) Author() gitdomain.Author {
	email := self.GitUserEmail
	name := self.GitUserName
	return gitdomain.Author(fmt.Sprintf("%s <%s>", name, email))
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *ValidatedConfigData) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.MainBranch
}
