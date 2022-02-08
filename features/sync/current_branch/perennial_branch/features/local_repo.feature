Feature: sync the current perennial branch (without remote repo)

  Background:
    Given my repo does not have an origin remote
    And my repo has the local perennial branches "production" and "qa"
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | main commit  | main_file  |
      | qa     | local    | local commit | local_file |
    And I am on the "qa" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND       |
      | qa     | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And all branches are now synchronized
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And now the initial commits exist
