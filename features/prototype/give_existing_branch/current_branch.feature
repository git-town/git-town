Feature: prototype the current branch

  Background:
    Given a local Git repo
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And an uncommitted file
    And the current branch is "feature"
    When I run "git-town prototype feature"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature" is now a prototype branch
      """
    And the prototype branches are now "feature"
    And the current branch is still "feature"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND       |
      | feature | git add -A    |
      |         | git stash     |
      |         | git stash pop |
    And there are now no prototype branches
    And the current branch is still "feature"
    And the uncommitted file still exists
    And the initial branches and lineage exist
