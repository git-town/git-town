Feature: make another observed branch an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | observed | observed |        | local, origin |
    When I run "git-town observe observed"

  Scenario: result
    Then Git Town runs no commands
    And it prints the error:
      """
      branch "observed" is already observed
      """
    And the observed branches are still "observed"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the observed branches are still "observed"
