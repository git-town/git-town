Feature: observe a local feature branch

  Background:
    Given the current branch is a feature branch "branch"
    And an uncommitted file
    When I run "git-town observe"

  @this
  Scenario: result
    Then it runs no commands
    And the current branch is still "branch"
    And the uncommitted file still exists
    And branch "branch" is now observed

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | branch | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "branch"
    And the uncommitted file still exists
    And there are now no observed branches
