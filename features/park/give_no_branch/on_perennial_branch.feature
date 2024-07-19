Feature: cannot park perennial branches

  Background:
    Given a Git repo clone
    And the branch
      | NAME      | TYPE      | LOCATIONS |
      | perennial | perennial | local     |
    And the current branch is "perennial"
    And an uncommitted file
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot park perennial branches
      """
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And the uncommitted file still exists
    And there are still no parked branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "perennial"
    And the uncommitted file still exists
    And the perennial branches are still "perennial"
    And there are still no parked branches
