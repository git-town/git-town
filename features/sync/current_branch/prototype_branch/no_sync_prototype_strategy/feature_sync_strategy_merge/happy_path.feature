Feature: sync the current prototype branch with tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION | MESSAGE            |
      | main      | local    | main local commit  |
      |           | local    | main origin commit |
      | prototype | local    | local commit       |
      |           | origin   | origin commit      |
    And the current branch is "prototype"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                           |
      | prototype | git fetch --prune --tags                          |
      |           | git checkout main                                 |
      | main      | git -c rebase.updateRefs=false rebase origin/main |
      |           | git push                                          |
      |           | git checkout prototype                            |
      | prototype | git merge --no-edit --ff main                     |
      |           | git merge --no-edit --ff origin/prototype         |
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE                                                        |
      | main      | local, origin | main local commit                                              |
      |           |               | main origin commit                                             |
      | prototype | local         | local commit                                                   |
      |           |               | Merge branch 'main' into prototype                             |
      |           | local, origin | origin commit                                                  |
      |           | local         | Merge remote-tracking branch 'origin/prototype' into prototype |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                           |
      | prototype | git reset --hard {{ sha-initial 'local commit' }} |
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE            |
      | main      | local, origin | main local commit  |
      |           |               | main origin commit |
      | prototype | local         | local commit       |
      |           | origin        | origin commit      |
    And the initial branches and lineage exist now
