Feature: sync the current observed branch

  Background:
    Given the current branch is an observed branch "other"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   |
      | other  | local         | local commit  | local_file  |
      |        | origin        | origin commit | origin_file |
      | main   | local, origin | main commit   | main_file   |
    And the current branch is "other"
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
      |        | git rebase origin/other  |
      |        | git push --tags          |
    And the current branch is still "other"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | other  | local, origin | origin commit |
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
