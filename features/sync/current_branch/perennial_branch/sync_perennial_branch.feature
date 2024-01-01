Feature: sync the current perennial branch

  Background:
    Given the perennial branches "production" and "qa"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   |
      | qa     | local         | local commit  | local_file  |
      |        | origin        | origin commit | origin_file |
      | main   | local, origin | main commit   | main_file   |
    And the current branch is "qa"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | qa     | git fetch --prune --tags |
      |        | git rebase origin/qa     |
      |        | git push                 |
      |        | git push --tags          |
    And all branches are now synchronized
    And the current branch is still "qa"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | qa     | local, origin | origin commit |
      |        |               | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "qa"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | qa     | local, origin | origin commit |
      |        |               | local commit  |
    And the initial branches and lineage exist
