package steps

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/cucumber/godog"
)

// ConfigurationSteps defines Cucumber step implementations around configuration.
// nolint:funlen
func ConfigurationSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^Git Town is no longer configured for this repository$`, func() error {
		outcome, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Run("git", "config", "--local", "--get-regex", "git-town")
		exitError := err.(*exec.ExitError)
		if exitError.ExitCode() != 1 {
			return fmt.Errorf("git config should return exit code 1 if no matching configuration found")
		}
		if outcome.OutputSanitized() != "" {
			return fmt.Errorf("expected no local Git Town configuration but got %q: %w", outcome.Output(), err)
		}
		return nil
	})

	suite.Step(`^I haven't configured Git Town yet$`, func() error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).DeleteMainBranchConfiguration()
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).DeletePerennialBranchConfiguration()
		return nil
	})

	suite.Step(`^my repo is now configured with no perennial branches$`, func() error {
		branches := fs.activeScenarioState.gitEnvironment.DeveloperRepo.FreshConfiguration().GetPerennialBranches()
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is set to "(true|false)"$`, func(value string) error {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot parse %q into bool: %w", value, err)
		}
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).SetNewBranchPush(b, false)
		return nil
	})

	suite.Step(`^the main branch is configured as "([^"]+)"$`, func(name string) error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).SetMainBranch(name)
		return nil
	})

	suite.Step(`^the main branch is now configured as "([^"]+)"$`, func(name string) error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.FreshConfiguration().GetMainBranch()
		if actual != name {
			return fmt.Errorf("expected %q, got %q", name, actual)
		}
		return nil
	})

	suite.Step(`^the main branch name is not configured$`, func() error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).DeleteMainBranchConfiguration()
		return nil
	})

	suite.Step(`^the perennial branches are not configured$`, func() error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).DeletePerennialBranchConfiguration()
		return nil
	})

	suite.Step(`^the perennial branches are configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration(true).AddToPerennialBranches(branch1, branch2)
		return nil
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)"$`, func(name string) error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.FreshConfiguration().GetPerennialBranches()
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if actual[0] != name {
			return fmt.Errorf("expected %q, got %q", name, actual[0])
		}
		return nil
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.FreshConfiguration().GetPerennialBranches()
		if len(actual) != 2 {
			return fmt.Errorf("expected 2 perennial branches, got %q", actual)
		}
		if actual[0] != branch1 || actual[1] != branch2 {
			return fmt.Errorf("expected %q, got %q", []string{branch1, branch2}, actual)
		}
		return nil
	})

	suite.Step(`^my repo is now configured with no perennial branches$`, func() error {
		actual := fs.activeScenarioState.gitEnvironment.DeveloperRepo.FreshConfiguration().GetPerennialBranches()
		if len(actual) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", actual)
		}
		return nil
	})
}
