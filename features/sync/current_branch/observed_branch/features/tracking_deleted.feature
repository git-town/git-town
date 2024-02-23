Feature: remove the observed branch as soon as the tracking branch is gone, even if it has unpushed commits

  Background:
    Given the current branch is an observed branch "other"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  |
      | main   | local, origin | main commit  | main_file  |
      | other  | local         | local commit | local_file |
    And origin deletes the "other" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git branch -D other      |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                              |
      | main   | git branch other {{ sha-before-run 'local commit' }} |
      |        | git checkout other                                   |
    And the current branch is now "other"
    And the initial commits exist
    And the initial branches and lineage exist
