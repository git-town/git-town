Feature: make multiple branches prototype

  Background:
    Given a Git repo clone
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
      | feature-2 | feature | main   | local, origin |
      | feature-3 | feature | main   | local, origin |
    And an uncommitted file
    When I run "git-town prototype feature-1 feature-2 feature-3"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature-1" is now a prototype branch
      """
    And branch "feature-1" is now prototype
    And it prints:
      """
      branch "feature-2" is now a prototype branch
      """
    And branch "feature-2" is now prototype
    And it prints:
      """
      branch "feature-3" is now a prototype branch
      """
    And branch "feature-3" is now prototype
    And the current branch is still "main"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And there are now no observed branches
    And the current branch is still "main"
    And the uncommitted file still exists
