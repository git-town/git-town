Feature: cannot observe the main branch

  Background:
    Given a Git repo clone
    And an uncommitted file
    When I run "git-town observe"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot observe the main branch
      """
    And the current branch is still "main"
    And the main branch is still "main"
    And the uncommitted file still exists
    And there are still no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And the main branch is still "main"
    And there are still no observed branches
