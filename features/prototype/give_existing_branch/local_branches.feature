@this
Feature: make multiple branches prototype

  Background:
    Given a Git repo clone
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | feature      | feature      | main   | local, origin |
      | contribution | contribution |        | local, origin |
      | observed     | observed     |        | local, origin |
      | parked       | parked       | main   | local, origin |
    And an uncommitted file
    When I run "git-town prototype feature contribution observed parked"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature" is now a prototype branch
      """
    And branch "feature" is now prototype
    And it prints:
      """
      branch "contribution" is now a prototype branch
      """
    And branch "contribution" is now prototype
    And it prints:
      """
      branch "observed" is now a prototype branch
      """
    And branch "observed" is now prototype
    And it prints:
      """
      branch "parked" is now a prototype branch
      """
    And branch "parked" is now prototype
    And branch "parked" is still
    And the current branch is still "main"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And there are now no prototype branches
    And the current branch is still "main"
    And the uncommitted file still exists
    And the initial branches and lineage exist
