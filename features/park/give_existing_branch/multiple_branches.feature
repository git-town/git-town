Feature: parking multiple other branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | feature      | feature      | main   | local         |
      | contribution | contribution |        | local, origin |
      | observed     | observed     | main   | local, origin |
      | prototype    | prototype    | main   | local         |
    When I run "git-town park feature contribution observed prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "feature" is now parked
      """
    And branch "feature" now has type "parked"
    And branch "contribution" now has type "parked"
    And branch "observed" now has type "parked"
    And branch "prototype" now has type "parked"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "feature" now has type "feature"
    And branch "contribution" now has type "contribution"
    And branch "observed" now has type "observed"
    And branch "prototype" now has type "prototype"
