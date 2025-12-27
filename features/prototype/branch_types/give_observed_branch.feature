Feature: prototype another observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    When I run "git-town prototype observed"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch observed is now a prototype branch
      """
    And branch "observed" now has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "observed" now has type "observed"
