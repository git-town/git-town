Feature: syncing a stacked feature branch using --no-push

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    When I run "git-town sync --no-push"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                |
      | child  | git fetch --prune --tags               |
      |        | git checkout main                      |
      | main   | git rebase origin/main                 |
      |        | git checkout parent                    |
      | parent | git merge --no-edit --ff origin/parent |
      |        | git merge --no-edit --ff main          |
      |        | git checkout child                     |
      | child  | git merge --no-edit --ff origin/child  |
      |        | git merge --no-edit --ff parent        |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                  |
      | main   | local, origin | origin main commit                                       |
      |        | local         | local main commit                                        |
      | child  | local         | local child commit                                       |
      |        | local, origin | origin child commit                                      |
      |        | local         | Merge remote-tracking branch 'origin/child' into child   |
      |        |               | local parent commit                                      |
      |        |               | origin parent commit                                     |
      |        |               | Merge remote-tracking branch 'origin/parent' into parent |
      |        |               | origin main commit                                       |
      |        |               | local main commit                                        |
      |        |               | Merge branch 'main' into parent                          |
      |        |               | Merge branch 'parent' into child                         |
      | parent | local         | local parent commit                                      |
      |        | local, origin | origin parent commit                                     |
      |        | local         | Merge remote-tracking branch 'origin/parent' into parent |
      |        |               | origin main commit                                       |
      |        |               | local main commit                                        |
      |        |               | Merge branch 'main' into parent                          |
    And the initial branches and lineage exist

  @this
  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                          |
      | child  | git reset --hard {{ sha 'local child commit' }}  |
      |        | git checkout main                                |
      | main   | git reset --hard {{ sha 'local main commit' }}   |
      |        | git checkout parent                              |
      | parent | git reset --hard {{ sha 'local parent commit' }} |
      |        | git checkout child                               |
    And the current branch is still "child"
    And the initial commits exist
    And the initial branches and lineage exist
