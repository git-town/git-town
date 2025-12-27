Feature: making multiple branches feature branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution |        | local, origin |
      | observed     | observed     |        | local, origin |
      | parked       | parked       | main   | local         |
    When I run "git-town feature contribution observed parked"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch contribution is now a feature branch
      """
    And Git Town prints:
      """
      branch observed is now a feature branch
      """
    And Git Town prints:
      """
      branch parked is now a feature branch
      """
    And branch "contribution" now has type "feature"
    And branch "observed" now has type "feature"
    And branch "parked" now has type "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "contribution" now has type "contribution"
    And branch "observed" now has type "observed"
    And branch "parked" now has type "parked"
