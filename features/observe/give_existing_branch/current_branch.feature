Feature: make the current observed branch an observed branch

  Background:
    Given a local Git repo
    And the branch
      | NAME     | TYPE     | LOCATIONS |
      | observed | observed | local     |
    And the current branch is "observed"
    When I run "git-town observe observed"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "observed" is already observed
      """
    And the observed branches are still "observed"
    And the current branch is still "observed"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "observed"
    And the observed branches are still "observed"
