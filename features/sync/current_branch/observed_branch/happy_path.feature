Feature: sync the current observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE       | FILE NAME   |
      | main     | local, origin | main commit   | main_file   |
      | observed | local         | local commit  | local_file  |
      |          | origin        | origin commit | origin_file |
    And the current branch is "observed"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                               |
      | observed | git fetch --prune --tags                              |
      |          | git -c rebase.updateRefs=false rebase origin/observed |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE       |
      | main     | local, origin | main commit   |
      | observed | local, origin | origin commit |
      |          | local         | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | observed | git reset --hard {{ sha-initial 'local commit' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
