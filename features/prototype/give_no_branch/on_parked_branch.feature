Feature: making a parked branch a prototype

  Background:
    Given a Git repo with origin
    And the branch
      | NAME   | TYPE   | PARENT | LOCATIONS |
      | branch | parked | main   | local     |
    And the current branch is "branch"
    And an uncommitted file
    When I run "git-town prototype"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "branch" is now a prototype branch
      """
    And the current branch is still "branch"
    And branch "branch" is now prototype
    And branch "branch" is still parked
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | branch | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "branch"
    And branch "branch" is still parked
    And there are now no prototype branches
    And the uncommitted file still exists
