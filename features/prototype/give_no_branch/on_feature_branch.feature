Feature: make the current feature branch a prototype branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
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
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | branch | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "branch"
    And there are now no prototype branches
    And the initial branches exist
    And the uncommitted file still exists
