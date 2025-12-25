@messyoutput
Feature: beam a commit from a stack with dependent changes into a prepended branch

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT           |
      | main   | local, origin | main commit | file      | line 1\nline 2\nline 3 |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT                                                                 |
      | old    | local, origin | commit 1 | file      | line 1: commit-1 changes\nline 2\nline 3                                     |
      | old    | local, origin | commit 2 | file      | line 1: commit-1 changes\nline 2: commit-2 changes\nline 3                   |
      | old    | local, origin | commit 3 | file      | line 1: commit-1 changes\nline 2: commit-2 changes\nline 3: commit-3 changes |
    And the current branch is "old"
    When I run "git-town prepend new --beam" and enter into the dialog:
      | DIALOG          | KEYS             | COMMENT         |
      | commits to beam | down space enter | select commit 2 |
    Then Git Town runs the commands
      | BRANCH | COMMAND                                      |
      | old    | git checkout -b new main                     |
      | new    | git cherry-pick {{ sha-initial 'commit 2' }} |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And file "file" now has content:
      """
      <<<<<<< HEAD
      line 1
      line 2
      =======
      line 1: commit-1 changes
      line 2: commit-2 changes
      >>>>>>> {{ sha-short 'commit 2' }} (commit 2)
      line 3
      """
    When I resolve the conflict in "file" with:
      """
      line 1
      line 2: commit-2 changes
      line 3
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                                 |
      | new    | GIT_EDITOR=true git cherry-pick --continue                                                              |
      |        | git checkout old                                                                                        |
      | old    | git -c rebase.updateRefs=false rebase --onto {{ sha-initial 'commit 2' }}^ {{ sha-initial 'commit 2' }} |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And Git Town prints something like:
      """
      error: could not apply .* commit 3
      """
    And file "file" now has content:
      """
      line 1: commit-1 changes
      <<<<<<< HEAD
      line 2
      line 3
      =======
      line 2: commit-2 changes
      line 3: commit-3 changes
      >>>>>>> {{ sha-short 'commit 3' }} (commit 3)
      """
    When I resolve the conflict in "file" with:
      """
      line 1: commit-1 changes
      line 2
      line 3: commit-3 changes
      """
    And I run "git add file"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | old    | GIT_EDITOR=true git rebase --continue     |
      |        | git -c rebase.updateRefs=false rebase new |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And Git Town prints something like:
      """
      error: could not apply .* commit 1
      """
    And file "file" now has content:
      """
      <<<<<<< HEAD
      line 1
      line 2: commit-2 changes
      =======
      line 1: commit-1 changes
      line 2
      >>>>>>> {{ sha-short 'commit 1' }} (commit 1)
      line 3
      """
    When I resolve the conflict in "file" with:
      """
      line 1: commit-1 changes
      line 2: commit-2 changes
      line 3
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | old    | GIT_EDITOR=true git rebase --continue |
    And Git Town prints something like:
      """
      could not apply .* commit 3
      """
    And file "file" now has content:
      """
      line 1: commit-1 changes
      <<<<<<< HEAD
      line 2: commit-2 changes
      line 3
      =======
      line 2
      line 3: commit-3 changes
      >>>>>>> {{ sha-short 'commit 3' }} (commit 3)
      """
    When I resolve the conflict in "file" with:
      """
      line 1: commit-1 changes
      line 2: commit-2 changes
      line 3: commit-3 changes
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | old    | GIT_EDITOR=true git rebase --continue           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout new                                |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                                                                 |
      | main   | local, origin | main commit | file      | line 1\nline 2\nline 3                                                       |
      | new    | local         | commit 2    | file      | line 1\nline 2: commit-2 changes\nline 3                                     |
      | old    | local, origin | commit 1    | file      | line 1: commit-1 changes\nline 2: commit-2 changes\nline 3                   |
      |        |               | commit 3    | file      | line 1: commit-1 changes\nline 2: commit-2 changes\nline 3: commit-3 changes |
      |        | origin        | commit 2    | file      | line 1\nline 2: commit-2 changes\nline 3                                     |
    And this lineage exists now
      """
      main
        new
          old
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | new    | git checkout old                                |
      | old    | git reset --hard {{ sha 'commit 3' }}           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git branch -D new                               |
    And the initial commits exist now
    And the initial lineage exists now

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

  Scenario: amend the beamed commit and sync
    When I amend this commit
      | BRANCH | LOCATION | MESSAGE          | FILE NAME | FILE CONTENT                                     |
      | new    | local    | commit 2 amended | file      | line 1\nline 2: amended commit-2 changes\nline 3 |
    And the current branch is "old"
    And I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                          |
      | old    | git fetch --prune --tags                                                         |
      |        | git checkout new                                                                 |
      | new    | git push -u origin new                                                           |
      |        | git checkout old                                                                 |
      | old    | git -c rebase.updateRefs=false rebase --onto new {{ sha-before-run 'commit 2' }} |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And file "file" now has content:
      """
      <<<<<<< HEAD
      line 1
      line 2: amended commit-2 changes
      =======
      line 1: commit-1 changes
      line 2: commit-2 changes
      >>>>>>> {{ sha-short 'commit 1' }} (commit 1)
      line 3
      """
    When I resolve the conflict in "file" with:
      """
      line 1: commit-1 changes
      line 2: amended commit-2 changes
      line 3
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | old    | GIT_EDITOR=true git rebase --continue |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And file "file" now has content:
      """
      line 1: commit-1 changes
      <<<<<<< HEAD
      line 2: amended commit-2 changes
      line 3
      =======
      line 2: commit-2 changes
      line 3: commit-3 changes
      >>>>>>> {{ sha-short 'commit 3' }} (commit 3)
      """
    When I resolve the conflict in "file" with:
      """
      line 1: commit-1 changes
      line 2: amended commit-2 changes
      line 3: commit-3 changes
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | old    | GIT_EDITOR=true git rebase --continue           |
      |        | git push --force-with-lease --force-if-includes |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE          | FILE NAME | FILE CONTENT                                                                         |
      | main   | local, origin | main commit      | file      | line 1\nline 2\nline 3                                                               |
      | new    | local, origin | commit 2 amended | file      | line 1\nline 2: amended commit-2 changes\nline 3                                     |
      | old    | local, origin | commit 1         | file      | line 1: commit-1 changes\nline 2: amended commit-2 changes\nline 3                   |
      |        |               | commit 3         | file      | line 1: commit-1 changes\nline 2: amended commit-2 changes\nline 3: commit-3 changes |
