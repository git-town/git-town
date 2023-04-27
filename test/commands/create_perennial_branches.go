package commands

import "fmt"

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func CreatePerennialBranches(repo Repo, names ...string) error {
	for _, name := range names {
		err := CreateBranch(repo, name, "main")
		if err != nil {
			return fmt.Errorf("cannot create perennial branch %q in repo %q: %w", name, repo.Dir(), err)
		}
	}
	return repo.Config().AddToPerennialBranches(names...)
}
