Feature: make the current observed branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    And the current branch is "observed"
    When I run "git-town contribute"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "observed" is now a contribution branch
      """
    And branch "observed" now has type "contribution"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "observed" now has type "observed"
