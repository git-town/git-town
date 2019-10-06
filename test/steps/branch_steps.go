package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/pkg/errors"
)

// BranchSteps defines Cucumber step implementations around Git branches.
func BranchSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I am on the "([^"]*)" branch$`, func(branchName string) error {
		output, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Run("git", "checkout", branchName)
		if err != nil {
			return errors.Wrapf(err, "cannot change to branch %q: %s", branchName, output)
		}
		if fs.activeScenarioState.debug {
			fmt.Println(output)
		}
		return nil
	})
}
