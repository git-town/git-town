Feature: making multiple branches a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS |
      | contribution | contribution |        | local     |
      | observed     | observed     |        | local     |
      | parked       | parked       | main   | local     |
    When I run "git-town hack contribution observed parked"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "contribution" is now a feature branch
      """
    And branch "contribution" now has type "feature"
    And Git Town prints:
      """
      branch "observed" is now a feature branch
      """
    And branch "observed" now has type "feature"
    And Git Town prints:
      """
      branch "parked" is now a feature branch
      """
    And branch "parked" now has type "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "contribution" now has type "contribution"
    And branch "observed" now has type "observed"
    And branch "parked" now has type "parked"
