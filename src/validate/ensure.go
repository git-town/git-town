package validate

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

// CobraEnsure wraps ensure into a Cobra-compatible format.
func CobraEnsure(repo *git.ProdRepo, validators ...validationCondition) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return ensure(repo, validators...)
	}
}

// ensure checks that the given repo conforms to the given validation conditions.
func ensure(repo *git.ProdRepo, validators ...validationCondition) error {
	for _, validator := range validators {
		if err := validator(repo); err != nil {
			return err
		}
	}
	return nil
}

// validationCondition verifies that the given Git repo conforms to a particular condition.
type validationCondition func(*git.ProdRepo) error
