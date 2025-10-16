Feature: sync the current perennial branch using the rebase sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | LOCATIONS     |
      | production | perennial | local, origin |
      | qa         | perennial | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   |
      | qa     | local         | local commit  | local_file  |
      |        | origin        | origin commit | origin_file |
      | main   | local, origin | main commit   | main_file   |
    And Git setting "git-town.sync-perennial-strategy" is "rebase"
    And the current branch is "qa"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | qa     | git fetch --prune --tags                        |
      |        | git -c rebase.updateRefs=false rebase origin/qa |
      |        | git push                                        |
      |        | git push --tags                                 |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | qa     | local, origin | origin commit |
      |        |               | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | qa     | local, origin | origin commit |
      |        |               | local commit  |
