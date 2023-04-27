package commands

import "fmt"

// CreatePerennialBranches creates perennial branches with the given names in this repository.
func CreatePerennialBranches(r *TestCommands, names ...string) error {
	for _, name := range names {
		err := CreateBranch(r, name, "main")
		if err != nil {
			return fmt.Errorf("cannot create perennial branch %q in repo %q: %w", name, r.WorkingDir, err)
		}
	}
	return r.Config.AddToPerennialBranches(names...)
}
