Feature: remove a contribution branch as soon as its tracking branch is gone, even if it has unpushed commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS     |
      | contribution | contribution | local, origin |
    And the commits
      | BRANCH       | LOCATION      | MESSAGE      | FILE NAME  |
      | main         | local, origin | main commit  | main_file  |
      | contribution | local         | local commit | local_file |
    And origin deletes the "contribution" branch
    And the current branch is "contribution"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                    |
      | contribution | git fetch --prune --tags   |
      |              | git checkout main          |
      | main         | git branch -D contribution |
    And Git Town prints:
      """
      deleted branch "contribution"
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                  |
      | main   | git branch contribution {{ sha-initial 'local commit' }} |
      |        | git checkout contribution                                |
    And the initial branches and lineage exist now
    And the initial commits exist now
