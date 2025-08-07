Feature: syncing a stacked feature branch using --no-push

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
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
      | BRANCH | COMMAND                                             |
      | child  | git fetch --prune --tags                            |
      |        | git checkout main                                   |
      | main   | git -c rebase.updateRefs=false rebase origin/main   |
      |        | git checkout parent                                 |
      | parent | git -c rebase.updateRefs=false rebase origin/parent |
      |        | git -c rebase.updateRefs=false rebase main          |
      |        | git checkout child                                  |
      | child  | git -c rebase.updateRefs=false rebase origin/child  |
      |        | git -c rebase.updateRefs=false rebase parent        |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        | local         | local main commit    |
      | child  | local         | origin child commit  |
      |        |               | local child commit   |
      |        | origin        | origin child commit  |
      | parent | local         | origin parent commit |
      |        |               | local parent commit  |
      |        | origin        | origin parent commit |
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
