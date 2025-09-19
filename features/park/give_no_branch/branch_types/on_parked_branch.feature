Feature: parking the current parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS |
      | parked | parked | main   | local     |
    And the current branch is "parked"
    When I run "git-town park"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "parked" is already parked
      """
    And branch "parked" still has type "parked"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "parked" still has type "parked"
