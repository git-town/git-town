package steps

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/test"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
)

// BranchSteps defines Cucumber step implementations around Git branches.
// nolint:funlen
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
		diff, errCount := table.EqualGherkin(input)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branch hierarchy\n\n", errCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see the diff above")
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

	suite.Step(`^I (?:end up|am still) on the "([^"]*)" branch$`, func(expected string) error {
		actual, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CurrentBranch()
		if err != nil {
			return fmt.Errorf("cannot determine current branch of developer repo: %w", err)
		}
		if actual != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})

	suite.Step(`^my coworker has a feature branch named "([^"]*)"$`, func(branch string) error {
		return fs.activeScenarioState.gitEnvironment.OriginRepo.CreateBranch(branch)
	})

	suite.Step(`^my repository has a feature branch named "([^"]*)"$`, func(branch string) error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranch(branch)
	})

	suite.Step(`^my repository has a feature branch named "([^"]+)" as a child of "([^"]+)"$`, func(childBranch, parentBranch string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateChildFeatureBranch(childBranch, parentBranch)
		if err != nil {
			return fmt.Errorf("cannot create feature branch %q: %w", childBranch, err)
		}
		return nil
	})

	suite.Step(`^my repository has the feature branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranch(branch1)
		if err != nil {
			return err
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranch(branch2)
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

	suite.Step(`^the perennial branches are configured as "([^"]+)"$`, func(name string) error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration().AddToPerennialBranches(name)
		return nil
	})
}
