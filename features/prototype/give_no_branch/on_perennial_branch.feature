Feature: cannot prototype perennial branches

  Background:
    Given a Git repo clone
    And the branch
      | NAME   | TYPE      | LOCATIONS |
      | branch | perennial | local     |
    Given the current branch is "branch"
    And an uncommitted file
    When I run "git-town prototype"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot prototype perennial branches
      """
    And the current branch is still "branch"
    And the perennial branches are still "branch"
    And there are still no prototype branches
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "branch"
    And the perennial branches are still "branch"
    And there are still no prototype branches
    And the uncommitted file still exists
