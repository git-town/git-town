package configdomain

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

// ValidatedConfigData is Git Town configuration where all essential values are guaranteed to exist and have meaningful values.
// This is ensured by querying from the user if needed.
type ValidatedConfigData struct {
	MainBranch gitdomain.LocalBranchName
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *ValidatedConfigData) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.MainBranch
}
