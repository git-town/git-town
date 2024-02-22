Feature: sync the current observed branch

  Background:
    Given the current branch is an observed branch "other"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   |
      | main   | local, origin | main commit   | main_file   |
      | other  | local         | local commit  | local_file  |
      |        | origin        | origin commit | origin_file |
    When I run "git-town sync"

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
      |        | local         | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                              |
      | other  | git reset --hard {{ sha-before-run 'local commit' }} |
    And the current branch is still "other"
    And the initial commits exist
    And the initial branches and lineage exist
