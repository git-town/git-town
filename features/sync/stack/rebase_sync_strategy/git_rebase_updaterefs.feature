Feature: stacked changes

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
    And the current branch is "child"
    And Git Town setting "sync-feature-strategy" is "rebase"
    And global Git setting "rebase.updateRefs" is "true"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git fetch --prune --tags                        |
      |        | git checkout main                               |
      | main   | git rebase origin/main --no-update-refs         |
      |        | git push                                        |
      |        | git checkout parent                             |
      | parent | git rebase main --no-update-refs                |
      |        | git push --force-with-lease --force-if-includes |
      |        | git rebase origin/parent --no-update-refs       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout child                              |
      | child  | git rebase parent --no-update-refs              |
      |        | git push --force-with-lease --force-if-includes |
      |        | git rebase origin/child --no-update-refs        |
      |        | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
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

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                         |
      | child  | git reset --hard {{ sha-before-run 'local child commit' }}                                      |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin child commit' }}:child   |
      |        | git checkout parent                                                                             |
      | parent | git reset --hard {{ sha-before-run 'local parent commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin parent commit' }}:parent |
      |        | git checkout child                                                                              |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | child  | local         | local child commit   |
      |        | origin        | origin child commit  |
      | parent | local         | local parent commit  |
      |        | origin        | origin parent commit |
    And the initial branches and lineage exist now
