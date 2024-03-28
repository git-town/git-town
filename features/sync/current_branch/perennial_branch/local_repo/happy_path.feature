Feature: sync the current perennial branch in a local repo

  Background:
    Given my repo does not have an origin
    And the local perennial branches "production" and "qa"
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | main commit  | main_file  |
      | qa     | local    | local commit | local_file |
    And the current branch is "qa"
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND       |
      | qa     | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And all branches are now synchronized
    And the current branch is still "qa"
    And the uncommitted file still exists
    And the initial commits exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | qa     | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "qa"
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | main   | local    | main commit  |
      | qa     | local    | local commit |
    And the initial branches and lineage exist
