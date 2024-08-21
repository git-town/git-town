Feature: make the current observed branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    And the current branch is "observed"
    When I run "git-town contribute"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "observed" is now a contribution branch
      """
    And branch "observed" is now a contribution branch
    And the current branch is still "observed"
    And there are now no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "observed"
    And branch "observed" is now observed
    And there are now no contribution branches
