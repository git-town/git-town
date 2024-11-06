Feature: observe the current perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS |
      | perennial | perennial | local     |
    And the current branch is "perennial"
    And an uncommitted file
    When I run "git-town observe"

  Scenario: result
    Then Git Town runs no commands
    And it prints the error:
      """
      cannot observe perennial branches
      """
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no observed branches
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no observed branches
    And the uncommitted file still exists
