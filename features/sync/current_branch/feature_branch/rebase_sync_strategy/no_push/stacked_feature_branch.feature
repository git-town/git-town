Feature: syncing a stacked feature branch using --no-push

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "child"
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
    And Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town sync --no-push"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | child  | git fetch --prune --tags                  |
      |        | git checkout main                         |
      | main   | git rebase origin/main --no-update-refs   |
      |        | git checkout parent                       |
      | parent | git rebase main --no-update-refs          |
      |        | git rebase origin/parent --no-update-refs |
      |        | git checkout child                        |
      | child  | git rebase parent --no-update-refs        |
      |        | git rebase origin/child --no-update-refs  |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        | local         | local main commit    |
      | child  | local, origin | origin child commit  |
      |        | local         | origin parent commit |
      |        |               | origin main commit   |
      |        |               | local main commit    |
      |        |               | local parent commit  |
      |        |               | local child commit   |
      | parent | local, origin | origin parent commit |
      |        | local         | origin main commit   |
      |        |               | local main commit    |
      |        |               | local parent commit  |
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

    And the current branch is still "child"
    And the initial commits exist now
    And the initial branches and lineage exist now
