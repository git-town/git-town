Feature: make another perennial branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME      | TYPE      | LOCATIONS     |
      | perennial | perennial | local, origin |
    And an uncommitted file
    When I run "git-town contribute perennial"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot make perennial branches contribution branches
      """
    And the perennial branches are still "perennial"
    And there are now no contribution branches
    And the current branch is still "main"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And the perennial branches are still "perennial"
    And the uncommitted file still exists
