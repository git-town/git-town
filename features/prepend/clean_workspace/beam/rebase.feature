Feature: prepend a branch to a feature branch using the "rebase" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  |
      | old    | local, origin | commit 1 |
      | old    | local, origin | commit 2 |
      | old    | local, origin | commit 3 |
      | old    | local, origin | commit 4 |
    And the current branch is "old"
    And Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town prepend parent --beam" and enter into the dialog:
      | DIALOG                 | KEYS                             |
      | select commits 2 and 4 | down space down down space enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | old    | git fetch --prune --tags                        |
      |        | git checkout main                               |
      | main   | git rebase origin/main --no-update-refs         |
      |        | git checkout old                                |
      | old    | git rebase main --no-update-refs                |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout -b parent main                     |
      | parent | git cherry-pick {{ sha-before-run 'commit 2' }} |
      |        | git cherry-pick {{ sha-before-run 'commit 4' }} |
      |        | git checkout old                                |
      | old    | git rebase parent --no-update-refs              |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout parent                             |
    And the current branch is now "parent"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | old    | local, origin | commit 1 |
      |        |               | commit 3 |
      |        | origin        | commit 2 |
      |        |               | commit 4 |
      | parent | local         | commit 2 |
      |        |               | commit 4 |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |
    When I run "git town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | parent | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git checkout parent                     |
      | parent | git rebase main --no-update-refs        |
      |        | git push -u origin parent               |
    And the current branch is now "parent"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | old    | local, origin | commit 1 |
      |        |               | commit 3 |
      | parent | local, origin | commit 2 |
      |        |               | commit 4 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And the current branch is now "old"
    And the initial commits exist now
    And the initial lineage exists now
