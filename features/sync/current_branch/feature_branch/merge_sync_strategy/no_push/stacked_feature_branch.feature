Feature: syncing a stacked feature branch using --no-push

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | (none)  |        | local, origin |
      | parent  | feature | main   | local, origin |
      | child   | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
    And the current branch is "child"
    When I run "git-town sync --no-push"

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
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                  |
      | main   | local, origin | origin main commit                                       |
      |        | local         | local main commit                                        |
      | child  | local         | local child commit                                       |
      |        |               | Merge branch 'parent' into child                         |
      |        | local, origin | origin child commit                                      |
      |        | local         | Merge remote-tracking branch 'origin/child' into child   |
      | parent | local         | local parent commit                                      |
      |        |               | Merge branch 'main' into parent                          |
      |        | local, origin | origin parent commit                                     |
      |        | local         | Merge remote-tracking branch 'origin/parent' into parent |
    And the initial branches and lineage exist now

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
    And the initial commits exist now
    And the initial branches and lineage exist now
