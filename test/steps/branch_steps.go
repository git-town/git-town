package steps

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/test"
)

// BranchSteps defines Cucumber step implementations around Git branches.
// nolint:funlen,gocognit
func BranchSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^all branches are now synchronized$`, func() error {
		outOfSync, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.HasBranchesOutOfSync()
		if err != nil {
			return err
		}
		if outOfSync {
			return fmt.Errorf("expected no branches out of sync")
		}
		return nil
	})

	suite.Step(`^Git Town is (?:now|still) aware of this branch hierarchy$`, func(input *messages.PickleStepArgument_PickleTable) error {
		gitConfig := git.NewConfiguration(fs.activeScenarioState.gitEnvironment.DeveloperShell)
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

	suite.Step(`^Git Town now has no branch hierarchy information$`, func() error {
		has := git.NewConfiguration(fs.activeScenarioState.gitEnvironment.DeveloperShell).HasBranchInformation()
		if has {
			return fmt.Errorf("unexpected Git Town branch hierarchy information")
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

	suite.Step(`^I am on the "([^"]*)" branch with "([^"]*)" as the previous Git branch$`, func(current, previous string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CheckoutBranch(previous)
		if err != nil {
			return err
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CheckoutBranch(current)
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

	suite.Step(`^I don't have a main branch name configured$`, func() error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.DeleteMainBranchConfiguration()
	})

	suite.Step(`^my code base has a feature branch named "([^"]*)"$`, func(name string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranch(name)
		if err != nil {
			return err
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(name)
	})

	suite.Step(`^my code base has a feature branch named "([^"]*)" as a child of "([^"]*)"$`, func(branch, parent string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateChildFeatureBranch(branch, parent)
		if err != nil {
			return err
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch)
	})

	suite.Step(`^my (?:coworker|origin) has a feature branch named "([^"]*)"$`, func(branch string) error {
		return fs.activeScenarioState.gitEnvironment.OriginRepo.CreateBranch(branch, "main")
	})

	suite.Step(`^my repository has a branch "([^"]*)"$`, func(branch string) error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateBranch(branch, "main")
	})

	suite.Step(`^my repository has a (local )?feature branch named "([^"]*)"$`, func(localStr, branch string) error {
		isLocal := localStr != ""
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranch(branch)
		if err != nil {
			return err
		}
		if !isLocal {
			return fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch)
		}
		return nil
	})

	suite.Step(`^my repository has a feature branch named "([^"]*)" with no parent$`, func(branch string) error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranchNoParent(branch)
	})

	suite.Step(`^my repository has a feature branch named "([^"]+)" as a child of "([^"]+)"$`, func(childBranch, parentBranch string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateChildFeatureBranch(childBranch, parentBranch)
		if err != nil {
			return fmt.Errorf("cannot create feature branch %q: %w", childBranch, err)
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(childBranch)
	})

	suite.Step(`^my repository has the branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateBranch(branch1, "main")
		if err != nil {
			return err
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateBranch(branch2, "main")
	})

	suite.Step(`^my repository has the (local )?feature branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1, branch2 string) error {
		isLocal := localStr != ""
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranch(branch1)
		if err != nil {
			return err
		}
		err = fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFeatureBranch(branch2)
		if err != nil {
			return err
		}
		if !isLocal {
			err = fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch1)
			if err != nil {
				return err
			}
			return fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch2)
		}
		return nil
	})

	suite.Step(`^my repository has the perennial branch "([^"]+)"`, func(branch1 string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreatePerennialBranches(branch1)
		if err != nil {
			return fmt.Errorf("cannot create perennial branches: %w", err)
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch1)
	})

	suite.Step(`^my repository has the (local )?perennial branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1, branch2 string) error {
		isLocal := localStr != ""
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreatePerennialBranches(branch1, branch2)
		if err != nil {
			return fmt.Errorf("cannot create perennial branches: %w", err)
		}
		if !isLocal {
			err = fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch1)
			if err != nil {
				return err
			}
			return fs.activeScenarioState.gitEnvironment.DeveloperRepo.PushBranch(branch2)
		}
		return nil
	})

	suite.Step(`^the "([^"]*)" branch gets deleted on the remote$`, func(name string) error {
		return fs.activeScenarioState.gitEnvironment.OriginRepo.RemoveBranch(name)
	})

	suite.Step(`^the existing branches are$`, func(table *messages.PickleStepArgument_PickleTable) error {
		existing, err := fs.activeScenarioState.gitEnvironment.Branches()
		if err != nil {
			return err
		}
		// remove the master branch from the remote since it exists only as a performance optimization
		existing.RemoveText("master, ")
		existing.RemoveText(", master")
		diff, errCount := existing.EqualGherkin(table)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branches\n\n", errCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the previous Git branch is (?:now|still) "([^"]*)"$`, func(want string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CheckoutBranch("-")
		if err != nil {
			return err
		}
		have, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CurrentBranch()
		if err != nil {
			return err
		}
		if have != want {
			return fmt.Errorf("expected previous branch %q but got %q", want, have)
		}
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CheckoutBranch("-")
	})
}
