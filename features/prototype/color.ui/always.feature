Feature: prototype multiple other branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | feature      | feature      | main   | local, origin |
      | contribution | contribution |        | local, origin |
      | observed     | observed     |        | local, origin |
      | parked       | parked       | main   | local, origin |
    And Git setting "color.ui" is "always"
    When I run "git-town prototype feature contribution observed parked"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch feature is now a prototype branch
      """
    And Git Town prints:
      """
      branch contribution is now a prototype branch
      """
    And Git Town prints:
      """
      branch observed is now a prototype branch
      """
    And Git Town prints:
      """
      branch parked is now a prototype branch
      """
    And branch "contribution" now has type "prototype"
    And branch "feature" now has type "prototype"
    And branch "observed" now has type "prototype"
    And branch "parked" now has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And branch "contribution" now has type "contribution"
    And branch "feature" now has type "feature"
    And branch "observed" now has type "observed"
    And branch "parked" now has type "parked"
