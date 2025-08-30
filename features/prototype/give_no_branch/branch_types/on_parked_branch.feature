Feature: prototype the current parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS |
      | parked | parked | main   | local     |
    And the current branch is "parked"
    When I run "git-town prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "parked" is now a prototype branch
      """
    And branch "parked" now has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "parked" now has type "parked"
