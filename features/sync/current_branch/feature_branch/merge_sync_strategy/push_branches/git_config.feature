Feature: disable pushing through Git metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    And Git setting "git-town.push-branches" is "false"
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | child  | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout parent                               |
      | parent | git merge --no-edit --ff main                     |
      |        | git merge --no-edit --ff origin/parent            |
      |        | git checkout child                                |
      | child  | git merge --no-edit --ff parent                   |
      |        | git merge --no-edit --ff origin/child             |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                  |
      | main   | local, origin | origin main commit                                       |
      |        | local         | local main commit                                        |
      | parent | local         | local parent commit                                      |
      |        |               | Merge branch 'main' into parent                          |
      |        | local, origin | origin parent commit                                     |
      |        | local         | Merge remote-tracking branch 'origin/parent' into parent |
      | child  | local         | local child commit                                       |
      |        |               | Merge branch 'parent' into child                         |
      |        | local, origin | origin child commit                                      |
      |        | local         | Merge remote-tracking branch 'origin/child' into child   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                          |
      | child  | git reset --hard {{ sha 'local child commit' }}  |
      |        | git checkout main                                |
      | main   | git reset --hard {{ sha 'local main commit' }}   |
      |        | git checkout parent                              |
      | parent | git reset --hard {{ sha 'local parent commit' }} |
      |        | git checkout child                               |
    And the initial branches and lineage exist now
    And the initial commits exist now
