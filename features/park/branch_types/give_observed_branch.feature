Feature: park another observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    When I run "git-town park observed"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch observed is now parked
      """
    And branch "observed" now has type "parked"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "observed" now has type "observed"
