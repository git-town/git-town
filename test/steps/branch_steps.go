package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/test"
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
			return fmt.Errorf("found %d differences:\n%s", errCount, diff)
		}
		return nil
	})

	suite.Step(`^I am on the "([^"]*)" branch$`, func(branchName string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CheckoutBranch(branchName)
		if err != nil {
			return fmt.Errorf("cannot change to branch %q: %w", branchName, err)
		}
		return nil
	})

	suite.Step(`^I end up on the "([^"]*)" branch$`, func(expected string) error {
		actual, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CurrentBranch()
		if err != nil {
			return fmt.Errorf("cannot determine current branch of developer repo: %w", err)
		}
		if actual != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})

	suite.Step(`^my repository has a feature branch named "([^"]*)"$`, func(branch string) error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranch(branch)
	})

	suite.Step(`^my repository has the perennial branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreatePerennialBranches(branch1, branch2)
		if err != nil {
			return fmt.Errorf("cannot create perennial branches: %w", err)
		}
		err = fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch1)
		if err != nil {
			return fmt.Errorf("cannot push branch %q: %w", branch1, err)
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch2)
	})
}
