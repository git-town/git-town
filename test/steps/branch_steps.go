package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/src/git"
	"github.com/pkg/errors"
)

// BranchSteps defines Cucumber step implementations around Git branches.
func BranchSteps(suite *godog.Suite, fs *FeatureState) {

	suite.Step(`^Git Town is now aware of this branch hierarchy$`, func(data *gherkin.DataTable) error {
		branches := git.Config.GetAncestorBranches("new-child")
		fmt.Println(branches)
		return nil
	})

	suite.Step(`^I am on the "([^"]*)" branch$`, func(branchName string) error {
		output, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Run("git", "checkout", branchName)
		if err != nil {
			return errors.Wrapf(err, "cannot change to branch %q: %s", branchName, output)
		}
		return nil
	})

	suite.Step(`^I end up on the "([^"]*)" branch$`, func(expected string) error {
		actual, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CurrentBranch()
		if err != nil {
			return err
		}
		if actual != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})
}
