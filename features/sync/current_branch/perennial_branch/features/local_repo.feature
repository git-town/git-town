Feature: sync the current perennial branch (local repo)

  Background:
    Given my repo does not have an origin
    And the local perennial branches "production" and "qa"
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | main commit  | main_file  |
      | qa     | local    | local commit | local_file |
    And the current branch is "qa"
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND       |
      | qa     | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And all branches are now synchronized
    And the current branch is still "qa"
    And my workspace still contains my uncommitted file
    And now the initial commits exist
