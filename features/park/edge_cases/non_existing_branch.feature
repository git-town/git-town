Feature: trying to park a non-existing branch

  Background:
    Given the current branch is a feature branch "branch"
    And an uncommitted file
    When I run "git-town park branch non-existing"

  Scenario: result
    Then it runs no commands
    And the current branch is still "branch"
    And the uncommitted file still exists
    And it prints the error:
      """
      there is no branch "non-existing"
      """
    And there are still no parked branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And there are still no parked branches
    And the current branch is still "branch"
