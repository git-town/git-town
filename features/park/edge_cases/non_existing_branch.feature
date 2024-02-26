Feature: trying to park a non-existing branch

  Background:
    Given an uncommitted file
    When I run "git-town park non-existing"

  Scenario: result
    Then it runs no commands
    And the current branch is still "main"
    And the main branch is still "main"
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
    And the current branch is still "main"
    And the main branch is still "main"
