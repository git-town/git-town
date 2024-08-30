Feature: observe the current observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS |
      | observed | observed | local     |
    And the current branch is "observed"
    When I run "git-town observe"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "observed" is already observed
      """
    And the observed branches are still "observed"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the observed branches are still "observed"
