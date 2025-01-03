Feature: making the current observed branch a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE   | LOCATIONS |
      | observed | (none) | local     |
    And the current branch is "observed"
    And local Git setting "git-town.observed-branches" is "observed"
    When I run "git-town hack"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "observed" is now a feature branch
      """
    And branch "observed" now has type "feature"
    And local Git setting "git-town-branch.observed.branchtype" is now "feature"
    And local Git setting "git-town.observed-branches" is still "observed"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "observed" now has type "observed"
    And local Git setting "git-town-branch.observed.branchtype" now doesn't exist
