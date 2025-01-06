Feature: sync the main branch

  Background:
    Given a Git repo with origin
    And the commits
      | LOCATION | MESSAGE       | FILE NAME   |
      | local    | local commit  | local_file  |
      | origin   | origin commit | origin_file |
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | main   | git fetch --prune --tags                |
      |        | git rebase origin/main --no-update-refs |
      |        | git push                                |
      |        | git push --tags                         |
    And the current branch is still "main"
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | origin commit |
      |        |               | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And the current branch is still "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | origin commit |
      |        |               | local commit  |
    And the initial branches and lineage exist now
