Feature: make an existing contribution branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution | main   | local, origin |
    And the current branch is "contribution"
    And an uncommitted file
    When I run "git-town contribute contribution"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "contribution" is already a contribution branch
      """
    And the contribution branches are still "contribution"
    And the current branch is still "contribution"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "contribution"
    And the contribution branches are still "contribution"
    And the uncommitted file still exists
