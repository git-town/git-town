Feature: sync the current prototype branch that has a tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION | MESSAGE       | FILE NAME   |
      | main      | local    | main commit   | main_file   |
      | prototype | local    | local commit  | local_file  |
      |           | origin   | origin commit | origin_file |
    And Git setting "git-town.sync-prototype-strategy" is "rebase"
    And the current branch is "prototype"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                      |
      | prototype | git fetch --prune --tags                                                     |
      |           | git checkout main                                                            |
      | main      | git -c rebase.updateRefs=false rebase origin/main                            |
      |           | git push                                                                     |
      |           | git checkout prototype                                                       |
      | prototype | git push --force-with-lease --force-if-includes                              |
      |           | git -c rebase.updateRefs=false rebase origin/prototype                       |
      |           | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |           | git push --force-with-lease --force-if-includes                              |
    And no rebase is now in progress
    And the initial branches and lineage exist now
    And all branches are now synchronized
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE       |
      | main      | local, origin | main commit   |
      | prototype | local, origin | origin commit |
      |           |               | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                                  |
      | prototype | git reset --hard {{ sha 'local commit' }}                                                |
      |           | git push --force-with-lease origin {{ sha-in-origin-initial 'origin commit' }}:prototype |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE       | FILE NAME   |
      | main      | local, origin | main commit   | main_file   |
      | prototype | local         | local commit  | local_file  |
      |           | origin        | origin commit | origin_file |
