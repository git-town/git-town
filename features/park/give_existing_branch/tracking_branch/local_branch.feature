Feature: park another local feature branch

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    When I run "git-town park feature"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "feature" is now parked
      """
    And the parked branches are now "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "feature" now has type "feature"
    And there are now no parked branches
