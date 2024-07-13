Feature: syncing a stacked feature branch using --no-push

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
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
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git fetch --prune --tags                        |
      |        | git checkout main                               |
      | main   | git rebase origin/main                          |
      |        | git checkout parent                             |
      | parent | git rebase main                                 |
      |        | git push --force-with-lease --force-if-includes |
      |        | git rebase origin/parent                        |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout child                              |
      | child  | git rebase parent                               |
      |        | git push --force-with-lease --force-if-includes |
      |        | git rebase origin/child                         |
      |        | git push --force-with-lease --force-if-includes |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        | local         | local main commit    |
      | child  | local, origin | origin child commit  |
      |        |               | origin parent commit |
      |        |               | origin main commit   |
      |        |               | local main commit    |
      |        |               | local parent commit  |
      |        |               | local child commit   |
      | parent | local, origin | origin parent commit |
      |        |               | origin main commit   |
      |        |               | local main commit    |
      |        |               | local parent commit  |
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                              |
      | child  | git reset --hard {{ sha 'local child commit' }}                                      |
      |        | git push --force-with-lease origin {{ sha 'origin child commit' }}:child             |
      |        | git checkout parent                                                                  |
      | parent | git reset --hard {{ sha 'local parent commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin 'origin parent commit' }}:parent |
      |        | git checkout main                                                                    |
      | main   | git reset --hard {{ sha 'local main commit' }}                                       |
      |        | git checkout child                                                                   |

    And the current branch is still "child"
    And the initial commits exist
    And the initial branches and lineage exist
