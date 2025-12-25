Feature: observe the current observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    And the current branch is "observed"
    When I run "git-town observe"

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
