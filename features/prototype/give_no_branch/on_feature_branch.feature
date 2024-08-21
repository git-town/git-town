Feature: prototype the current feature branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the current branch is "branch"
    When I run "git-town prototype"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "branch" is now a prototype branch
      """
    And the current branch is still "branch"
    And branch "branch" is now prototype

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
