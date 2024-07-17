Feature: observing a parked branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS |
      | branch | parked | main   | local     |
    Given the current branch is "branch"
    And an uncommitted file
    When I run "git-town observe"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "branch" is now an observed branch
      """
    And the current branch is still "branch"
    And branch "branch" is now observed
    And there are now no parked branches
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | branch | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "branch"
    And branch "branch" is now parked
    And there are now no observed branches
    And the uncommitted file still exists
