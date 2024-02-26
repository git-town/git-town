Feature: park a remote branch by local name

  Background:
    Given a known remote feature branch "feature"
    And an uncommitted file
    When I run "git-town park feature"

  Scenario: result
    Then it runs no commands
    And branch "feature" is now parked
    And the current branch is still "main"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And there are now no parked branches
    And the current branch is still "main"
    And the uncommitted file still exists
