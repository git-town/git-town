Feature: prototype the current observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS |
      | observed | observed | local     |
    And the current branch is "observed"
    When I run "git-town prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "observed" is now a prototype branch
      """
    And branch "observed" is now prototype
    And there are now no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "observed"
    And branch "observed" is now observed
    And there are now no parked branches
