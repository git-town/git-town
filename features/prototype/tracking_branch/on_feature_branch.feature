Feature: prototype the current feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch feature is now a prototype branch
      """
    And branch "feature" now has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "feature" now has type "feature"
    And the initial branches exist now
