Feature: parking the current contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS |
      | contribution | contribution | local     |
    And the current branch is "contribution"
    When I run "git-town park"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "contribution" now has type "parked"
      """
    And the current branch is still "contribution"
    And branch "contribution" now has type "parked"
    And there are now no contribution branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "contribution"
    And branch "contribution" now has type "contribution"
    And there are now no parked branches
