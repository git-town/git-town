Feature: prototype the current local branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town prototype"

  Scenario: result
    Then Git Town runs no commands
    And it prints:
      """
      branch "feature" is now a prototype branch
      """
    And the prototype branches are now "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "feature"
    And there are now no parked branches
