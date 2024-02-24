Feature: parking an observed branch

  Background:
    Given the current branch is a observed branch "observed"
    And an uncommitted file
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And the current branch is still "observed"
    And branch "observed" is still observed
    And it prints the error:
      """
      Cannot park observed branches
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "observed"
    And branch "observed" is still observed
