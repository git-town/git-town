Feature: make the current observed branch an observed branch

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE     | LOCATIONS |
      | observed | observed | local     |
    And the current branch is "observed"
    When I run "git-town observe observed"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      branch "observed" is already observed
      """
    And branch "observed" still has type "observed"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "observed" still has type "observed"
