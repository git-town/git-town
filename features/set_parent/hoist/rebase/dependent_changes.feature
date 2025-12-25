Feature: remove a branch and all its children from a stack with dependent changes

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT           |
      | main   | local, origin | main commit | file      | line 1\nline 2\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                             |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\nline 2\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                               |
      | branch-2 | local, origin | branch-2 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                 |
      | branch-3 | local, origin | branch-3 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3: branch-3 changes |
    And local Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "branch-2"
    When I run "git-town set-parent main"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      """
    And a rebase is now in progress
    And file "file" now has content:
      """
      <<<<<<< HEAD
      line 1
      line 2
      =======
      line 1: branch-1 changes
      line 2: branch-2 changes
      >>>>>>> {{ sha-short 'branch-2 commit' }} (branch-2 commit)
      line 3
      """

  Scenario: resolve and continue
    When I resolve the conflict in "file" with:
      """
      line 1
      line 2: branch-2 changes
      line 3
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-2 | GIT_EDITOR=true git rebase --continue                          |
      |          | git push --force-with-lease --force-if-includes                |
      |          | git checkout branch-3                                          |
      | branch-3 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-2 branch-1 |
    And a rebase is now in progress
    And file "file" now has content:
      """
      <<<<<<< HEAD
      line 1
      =======
      line 1: branch-1 changes
      >>>>>>> {{ sha-initial-short 'branch-2 commit' }} (branch-2 commit)
      line 2: branch-2 changes
      line 3
      """
    When I resolve the conflict in "file" with:
      """
      line 1
      line 2: branch-2 changes
      line 3
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-3 | GIT_EDITOR=true git rebase --continue           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                               |
      | main     | local, origin | main commit     | file      | line 1\nline 2\nline 3                                     |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\nline 2\nline 3                   |
      | branch-2 | local, origin | branch-2 commit | file      | line 1\nline 2: branch-2 changes\nline 3                   |
      | branch-3 | local, origin | branch-3 commit | file      | line 1\nline 2: branch-2 changes\nline 3: branch-3 changes |
    And this lineage exists now
      """
      main
        branch-1
        branch-2
          branch-3
      """
