Feature: make a child branch a sibling in a stack with dependent changes

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
    And the current branch is "branch-3"
    When I run "git-town set-parent branch-1"

  Scenario: result
    And Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-3 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-1 branch-2 |
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      """
    And a rebase is now in progress
    And file "file" now has content:
      """
      line 1: branch-1 changes
      <<<<<<< HEAD
      line 2
      line 3
      =======
      line 2: branch-2 changes
      line 3: branch-3 changes
      >>>>>>> {{ sha-short 'branch-3 commit' }} (branch-3 commit)
      """

  Scenario: resolve and continue
    When I resolve the conflict in "file" with:
      """
      line 1: branch-1 changes
      line 2
      line 3: branch-3 changes
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-3 | GIT_EDITOR=true git rebase --continue           |
      |          | git push --force-with-lease --force-if-includes |
    And no rebase is now in progress
    And this lineage exists now
      """
      main
        branch-1
          branch-2
          branch-3
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                               |
      | main     | local, origin | main commit     | file      | line 1\nline 2\nline 3                                     |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\nline 2\nline 3                   |
      | branch-2 | local, origin | branch-2 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3 |
      | branch-3 | local, origin | branch-3 commit | file      | line 1: branch-1 changes\nline 2\nline 3: branch-3 changes |
