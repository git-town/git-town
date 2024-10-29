Feature: sync a grandchild feature branch using the "compress" strategy

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
    And wait 1 second to ensure new Git timestamps
    And Git Town setting "sync-feature-strategy" is "compress"
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | child  | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git push                                |
      |        | git checkout parent                     |
      | parent | git merge --no-edit --ff origin/parent  |
      |        | git merge --no-edit --ff main           |
      |        | git reset --soft main                   |
      |        | git commit -m "local parent commit"     |
      |        | git push --force-with-lease             |
      |        | git checkout child                      |
      | child  | git merge --no-edit --ff origin/child   |
      |        | git merge --no-edit --ff parent         |
      |        | git reset --soft parent                 |
      |        | git commit -m "local child commit"      |
      |        | git push --force-with-lease             |
    And all branches are now synchronized
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | local, origin | origin main commit  |
      |        |               | local main commit   |
      | child  | local, origin | local child commit  |
      | parent | local, origin | local parent commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
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
