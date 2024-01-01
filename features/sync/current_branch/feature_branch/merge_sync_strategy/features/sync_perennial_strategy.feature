Feature: with sync-perennial-strategy set to "merge"

  Background:
    Given Git Town setting "sync-perennial-strategy" is "merge"
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git merge --no-edit origin/main    |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | local, origin | local main commit                                          |
      |         |               | origin main commit                                         |
      |         |               | Merge remote-tracking branch 'origin/main'                 |
      | feature | local, origin | local feature commit                                       |
      |         |               | origin feature commit                                      |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |
      |         |               | local main commit                                          |
      |         |               | origin main commit                                         |
      |         |               | Merge remote-tracking branch 'origin/main'                 |
      |         |               | Merge branch 'main' into feature                           |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | git reset --hard {{ sha 'local feature commit' }}                            |
      |         | git push --force-with-lease origin {{ sha 'origin feature commit' }}:feature |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                    |
      | main    | local, origin | local main commit                          |
      |         |               | origin main commit                         |
      |         |               | Merge remote-tracking branch 'origin/main' |
      | feature | local         | local feature commit                       |
      |         | origin        | origin feature commit                      |
    And the initial branches and lineage exist
