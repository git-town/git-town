Feature: parking an observed branch

  Background:
    Given the current branch is an observed branch "observed"
    And an uncommitted file
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And the current branch is still "observed"
    And branch "observed" is now parked
    And the uncommitted file still exists
    And there are now no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND       |
      | observed | git add -A    |
      |          | git stash     |
      |          | git stash pop |
    And the current branch is still "observed"
    And the uncommitted file still exists
    And branch "observed" is now an observed branch
    And there are now no parked branches
