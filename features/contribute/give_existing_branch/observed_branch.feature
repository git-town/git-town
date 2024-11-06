Feature: make another observed branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | observed | observed |        | local, origin |
    When I run "git-town contribute observed"

  Scenario: result
    Then Git Town runs no commands
    And it prints:
      """
      branch "observed" is now a contribution branch
      """
    And the contribution branches are now "observed"
    And there are now no observed branches
    And the current branch is still "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "main"
    And the observed branches are now "observed"
