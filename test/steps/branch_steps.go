package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/test"
	"github.com/pkg/errors"
)

// BranchSteps defines Cucumber step implementations around Git branches.
func BranchSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^Git Town is now aware of this branch hierarchy$`, func(input *gherkin.DataTable) error {
		gitConfig := git.NewConfiguration(fs.activeScenarioState.gitEnvironment.DeveloperRepo.Dir)
		table := test.DataTable{}
		table.AddRow("BRANCH", "PARENT")
		for _, row := range input.Rows[1:] {
			branch := row.Cells[0].Value
			parentBranch := gitConfig.GetParentBranch(branch)
			table.AddRow(branch, parentBranch)
		}
		diff, errCount := table.Equal(input)
		if errCount > 0 {
			fmt.Println(diff)
			return fmt.Errorf("%d differences", errCount)
		}
		return nil
	})

	suite.Step(`^I am on the "([^"]*)" branch$`, func(branchName string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CheckoutBranch(branchName)
		if err != nil {
			return errors.Wrapf(err, "cannot change to branch %q", branchName)
		}
		return nil
	})

	suite.Step(`^I end up on the "([^"]*)" branch$`, func(expected string) error {
		actual, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CurrentBranch()
		if err != nil {
			return errors.Wrap(err, "cannot determine current branch of developer repo")
		}
		if actual != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})

	suite.Step(`^my code base has the perennial branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreatePerennialBranches(branch1, branch2)
		if err != nil {
			return errors.Wrap(err, "cannot create perennial branches")
		}
		err = fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch1)
		if err != nil {
			return errors.Wrapf(err, "cannot push branch %q", branch1)
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch2)
	})

	suite.Step(`^my repository has a feature branch named "([^"]*)"$`, func(branch string) error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranch(branch)
	})
}
