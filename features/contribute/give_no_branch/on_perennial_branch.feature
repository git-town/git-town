Feature: cannot make the current perennial branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS |
      | perennial | perennial | local     |
    And the current branch is "perennial"
    When I run "git-town contribute"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot make perennial branches contribution branches
      """
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no contribution branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no contribution branches
