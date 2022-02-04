Feature: sync the current perennial branch

  Background:
    Given my repo has the perennial branches "production" and "qa"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   |
      | qa     | local         | local commit  | local_file  |
      |        | remote        | remote commit | remote_file |
      | main   | local, remote | main commit   | main_file   |
    And I am on the "qa" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | qa     | git fetch --prune --tags |
      |        | git rebase origin/qa     |
      |        | git push                 |
      |        | git push --tags          |
    And all branches are now synchronized
    And I am still on the "qa" branch
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, remote | main commit   |
      | qa     | local, remote | remote commit |
      |        |               | local commit  |
