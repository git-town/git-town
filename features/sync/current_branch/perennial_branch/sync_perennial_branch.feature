Feature: sync the current perennial branch

  Background:
    Given my repo has the perennial branches "production" and "qa"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   |
      | qa     | local         | local commit  | local_file  |
      |        | origin        | origin commit | origin_file |
      | main   | local, origin | main commit   | main_file   |
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
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | qa     | local, origin | origin commit |
      |        |               | local commit  |
