Feature: cannot prototype perennial branches

  Background:
    Given a Git repo clone
    And the branch
      | NAME   | TYPE      | LOCATIONS |
      | branch | perennial | local     |
    Given the current branch is "branch"
    And an uncommitted file
    When I run "git-town prototype"

  @debug @this
  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot prototype perennial branches
      """
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no prototype branches
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "perennial"
    And the perennial branches are still "perennial"
    And there are still no prototype branches
    And the uncommitted file still exists
