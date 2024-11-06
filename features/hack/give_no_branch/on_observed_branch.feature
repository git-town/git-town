Feature: making the current observed branch a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT | LOCATIONS |
      | observed | observed |        | local     |
    And the current branch is "observed"
    When I run "git-town hack"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "observed" is now a feature branch
      """
    And branch "observed" is now a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "observed" is now observed
