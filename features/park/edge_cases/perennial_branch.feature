Feature: parking a perennial branch

  Background:
    Given the current branch is a perennial branch "perennial"
    And an uncommitted file
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And the uncommitted file still exists
    And it prints the error:
      """
      cannot park perennial branches
      """
    And there are still no parked branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "perennial"
    And the uncommitted file still exists
    And the perennial branches are still "perennial"
    And there are still no parked branches
