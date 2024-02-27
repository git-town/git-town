Feature: cannot make perennial branches contribution branches

  Background:
    Given the current branch is a perennial branch "perennial"
    And an uncommitted file
    When I run "git-town contribute"

  Scenario: result
    Then it runs no commands
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no contribution branches
    And it prints the error:
      """
      cannot make perennial branches contribution branches
      """
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no contribution branches
    And the uncommitted file still exists
