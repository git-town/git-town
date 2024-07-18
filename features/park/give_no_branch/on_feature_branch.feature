Feature: parking the current branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the current branch is "branch"
    And an uncommitted file
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "branch" is now parked
      """
    And the current branch is still "branch"
    And branch "branch" is now parked
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | branch | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "branch"
    And there are now no parked branches
    And the uncommitted file still exists
