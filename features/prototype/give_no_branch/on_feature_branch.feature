Feature: prototype the current feature branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town prototype"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature" is now a prototype branch
      """
    And the current branch is still "feature"
    And branch "feature" is now prototype

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And there are now no prototype branches
    And the initial branches exist
