Feature: prototype another local feature branch

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    When I run "git-town prototype feature"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "feature" is now a prototype branch
      """
    And the prototype branches are now "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "feature" is now a feature branch
    And there are now no prototype branches
