Feature: make another observed branch an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | observed | observed |        | local, origin |
    When I run "git-town observe observed"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "observed" is already observed
      """
    And branch "observed" still has type "observed"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "observed" still has type "observed"
