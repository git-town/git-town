Feature: remove a parked branch as soon as the tracking branch is gone, even if it has unpushed commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  |
      | main   | local, origin | main commit  | main_file  |
      | parked | local         | local commit | local_file |
    And origin deletes the "parked" branch
    And the current branch is "parked"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | parked | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git branch -D parked     |
    And Git Town prints:
      """
      deleted branch "parked"
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                            |
      | main   | git branch parked {{ sha-initial 'local commit' }} |
      |        | git checkout parked                                |
    And the initial branches and lineage exist now
    And the initial commits exist now
