Feature: parking the main branch

  Background:
    Given an uncommitted file
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And the current branch is still "main"
    And the main branch is still "main"
    And the uncommitted file still exists
    And it prints the error:
      """
      cannot park the main branch
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And the main branch is still "main"
    And there are now no parked branches
