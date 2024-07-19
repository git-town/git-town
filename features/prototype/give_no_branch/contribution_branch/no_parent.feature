Feature: make a contribution branch without parent a prototype

  Background:
    Given a Git repo clone
    And the branch
      | NAME   | TYPE         | LOCATIONS |
      | branch | contribution | origin    |
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
    And there are now no contribution branches
    And the uncommitted file still exists

  @this
  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | branch | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "branch"
    And the initial branches and lineage exist
    And there are now no prototype branches
    And the uncommitted file still exists
