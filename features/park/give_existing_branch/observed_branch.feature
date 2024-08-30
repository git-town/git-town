Feature: park another observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    When I run "git-town park observed"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "observed" is now parked
      """
    And the parked branches are now "observed"
    And there are now no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the observed branches are now "observed"
    And there are now no parked branches
