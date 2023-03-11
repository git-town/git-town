package validate

import "github.com/git-town/git-town/v7/src/git"

// KnowsBranchesAncestry asserts that the entire ancestry for all given branches
// is known to Git Town.
// Missing ancestry information is queried from the user.
func KnowsBranchesAncestry(branches []string, repo *git.ProdRepo) error {
	for _, branch := range branches {
		err := KnowsBranchAncestry(branch, repo.Config.MainBranch(), repo)
		if err != nil {
			return err
		}
	}
	return nil
}
