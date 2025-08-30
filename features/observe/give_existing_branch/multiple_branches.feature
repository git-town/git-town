Feature: observing multiple other branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | feature      | feature      | main   | local, origin |
      | contribution | contribution |        | local, origin |
      | parked       | parked       | main   | local, origin |
      | prototype    | prototype    | main   | local, origin |
    When I run "git-town observe feature contribution parked prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "feature" is now an observed branch
      """
    And branch "feature" now has type "observed"
    And Git Town prints:
      """
      branch "contribution" is now an observed branch
      """
    And branch "contribution" now has type "observed"
    And Git Town prints:
      """
      branch "parked" is now an observed branch
      """
    And branch "parked" now has type "observed"
    And Git Town prints:
      """
      branch "prototype" is now an observed branch
      """
    And branch "prototype" now has type "observed"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "feature" now has type "feature"
    And branch "contribution" now has type "contribution"
    And branch "parked" now has type "parked"
    And branch "prototype" now has type "prototype"
    And the initial branches exist now
