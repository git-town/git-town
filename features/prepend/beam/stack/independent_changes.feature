@messyoutput
Feature: beam a commit from a stack with independent changes into a prepended branch

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT               |
      | main   | local, origin | main commit | file      | line 1\n\nline 2\n\nline 3 |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT                                                                     |
      | old    | local, origin | commit 1 | file      | line 1: commit-1 changes\n\nline 2\n\nline 3                                     |
      | old    | local, origin | commit 2 | file      | line 1: commit-1 changes\n\nline 2: commit-2 changes\n\nline 3                   |
      | old    | local, origin | commit 3 | file      | line 1: commit-1 changes\n\nline 2: commit-2 changes\n\nline 3: commit-3 changes |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "old"
    When I run "git-town prepend new --beam" and enter into the dialog:
      | DIALOG          | KEYS             | COMMENT         |
      | commits to beam | down space enter | select commit 2 |
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                 |
      | old    | git checkout -b new main                                                                                |
      | new    | git cherry-pick {{ sha-initial 'commit 2' }}                                                            |
      |        | git checkout old                                                                                        |
      | old    | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 2' }}^ {{ sha-initial 'commit 2' }} |
      |        | git -c rebase.updateRefs=false rebase new                                                               |
      |        | git push --force-with-lease --force-if-includes                                                         |
      |        | git checkout new                                                                                        |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        new
          old
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                                                                     |
      | main   | local, origin | main commit | file      | line 1\n\nline 2\n\nline 3                                                       |
      | new    | local         | commit 2    | file      | line 1\n\nline 2: commit-2 changes\n\nline 3                                     |
      | old    | local, origin | commit 1    | file      | line 1: commit-1 changes\n\nline 2: commit-2 changes\n\nline 3                   |
      |        |               | commit 3    | file      | line 1: commit-1 changes\n\nline 2: commit-2 changes\n\nline 3: commit-3 changes |
      |        | origin        | commit 2    | file      | line 1\n\nline 2: commit-2 changes\n\nline 3                                     |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | new    | git checkout old                                |
      | old    | git reset --hard {{ sha 'commit 3' }}           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git branch -D new                               |
    And the initial lineage exists now
    And the initial commits exist now

  Scenario: first sync after prepend
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | new    | git fetch --prune --tags |
      |        | git push -u origin new   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | commit 2    |
      | old    | local, origin | commit 1    |
      |        |               | commit 3    |
    And no uncommitted files exist now

  Scenario: sync and amend the beamed commit
    When I amend this commit
      | BRANCH | LOCATION | MESSAGE          | FILE NAME | FILE CONTENT                                         |
      | new    | local    | commit 2 amended | file      | line 1\n\nline 2: amended commit-2 changes\n\nline 3 |
    And the current branch is "old"
    And I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                          |
      | old    | git fetch --prune --tags                                                         |
      |        | git checkout new                                                                 |
      | new    | git push -u origin new                                                           |
      |        | git checkout old                                                                 |
      | old    | git -c rebase.updateRefs=false rebase --onto new {{ sha-before-run 'commit 2' }} |
      |        | git push --force-with-lease --force-if-includes                                  |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE          | FILE NAME | FILE CONTENT                                                                             |
      | main   | local, origin | main commit      | file      | line 1\n\nline 2\n\nline 3                                                               |
      | new    | local, origin | commit 2 amended | file      | line 1\n\nline 2: amended commit-2 changes\n\nline 3                                     |
      | old    | local, origin | commit 1         | file      | line 1: commit-1 changes\n\nline 2: amended commit-2 changes\n\nline 3                   |
      |        |               | commit 3         | file      | line 1: commit-1 changes\n\nline 2: amended commit-2 changes\n\nline 3: commit-3 changes |
