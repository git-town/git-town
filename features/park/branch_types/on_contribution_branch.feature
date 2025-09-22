Feature: parking the current contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS     |
      | contribution | contribution | local, origin |
    And the current branch is "contribution"
    When I run "git-town park"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "contribution" is now parked
      """
    And branch "contribution" now has type "parked"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "contribution" now has type "contribution"
