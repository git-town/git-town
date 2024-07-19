Feature: cannot make perennial branches contribution branches

  Background:
    Given a Git repo clone
    And the branch
      | NAME      | TYPE      | LOCATIONS |
      | perennial | perennial | local     |
    And the current branch is "perennial"
    And an uncommitted file
    When I run "git-town contribute"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot make perennial branches contribution branches
      """
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no contribution branches
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no contribution branches
    And the uncommitted file still exists
