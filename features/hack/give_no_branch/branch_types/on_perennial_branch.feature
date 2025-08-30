Feature: cannot make the current perennial branch a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE   | LOCATIONS |
      | perennial | (none) | local     |
    And the current branch is "perennial"
    And local Git setting "git-town.perennial-branches" is "perennial"
    When I run "git-town hack"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot make perennial branches feature branches
      """
    And branch "perennial" still has type "perennial"
    And local Git setting "git-town-branch.perennial.branchtype" still doesn't exist
    And local Git setting "git-town.perennial-branches" is still "perennial"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "perennial" still has type "perennial"
    And local Git setting "git-town-branch.perennial.branchtype" still doesn't exist
