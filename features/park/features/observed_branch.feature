Feature: parking an observed branch

  Background:
    Given the current branch is an observed branch "observed"
    And an uncommitted file
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "observed" is now parked
      """
    And the current branch is still "observed"
    And branch "observed" is now parked
    And there are now no observed branches
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND       |
      | observed | git add -A    |
      |          | git stash     |
      |          | git stash pop |
    And the current branch is still "observed"
    And branch "observed" is now observed
    And there are now no parked branches
    And the uncommitted file still exists
