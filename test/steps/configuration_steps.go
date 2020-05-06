package steps

import (
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
)

// ConfigurationSteps defines Cucumber step implementations around configuration.
// nolint:funlen,gocognit
func ConfigurationSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^Git Town is no longer configured for this repository$`, func() error {
		res, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.HasGitTownConfigNow()
		if err != nil {
			return err
		}
		if res {
			return fmt.Errorf("unexpected Git Town configuration")
		}
		return nil
	})

	suite.Step(`^I haven't configured Git Town yet$`, func() error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(false).DeleteMainBranchConfiguration()
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(false).DeletePerennialBranchConfiguration()
		return nil
	})

	suite.Step(`^my repo is now configured with no perennial branches$`, func() error {
		branches := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetPerennialBranches()
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is now (true|false)$`, func(text string) error {
		want, err := strconv.ParseBool(text)
		if err != nil {
			return fmt.Errorf("cannot parse %q into bool: %w", text, err)
		}
		have := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).ShouldNewBranchPush()
		if have != want {
			return fmt.Errorf("expected global new-branch-push-flag to be %t, but was %t", want, have)
		}
		return nil
	})

	suite.Step(`^the global new-branch-push-flag configuration is (true|false)$`, func(text string) error {
		b, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		_ = fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(false).SetNewBranchPush(b, true)
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is (true|false)$`, func(value string) error {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot parse %q into bool: %w", value, err)
		}
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(false).SetNewBranchPush(b, false)
		return nil
	})

	suite.Step(`^the main branch is configured as "([^"]+)"$`, func(name string) error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(false).SetMainBranch(name)
		return nil
	})

	suite.Step(`^the main branch is now configured as "([^"]+)"$`, func(name string) error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetMainBranch()
		if actual != name {
			return fmt.Errorf("expected %q, got %q", name, actual)
		}
		return nil
	})

	suite.Step(`^the main branch name is not configured$`, func() error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(false).DeleteMainBranchConfiguration()
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is set to (true|false)$`, func(text string) error {
		value, err := strconv.ParseBool(text)
		if err != nil {
			return fmt.Errorf("cannot parse %q into bool: %w", text, err)
		}
		_ = fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(false).SetNewBranchPush(value, false)
		return nil
	})

	suite.Step(`^the perennial branches are not configured$`, func() error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(false).DeletePerennialBranchConfiguration()
		return nil
	})

	suite.Step(`^the perennial branches are configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(false).AddToPerennialBranches(branch1, branch2)
		return nil
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)"$`, func(name string) error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetPerennialBranches()
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if actual[0] != name {
			return fmt.Errorf("expected %q, got %q", name, actual[0])
		}
		return nil
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetPerennialBranches()
		if len(actual) != 2 {
			return fmt.Errorf("expected 2 perennial branches, got %q", actual)
		}
		if actual[0] != branch1 || actual[1] != branch2 {
			return fmt.Errorf("expected %q, got %q", []string{branch1, branch2}, actual)
		}
		return nil
	})

	suite.Step(`^my repo is now configured with no perennial branches$`, func() error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).GetPerennialBranches()
		if len(actual) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", actual)
		}
		return nil
	})
}
